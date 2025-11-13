package configurationservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type ValidateConfigurationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewValidateConfigurationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidateConfigurationLogic {
	return &ValidateConfigurationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 验证某个配置的值是否有效
// 使用 CEL 表达式进行高级验证，支持 Program 缓存
func (l *ValidateConfigurationLogic) ValidateConfiguration(in *core.ValidateConfigurationRequest) (*core.ValidateConfigurationResponse, error) {
	// 从缓存获取配置
	result, ok := l.svcCtx.ConfigurationCache.Get(in.Key)
	if !ok {
		return &core.ValidateConfigurationResponse{
			IsValid: false,
			Message: "common.configuration.notFound",
		}, nil
	}

	// 验证分组是否有权限
	group, _, err := splitGroupAndValue(result)
	if err != nil {
		return &core.ValidateConfigurationResponse{
			IsValid: false,
			Message: err.Error(),
		}, nil
	}

	// 检查读权限
	permChecker := NewConfigurationPermissionChecker(l.svcCtx.Casbin, l.Logger)
	if err := permChecker.CheckReadPermission(l.ctx, group); err != nil {
		return &core.ValidateConfigurationResponse{
			IsValid: false,
			Message: "common.forbidden",
		}, nil
	}

	// 使用 ConfigValidator 进行验证
	validateResult, err := l.svcCtx.ConfigValidator.Validate(in.Key, in.Exp)
	if err != nil {
		l.Logger.Errorf("配置验证失败: %v", err)
		return &core.ValidateConfigurationResponse{
			IsValid: false,
			Message: err.Error(),
		}, nil
	}

	// 验证结果
	return &core.ValidateConfigurationResponse{
		IsValid: validateResult,
		Message: "common.success",
	}, nil
}
