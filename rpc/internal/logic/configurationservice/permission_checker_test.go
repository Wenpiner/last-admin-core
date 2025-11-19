package configurationservicelogic

import (
	"context"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/zeromicro/go-zero/core/logx"
)

// TestCheckReadPermissionPublicGroup 测试公开配置（/public 开头）是否跳过权限检查
func TestCheckReadPermissionPublicGroup(t *testing.T) {
	logger := logx.WithContext(nil)
	var mockCasbin *casbin.Enforcer

	checker := &ConfigurationPermissionChecker{
		casbin: mockCasbin,
		logger: logger,
	}

	err := checker.CheckReadPermission(context.Background(), "/public/config")
	if err != nil {
		t.Errorf("Expected nil for /public/config, got error: %v", err)
	}
}

// TestCheckWritePermissionPublicGroup 测试公开配置写权限
func TestCheckWritePermissionPublicGroup(t *testing.T) {
	logger := logx.WithContext(nil)
	var mockCasbin *casbin.Enforcer

	checker := &ConfigurationPermissionChecker{
		casbin: mockCasbin,
		logger: logger,
	}

	err := checker.CheckWritePermission(context.Background(), "/public/config")
	if err != nil {
		t.Errorf("Expected nil for /public/config write, got error: %v", err)
	}
}

// TestCheckReadPermissionWithWritePermission 测试权限继承：write 权限包含 read 权限
func TestCheckReadPermissionWithWritePermission(t *testing.T) {
	logger := logx.WithContext(nil)

	m, err := model.NewModelFromString(`
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && r.dom == p.dom && r.obj == p.obj && permissionMatch(p.act, r.act)
`)
	if err != nil {
		t.Fatalf("Failed to create casbin model: %v", err)
	}

	enforcer, err := casbin.NewEnforcer(m)
	if err != nil {
		t.Fatalf("Failed to create casbin enforcer: %v", err)
	}

	// 注册自定义权限匹配函数
	permissionMatchFunc := func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return false, nil
		}

		grantedAction, ok1 := args[0].(string)
		requiredAction, ok2 := args[1].(string)

		if !ok1 || !ok2 {
			return false, nil
		}

		if grantedAction == requiredAction {
			return true, nil
		}

		if requiredAction == "read" && grantedAction == "write" {
			return true, nil
		}

		if grantedAction == "admin" {
			return true, nil
		}

		return false, nil
	}
	enforcer.AddFunction("permissionMatch", permissionMatchFunc)

	checker := &ConfigurationPermissionChecker{
		casbin: enforcer,
		logger: logger,
	}

	enforcer.AddPolicy("role1", "configuration", "/app/config", "write")

	err = checker.checkPermission([]string{"role1"}, "/app/config", OperationWrite)
	if err != nil {
		t.Errorf("Expected nil for write permission, got error: %v", err)
	}

	err = checker.checkPermission([]string{"role1"}, "/app/config", OperationRead)
	if err != nil {
		t.Errorf("Expected nil for read permission (inherited from write), got error: %v", err)
	}

	err = checker.checkPermission([]string{"role1"}, "/other/config", OperationRead)
	if err == nil {
		t.Errorf("Expected error for non-existent permission, got nil")
	}
}

// TestCasbinInitializationWithCustomFunction 测试 Casbin 初始化时自定义函数的注册
func TestCasbinInitializationWithCustomFunction(t *testing.T) {
	// 创建 model
	m, err := model.NewModelFromString(`
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.dom == p.dom && keyMatch2(r.obj, p.obj) && permissionMatch(p.act, r.act)
`)
	if err != nil {
		t.Fatalf("Failed to create casbin model: %v", err)
	}

	// 创建 enforcer（不加载策略）
	enforcer, err := casbin.NewEnforcer(m)
	if err != nil {
		t.Fatalf("Failed to create casbin enforcer: %v", err)
	}

	// 注册自定义权限匹配函数
	permissionMatchFunc := func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return false, nil
		}

		grantedAction, ok1 := args[0].(string)
		requiredAction, ok2 := args[1].(string)

		if !ok1 || !ok2 {
			return false, nil
		}

		if grantedAction == requiredAction {
			return true, nil
		}

		if requiredAction == "read" && grantedAction == "write" {
			return true, nil
		}

		if grantedAction == "admin" {
			return true, nil
		}

		return false, nil
	}
	enforcer.AddFunction("permissionMatch", permissionMatchFunc)

	// 添加策略
	enforcer.AddPolicy("role1", "configuration", "/app/config", "write")

	// 第一次调用 Enforce - 这会编译 matcher 表达式
	allowed, err := enforcer.Enforce("role1", "configuration", "/app/config", "write")
	if err != nil {
		t.Fatalf("First Enforce call failed: %v", err)
	}
	if !allowed {
		t.Errorf("Expected allowed=true for write permission, got false")
	}

	// 第二次调用 Enforce - 这会使用缓存的 matcher 表达式
	allowed, err = enforcer.Enforce("role1", "configuration", "/app/config", "read")
	if err != nil {
		t.Fatalf("Second Enforce call failed: %v", err)
	}
	if !allowed {
		t.Errorf("Expected allowed=true for read permission (inherited from write), got false")
	}
}

// TestCasbinEnforceWithoutFunction 测试在没有注册自定义函数的情况下调用 Enforce
func TestCasbinEnforceWithoutFunction(t *testing.T) {
	// 创建 model
	m, err := model.NewModelFromString(`
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.dom == p.dom && keyMatch2(r.obj, p.obj) && permissionMatch(p.act, r.act)
`)
	if err != nil {
		t.Fatalf("Failed to create casbin model: %v", err)
	}

	// 创建 enforcer（不加载策略）
	enforcer, err := casbin.NewEnforcer(m)
	if err != nil {
		t.Fatalf("Failed to create casbin enforcer: %v", err)
	}

	// 注意：这里故意不注册 permissionMatch 函数

	// 添加策略
	enforcer.AddPolicy("role1", "configuration", "/app/config", "write")

	// 尝试调用 Enforce - 这应该会失败，因为 permissionMatch 函数没有被注册
	allowed, err := enforcer.Enforce("role1", "configuration", "/app/config", "write")
	if err == nil {
		t.Errorf("Expected error when permissionMatch function is not registered, but got allowed=%v", allowed)
	} else {
		t.Logf("Got expected error: %v", err)
	}
}

