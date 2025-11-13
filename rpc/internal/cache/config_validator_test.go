package cache

import (
	"testing"
)

func TestConfigValidator(t *testing.T) {
	// 创建配置缓存
	cache := NewConfigurationCache()

	// 初始化测试数据
	testConfigs := map[string]string{
		"enable_feature_x": "true",
		"min_user_level":   "50",
		"admin_email":      "admin@example.com",
		"allowed_user_ids": `[100, 201, 305]`,
		"feature_flags":    `{"beta": true, "max_items": 25, "theme": "dark"}`,
	}
	cache.SetAll(testConfigs)

	// 创建验证器
	validator, err := NewConfigValidator(cache)
	if err != nil {
		t.Fatalf("创建验证器失败: %v", err)
	}

	// 测试用例
	tests := []struct {
		name      string
		key       string
		exp       string
		wantValid bool
		wantErr   bool
	}{
		{
			name:      "布尔值比较 - true",
			key:       "enable_feature_x",
			exp:       `value == true`,
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "布尔值比较 - false",
			key:       "enable_feature_x",
			exp:       `value == false`,
			wantValid: false,
			wantErr:   false,
		},
		{
			name:      "数字比较 - 大于",
			key:       "min_user_level",
			exp:       `int(value) > 20`,
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "数字比较 - 小于",
			key:       "min_user_level",
			exp:       `int(value) < 40`,
			wantValid: false,
			wantErr:   false,
		},
		{
			name:      "JSON 对象字段访问",
			key:       "feature_flags",
			exp:       `value.theme == "dark"`,
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "JSON 对象字段比较",
			key:       "feature_flags",
			exp:       `value.max_items > 20`,
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "JSON 数组包含",
			key:       "allowed_user_ids",
			exp:       `201 in value`,
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "JSON 数组不包含",
			key:       "allowed_user_ids",
			exp:       `999 in value`,
			wantValid: false,
			wantErr:   false,
		},
		{
			name:      "字符串方法 - startsWith",
			key:       "admin_email",
			exp:       `value.startsWith("admin")`,
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "字符串方法 - contains",
			key:       "admin_email",
			exp:       `value.contains("@example")`,
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "配置不存在",
			key:       "non_existent_key",
			exp:       `value == "test"`,
			wantValid: false,
			wantErr:   true,
		},
		{
			name:      "表达式返回非bool",
			key:       "min_user_level",
			exp:       `int(value) + 10`,
			wantValid: false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validator.Validate(tt.key, tt.exp)

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result != tt.wantValid {
				t.Errorf("Validate() got %v, want %v", result, tt.wantValid)
			}
		})
	}
}

func TestConfigValidatorEvaluate(t *testing.T) {
	// 创建配置缓存
	cache := NewConfigurationCache()

	// 初始化测试数据
	testConfigs := map[string]string{
		"min_user_level": "50",
		"feature_flags":  `{"beta": true, "max_items": 25, "theme": "dark"}`,
		"allowed_user_ids": `[100, 201, 305]`,
	}
	cache.SetAll(testConfigs)

	// 创建验证器
	validator, err := NewConfigValidator(cache)
	if err != nil {
		t.Fatalf("创建验证器失败: %v", err)
	}

	// 测试 Evaluate 函数
	tests := []struct {
		name    string
		key     string
		exp     string
		wantErr bool
	}{
		{
			name:    "提取数字值",
			key:     "min_user_level",
			exp:     `int(value)`,
			wantErr: false,
		},
		{
			name:    "提取 JSON 字段",
			key:     "feature_flags",
			exp:     `value.theme`,
			wantErr: false,
		},
		{
			name:    "提取 JSON 数组元素",
			key:     "allowed_user_ids",
			exp:     `value[1]`,
			wantErr: false,
		},
		{
			name:    "三元运算符",
			key:     "feature_flags",
			exp:     `value.beta ? "Beta开启" : "Beta关闭"`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validator.Evaluate(tt.key, tt.exp)

			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result == nil {
				t.Errorf("Evaluate() got nil result")
			}
		})
	}
}

func TestProgramCache(t *testing.T) {
	// 创建配置缓存
	cache := NewConfigurationCache()
	cache.SetAll(map[string]string{
		"test_key": "test_value",
	})

	// 创建验证器
	validator, err := NewConfigValidator(cache)
	if err != nil {
		t.Fatalf("创建验证器失败: %v", err)
	}

	// 验证初始缓存大小为 0
	if size := validator.GetProgramCacheSize(); size != 0 {
		t.Errorf("初始缓存大小应为 0, 实际为 %d", size)
	}

	// 执行验证，应该缓存 Program
	_, _ = validator.Validate("test_key", `value == "test_value"`)

	// 验证缓存大小为 1
	if size := validator.GetProgramCacheSize(); size != 1 {
		t.Errorf("缓存大小应为 1, 实际为 %d", size)
	}

	// 执行相同的验证，应该使用缓存
	_, _ = validator.Validate("test_key", `value == "test_value"`)

	// 验证缓存大小仍为 1
	if size := validator.GetProgramCacheSize(); size != 1 {
		t.Errorf("缓存大小应为 1, 实际为 %d", size)
	}

	// 执行不同的验证，应该添加新的 Program
	_, _ = validator.Validate("test_key", `value != "other_value"`)

	// 验证缓存大小为 2
	if size := validator.GetProgramCacheSize(); size != 2 {
		t.Errorf("缓存大小应为 2, 实际为 %d", size)
	}

	// 清空缓存
	validator.ClearProgramCache()

	// 验证缓存大小为 0
	if size := validator.GetProgramCacheSize(); size != 0 {
		t.Errorf("清空后缓存大小应为 0, 实际为 %d", size)
	}
}

