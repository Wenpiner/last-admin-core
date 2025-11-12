package role

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/client/roleservice"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteRoleLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除角色
func NewDeleteRoleLogic(r *http.Request, svcCtx *svc.ServiceContext) *DeleteRoleLogic {
	return &DeleteRoleLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *DeleteRoleLogic) DeleteRole(req *types.ID32Request) (resp *types.BaseResponse, err error) {
	// 调用 RPC 服务进行删除
	rpcReq := &roleservice.ID32Request{Id: req.ID}
	rpcResp, err := l.svcCtx.RoleRpc.DeleteRole(l.ctx, rpcReq)
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
