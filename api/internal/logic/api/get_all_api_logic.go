package api

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllApiLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取所有API
func NewGetAllApiLogic(r *http.Request, svcCtx *svc.ServiceContext) *GetAllApiLogic {
	return &GetAllApiLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *GetAllApiLogic) GetAllApi() (resp *types.ApiAllResponse, err error) {
	apis, err := l.svcCtx.ApiRpc.ListApi(l.ctx, &core.ApiListRequest{
		Page: &core.BasePageRequest{
			PageNumber: 1,
			PageSize:   5000,
		},
	})
	if err != nil {
		return nil, err
	}

	resp = &types.ApiAllResponse{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "success",
		},
		Data: ConvertToApiInfo(apis.List),
	}

	return
}
