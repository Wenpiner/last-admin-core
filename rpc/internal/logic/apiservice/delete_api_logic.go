package apiservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-core/rpc/ent/api"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteApiLogic {
	return &DeleteApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除API
func (l *DeleteApiLogic) DeleteApi(in *core.ID32SRequest) (*core.BaseResponse, error) {


	// 执行删除操作
	_, err := l.svcCtx.DBEnt.API.Delete().Where(
		api.IDIn(in.Ids...),
	).Exec(l.ctx)

	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	return &core.BaseResponse{
		Message: "删除成功",
	}, nil
}
