package configurationservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/ent/configuration"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrUpdateConfigurationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrUpdateConfigurationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrUpdateConfigurationLogic {
	return &CreateOrUpdateConfigurationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建或更新配置
// 如果配置存在（根据 key），则更新；否则创建新配置
func (l *CreateOrUpdateConfigurationLogic) CreateOrUpdateConfiguration(in *core.ConfigurationInfo) (*core.ConfigurationInfo, error) {
	// 验证必填字段
	if in.Key == "" || in.Value == "" || in.Name == "" || in.Group == "" {
		return nil, errorx.NewInvalidArgumentError("key、value、name、group 不能为空")
	}

	// 检查写权限
	permChecker := NewConfigurationPermissionChecker(l.svcCtx.Casbin, l.Logger)
	if err := permChecker.CheckWritePermission(l.ctx, in.Group); err != nil {
		return nil, err
	}

	// 先查询配置是否存在
	existingConfig, err := l.svcCtx.DBEnt.Configuration.Query().
		Where(configuration.KeyEQ(in.Key)).
		First(l.ctx)

	if err != nil && !ent.IsNotFound(err) {
		// 数据库查询错误（不是"未找到"错误）
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	var result *core.ConfigurationInfo

	if existingConfig != nil {
		// 配置存在，执行更新
		updatedConfig, err := l.svcCtx.DBEnt.Configuration.UpdateOne(existingConfig).
			SetValue(in.Value).
			SetName(in.Name).
			SetGroup(in.Group).
			SetNillableDescription(in.Description).
			Save(l.ctx)
		if err != nil {
			return nil, errorhandler.DBEntError(l.Logger, err, in)
		}

		// 更新缓存
		l.svcCtx.ConfigurationCache.Set(updatedConfig.Key, updatedConfig.Value)

		result = ConvertConfigurationToConfigurationInfo(updatedConfig)
	} else {
		// 配置不存在，执行创建
		newConfig, err := l.svcCtx.DBEnt.Configuration.Create().
			SetKey(in.Key).
			SetValue(in.Value).
			SetName(in.Name).
			SetGroup(in.Group).
			SetNillableDescription(in.Description).
			Save(l.ctx)
		if err != nil {
			return nil, errorhandler.DBEntError(l.Logger, err, in)
		}

		// 添加到缓存
		l.svcCtx.ConfigurationCache.Set(newConfig.Key, newConfig.Value)

		result = ConvertConfigurationToConfigurationInfo(newConfig)
	}

	return result, nil
}
