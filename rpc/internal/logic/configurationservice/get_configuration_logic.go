package configurationservicelogic

import (
	"context"
	"fmt"
	"strings"

	"github.com/wenpiner/last-admin-core/rpc/ent/configuration"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConfigurationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConfigurationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConfigurationLogic {
	return &GetConfigurationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取配置
func (l *GetConfigurationLogic) GetConfiguration(in *core.StringRequest) (*core.ConfigurationInfo, error) {
	// 先从缓存获取
	if value, ok := l.svcCtx.ConfigurationCache.Get(in.Value); ok {
		// 分割 group 和 value
		group, value, err := splitGroupAndValue(value)
		if err != nil {
			return nil, err
		}

		// 检查读权限
		permChecker := NewConfigurationPermissionChecker(l.svcCtx.Casbin, l.Logger)
		if err := permChecker.CheckReadPermission(l.ctx, group); err != nil {
			return nil, err
		}

		return &core.ConfigurationInfo{
			Key:   in.Value,
			Value: value,
		}, nil
	}

	// 缓存未命中，从数据库查询
	config, err := l.svcCtx.DBEnt.Configuration.Query().
		Where(configuration.KeyEQ(in.Value)).
		First(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	// 检查读权限
	permChecker := NewConfigurationPermissionChecker(l.svcCtx.Casbin, l.Logger)
	if err := permChecker.CheckReadPermission(l.ctx, config.Group); err != nil {
		return nil, err
	}

	// 更新缓存
	v := fmt.Sprintf("%s<>%s", config.Group, config.Value)
	l.svcCtx.ConfigurationCache.Set(config.Key, v)

	return ConvertConfigurationToConfigurationInfo(config), nil
}

func splitGroupAndValue(s string) (string, string, error) {
	parts := strings.Split(s, "<>")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("配置项 '%s' 格式错误", s)
	}
	return parts[0], parts[1], nil
}
