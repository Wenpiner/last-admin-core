package position

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePositionLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除岗位
func NewDeletePositionLogic(r *http.Request, svcCtx *svc.ServiceContext) *DeletePositionLogic {
	return &DeletePositionLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *DeletePositionLogic) DeletePosition(req *types.ID32Request) (resp *types.BaseResponse, err error) {
	_, err = l.svcCtx.PositionRpc.DeletePosition(l.ctx, &core.ID32Request{Id: req.ID})
	if err != nil {
		return nil, err
	}

	resp = &types.BaseResponse{
		Code:    0,
		Message: "success",
	}

	return
}
