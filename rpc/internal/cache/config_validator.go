package cache

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
)

// ConfigValidator 封装了配置和CEL求值逻辑
// 支持 Program 缓存以提高性能
type ConfigValidator struct {
	cache          *ConfigurationCache // 配置缓存
	programCache   sync.Map            // CEL Program 缓存
	env            *cel.Env            // CEL 环境
	mu             sync.RWMutex        // 保护 env 的读写
}

// NewConfigValidator 创建一个新的验证器实例
func NewConfigValidator(cache *ConfigurationCache) (*ConfigValidator, error) {
	env, err := cel.NewEnv(
		cel.StdLib(),
		cel.Declarations(
			decls.NewVar("value", decls.Dyn),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("创建CEL环境失败: %w", err)
	}

	return &ConfigValidator{
		cache:        cache,
		programCache: sync.Map{},
		env:          env,
	}, nil
}

// prepareValue 从配置缓存中获取字符串值并尝试将其解析为JSON
// 如果不是JSON，则按原样返回字符串
func (cv *ConfigValidator) prepareValue(key string) (interface{}, error) {
	stringValue, ok := cv.cache.Get(key)
	if !ok {
		return nil, fmt.Errorf("配置项 '%s' 未找到", key)
	}

	// 分割 group 和 value
	_, value, err := splitGroupAndValue(stringValue)
	if err != nil {
		return nil, err
	}

	var evalValue interface{}
	var jsonData interface{}

	// 尝试将字符串解析为JSON
	err = json.Unmarshal([]byte(value), &jsonData)
	if err == nil {
		// 解析成功！
		evalValue = jsonData
	} else {
		// 解析失败，它只是一个普通字符串
		evalValue = value
	}
	return evalValue, nil
}

// getProgram 编译（或从缓存中检索）CEL表达式
// 使用 sync.Map 实现无锁缓存
func (cv *ConfigValidator) getProgram(exp string) (cel.Program, error) {
	// 先尝试从缓存获取
	if prog, ok := cv.programCache.Load(exp); ok {
		return prog.(cel.Program), nil
	}

	cv.mu.RLock()
	defer cv.mu.RUnlock()

	// 1. 解析表达式
	ast, iss := cv.env.Parse(exp)
	if iss.Err() != nil {
		return nil, fmt.Errorf("解析表达式失败: %w", iss.Err())
	}

	// 2. 类型检查
	checked, iss := cv.env.Check(ast)
	if iss.Err() != nil {
		return nil, fmt.Errorf("类型检查失败: %w", iss.Err())
	}

	// 3. 编译为 Program
	prog, err := cv.env.Program(checked)
	if err != nil {
		return nil, fmt.Errorf("创建Program失败: %w", err)
	}

	// 4. 缓存 Program
	cv.programCache.Store(exp, prog)
	return prog, nil
}

// Validate 验证函数，返回 (bool, error)
// 在运行时检查结果是否为 bool
func (cv *ConfigValidator) Validate(key string, exp string) (bool, error) {
	// 1. 获取编译后的CEL program
	prog, err := cv.getProgram(exp)
	if err != nil {
		return false, fmt.Errorf("获取 program 失败 (exp: '%s'): %w", exp, err)
	}

	// 2. 准备数据
	evalValue, err := cv.prepareValue(key)
	if err != nil {
		return false, err
	}

	// 3. 求值
	out, _, err := prog.Eval(map[string]interface{}{
		"value": evalValue,
	})
	if err != nil {
		return false, fmt.Errorf("表达式求值失败 (key: '%s', exp: '%s'): %w", key, exp, err)
	}

	// 4. 转换结果 (运行时检查)
	result, ok := out.Value().(bool)
	if !ok {
		// 表达式计算结果不是 bool
		return false, fmt.Errorf("表达式 '%s' 未返回bool, 实际返回: %T (%v)", exp, out.Value(), out.Value())
	}

	return result, nil
}

// Evaluate 求值函数，返回 (interface{}, error)
// 返回表达式计算的任何结果
func (cv *ConfigValidator) Evaluate(key string, exp string) (interface{}, error) {
	// 1. 获取编译后的CEL program
	prog, err := cv.getProgram(exp)
	if err != nil {
		return nil, fmt.Errorf("获取 program 失败 (exp: '%s'): %w", exp, err)
	}

	// 2. 准备数据
	evalValue, err := cv.prepareValue(key)
	if err != nil {
		return nil, err
	}

	// 3. 求值
	out, _, err := prog.Eval(map[string]interface{}{
		"value": evalValue,
	})
	if err != nil {
		return nil, fmt.Errorf("表达式求值失败 (key: '%s', exp: '%s'): %w", key, exp, err)
	}

	// 4. 直接返回结果
	return out.Value(), nil
}

// ClearProgramCache 清空 Program 缓存
// 用于测试或需要重新编译所有表达式的场景
func (cv *ConfigValidator) ClearProgramCache() {
	cv.programCache.Range(func(key, value interface{}) bool {
		cv.programCache.Delete(key)
		return true
	})
}

// GetProgramCacheSize 获取 Program 缓存大小
func (cv *ConfigValidator) GetProgramCacheSize() int {
	size := 0
	cv.programCache.Range(func(key, value interface{}) bool {
		size++
		return true
	})
	return size
}


func splitGroupAndValue(s string) (string, string, error) {
	parts := strings.Split(s, "<>")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("配置项 '%s' 格式错误", s)
	}
	return parts[0], parts[1], nil
}
