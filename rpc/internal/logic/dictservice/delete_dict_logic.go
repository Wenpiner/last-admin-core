package dictservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteDictLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDictLogic {
	return &DeleteDictLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除字典
func (l *DeleteDictLogic) DeleteDict(in *core.ID32Request) (*core.BaseResponse, error) {
	err := l.svcCtx.DBEnt.DictType.DeleteOneID(in.Id).Exec(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	return &core.BaseResponse{
		Message: "common.deleteSuccess",
	}, nil
}
