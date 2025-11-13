package configurationservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-core/rpc/ent/configuration"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteConfigurationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteConfigurationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteConfigurationLogic {
	return &DeleteConfigurationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除配置
func (l *DeleteConfigurationLogic) DeleteConfiguration(in *core.StringRequest) (*core.BaseResponse, error) {
	// 查询配置是否存在
	config, err := l.svcCtx.DBEnt.Configuration.Query().
		Where(configuration.KeyEQ(in.Value)).
		First(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	// 检查写权限
	permChecker := NewConfigurationPermissionChecker(l.svcCtx.Casbin, l.Logger)
	if err := permChecker.CheckWritePermission(l.ctx, config.Group); err != nil {
		return nil, err
	}

	_, err = l.svcCtx.DBEnt.Configuration.Delete().
		Where(configuration.KeyEQ(in.Value)).
		Exec(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	// 从缓存删除
	l.svcCtx.ConfigurationCache.Delete(in.Value)

	return &core.BaseResponse{
		Message: "common.deleteSuccess",
	}, nil
}
