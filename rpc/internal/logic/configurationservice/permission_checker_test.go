package configurationservicelogic

import (
	"context"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/zeromicro/go-zero/core/logx"
)

// TestCheckReadPermissionPublicGroup 测试公开配置（/public 开头）是否跳过权限检查
func TestCheckReadPermissionPublicGroup(t *testing.T) {
	// 创建一个模拟的 logger
	logger := logx.WithContext(nil)

	// 创建一个模拟的 casbin enforcer（不需要真实的数据库）
	var mockCasbin *casbin.Enforcer

	checker := &ConfigurationPermissionChecker{
		casbin: mockCasbin,
		logger: logger,
	}

	// 测试用例 1: 公开配置（/public 开头）应该返回 nil，即使没有角色信息
	err := checker.CheckReadPermission(context.Background(), "/public/config")
	if err != nil {
		t.Errorf("Expected nil for /public/config, got error: %v", err)
	}

	// 测试用例 2: 公开配置的子路径也应该返回 nil
	err = checker.CheckReadPermission(context.Background(), "/public/system/config")
	if err != nil {
		t.Errorf("Expected nil for /public/system/config, got error: %v", err)
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

	// 测试用例: 公开配置写操作也应该返回 nil，即使没有角色信息
	err := checker.CheckWritePermission(context.Background(), "/public/config")
	if err != nil {
		t.Errorf("Expected nil for /public/config write, got error: %v", err)
	}
}

// TestCheckReadPermissionPublicGroupVariations 测试 /public 前缀的各种变体
func TestCheckReadPermissionPublicGroupVariations(t *testing.T) {
	logger := logx.WithContext(nil)
	var mockCasbin *casbin.Enforcer

	checker := &ConfigurationPermissionChecker{
		casbin: mockCasbin,
		logger: logger,
	}

	testCases := []struct {
		group       string
		description string
	}{
		{"/public", "exact /public"},
		{"/public/", "/public/"},
		{"/public/config", "/public/config"},
		{"/public/system/config", "/public/system/config"},
		{"/public123", "/public123 (starts with /public)"},
		{"/publicconfig", "/publicconfig (starts with /public)"},
	}

	for _, tc := range testCases {
		err := checker.CheckReadPermission(context.Background(), tc.group)
		if err != nil {
			t.Errorf("%s: Expected nil, got error: %v", tc.description, err)
		}
	}
}
