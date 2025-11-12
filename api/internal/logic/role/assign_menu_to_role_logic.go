package role

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignMenuToRoleLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 为角色分配菜单
func NewAssignMenuToRoleLogic(r *http.Request, svcCtx *svc.ServiceContext) *AssignMenuToRoleLogic {
	return &AssignMenuToRoleLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *AssignMenuToRoleLogic) AssignMenuToRole(req *types.RoleMenuRequest) (resp *types.BaseResponse, err error) {
	_, err = l.svcCtx.RoleRpc.AssignMenu(l.ctx, &core.RoleMenuRequest{
		RoleId: &req.RoleId,
		MenuIds: req.MenuIds,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.BaseResponse{
		Code:    0,
		Message: "success",
	}

	return
}
