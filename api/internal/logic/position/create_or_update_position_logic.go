package position

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrUpdatePositionLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建或更新岗位
func NewCreateOrUpdatePositionLogic(r *http.Request, svcCtx *svc.ServiceContext) *CreateOrUpdatePositionLogic {
	return &CreateOrUpdatePositionLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *CreateOrUpdatePositionLogic) CreateOrUpdatePosition(req *types.PositionInfo) (resp *types.ModifyPositionResponse, err error) {
	// 调用 RPC 服务进行创建或更新
	rpcReq := ConvertApiPositionInfoToRpcPositionInfo(req)
	rpcResp, err := l.svcCtx.PositionRpc.CreateOrUpdatePosition(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 将 RPC 响应转换为 API 响应
	position := ConvertRpcPositionInfoToApiPositionInfo(rpcResp)
	resp = &types.ModifyPositionResponse{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "success",
		},
		Data: *position,
	}
	return resp, nil
}
