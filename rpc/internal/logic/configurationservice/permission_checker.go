package configurationservicelogic

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/wenpiner/last-admin-common/ctx/rolectx"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	// 权限资源类型
	ConfigurationResourceType = "configuration"
	// 权限操作
	OperationRead  = "read"
	OperationWrite = "write"
)

// ConfigurationPermissionChecker 配置权限检查器
type ConfigurationPermissionChecker struct {
	casbin *casbin.Enforcer
	logger logx.Logger
}

// NewConfigurationPermissionChecker 创建权限检查器
func NewConfigurationPermissionChecker(casbin *casbin.Enforcer, logger logx.Logger) *ConfigurationPermissionChecker {
	return &ConfigurationPermissionChecker{
		casbin: casbin,
		logger: logger,
	}
}

// CheckReadPermission 检查读权限
// 返回 error 表示无权限或检查失败
func (c *ConfigurationPermissionChecker) CheckReadPermission(ctx context.Context, group string) error {
	roleIds, ok := rolectx.GetRoleFromContext(ctx)
	if !ok || len(roleIds) == 0 {
		c.logger.Errorw("从上下文中读取Role失败", logx.Field("detail", "roleIds is empty"))
		return errorx.NewInvalidArgumentError("common.configuration.forbidden")
	}

	return c.checkPermission(roleIds, group, OperationRead)
}

// CheckWritePermission 检查写权限
// 返回 error 表示无权限或检查失败
func (c *ConfigurationPermissionChecker) CheckWritePermission(ctx context.Context, group string) error {
	roleIds, ok := rolectx.GetRoleFromContext(ctx)
	if !ok || len(roleIds) == 0 {
		c.logger.Errorw("从上下文中读取Role失败", logx.Field("detail", "roleIds is empty"))
		return errorx.NewInvalidArgumentError("common.configuration.forbidden")
	}

	return c.checkPermission(roleIds, group, OperationWrite)
}

// GetAllowedGroups 获取允许的配置分组列表
// operation: "read" 或 "write"，如果为空则返回所有操作的分组
// 返回允许的分组列表，如果无权限则返回空列表
func (c *ConfigurationPermissionChecker) GetAllowedGroups(ctx context.Context, operation string) ([]string, error) {
	roleIds, ok := rolectx.GetRoleFromContext(ctx)
	if !ok || len(roleIds) == 0 {
		c.logger.Errorw("从上下文中读取Role失败", logx.Field("detail", "roleIds is empty"))
		return []string{}, nil
	}

	// 使用 map 去重
	groupSet := make(map[string]bool)

	for _, roleId := range roleIds {
		// 获取该角色的所有配置权限策略
		// 策略格式: [roleId, "configuration", group, operation]
		policies, err := c.casbin.GetFilteredPolicy(0, roleId, ConfigurationResourceType)
		if err != nil {
			c.logger.Errorw("获取 Casbin 策略失败", logx.Field("detail", err.Error()))
			return nil, err
		}

		// 过滤分组
		for _, policy := range policies {
			// policy 格式: [roleId, "configuration", group, operation]
			if len(policy) >= 4 {
				// 如果指定了 operation，则只过滤该操作的分组
				// 如果未指定 operation（为空），则返回所有操作的分组
				if operation == "" || policy[3] == operation {
					groupSet[policy[2]] = true
				}
			}
		}
	}

	// 转换为切片
	groups := make([]string, 0, len(groupSet))
	for group := range groupSet {
		groups = append(groups, group)
	}

	return groups, nil
}

// checkPermission 内部方法：检查权限
// 如果任何一个角色拥有指定的权限，则返回 nil
func (c *ConfigurationPermissionChecker) checkPermission(roleIds []string, group string, operation string) error {
	for _, roleId := range roleIds {
		// 检查该角色是否拥有指定的权限
		// 策略格式: [roleId, "configuration", group, operation]
		policies, err := c.casbin.GetFilteredPolicy(0, roleId, ConfigurationResourceType, group, operation)
		if err != nil {
			c.logger.Errorw("获取 Casbin 策略失败", logx.Field("detail", err.Error()))
			return err
		}

		if len(policies) > 0 {
			// 找到了权限
			return nil
		}
	}

	// 没有找到任何权限
	return errorx.NewInvalidArgumentError("common.configuration.forbidden")
}

// CheckPermissionWithMessage 检查权限并返回自定义错误消息
func (c *ConfigurationPermissionChecker) CheckPermissionWithMessage(ctx context.Context, group string, operation string, errorMessage string) error {
	roleIds, ok := rolectx.GetRoleFromContext(ctx)
	if !ok || len(roleIds) == 0 {
		c.logger.Errorw("从上下文中读取Role失败", logx.Field("detail", "roleIds is empty"))
		return errorx.NewInvalidArgumentError(errorMessage)
	}

	err := c.checkPermission(roleIds, group, operation)
	if err != nil {
		return errorx.NewInvalidArgumentError(errorMessage)
	}
	return nil
}

// DebugPolicies 调试方法：打印所有策略（仅用于开发）
func (c *ConfigurationPermissionChecker) DebugPolicies(roleId string) {
	policies, err := c.casbin.GetFilteredPolicy(0, roleId)
	if err != nil {
		c.logger.Errorw("获取策略失败", logx.Field("detail", err.Error()))
		return
	}

	c.logger.Infow("角色策略", logx.Field("roleId", roleId), logx.Field("policies", fmt.Sprintf("%v", policies)))
}

