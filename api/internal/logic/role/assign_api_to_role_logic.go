package role

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignApiToRoleLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 为角色分配API
func NewAssignApiToRoleLogic(r *http.Request, svcCtx *svc.ServiceContext) *AssignApiToRoleLogic {
	return &AssignApiToRoleLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *AssignApiToRoleLogic) AssignApiToRole(req *types.RoleApiRequest) (resp *types.BaseResponse, err error) {
	rpcReq := &core.RoleApiRequest{
		RoleId: &req.RoleId,
		ApiIds: req.ApiIds,
	}
	_, err = l.svcCtx.RoleRpc.AssignApi(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.BaseResponse{
		Code:    0,
		Message: "success",
	}

	return
}
