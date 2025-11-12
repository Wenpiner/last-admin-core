package user

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/client/userservice"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除用户
func NewDeleteUserLogic(r *http.Request, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *DeleteUserLogic) DeleteUser(req *types.UUIDRequest) (resp *types.BaseResponse, err error) {
	// 调用 RPC 服务进行删除
	rpcReq := &userservice.UUIDRequest{Id: req.ID}
	rpcResp, err := l.svcCtx.UserRpc.DeleteUser(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 将 RPC 响应转换为 API 响应
	resp = &types.BaseResponse{
		Code:    0,
		Message: rpcResp.Message,
	}
	return resp, nil
}
