package api

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteApiLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除API
func NewDeleteApiLogic(r *http.Request, svcCtx *svc.ServiceContext) *DeleteApiLogic {
	return &DeleteApiLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *DeleteApiLogic) DeleteApi(req *types.ID32Request) (resp *types.BaseResponse, err error) {
	_, err = l.svcCtx.ApiRpc.DeleteApi(l.ctx, &core.ID32SRequest{Ids: []uint32{req.ID}})
	if err != nil {
		return nil, err
	}

	resp = &types.BaseResponse{
		Code:    0,
		Message: "success",
	}

	return
}
