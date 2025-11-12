package position

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPositionLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取岗位列表
func NewListPositionLogic(r *http.Request, svcCtx *svc.ServiceContext) *ListPositionLogic {
	return &ListPositionLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *ListPositionLogic) ListPosition(req *types.PositionListRequest) (resp *types.PositionListResponse, err error) {
	positionList, err := l.svcCtx.PositionRpc.ListPosition(l.ctx, &core.PositionListRequest{
		Page: &core.BasePageRequest{
			PageNumber: req.Page.CurrentPage,
			PageSize:   req.Page.PageSize,
		},
		PositionName: &req.PositionName,
		PositionCode: &req.PositionCode,
	})

	if err != nil {
		return nil, err
	}

	resp = ConvertRpcPositionListResponseToApiPositionListResponse(positionList)

	return
}
