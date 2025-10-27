package dictservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteDictItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteDictItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDictItemLogic {
	return &DeleteDictItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除字典子项
func (l *DeleteDictItemLogic) DeleteDictItem(in *core.ID32Request) (*core.BaseResponse, error) {
	err := l.svcCtx.DBEnt.DictItem.DeleteOneID(in.Id).Exec(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	return &core.BaseResponse{
		Message: "common.deleteSuccess",
	}, nil
}
