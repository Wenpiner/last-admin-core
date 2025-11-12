package role

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/logic/menu"
	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleMenuLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取角色菜单
func NewGetRoleMenuLogic(r *http.Request, svcCtx *svc.ServiceContext) *GetRoleMenuLogic {
	return &GetRoleMenuLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *GetRoleMenuLogic) GetRoleMenu(req *types.StringIDRequest) (resp *types.MenuListResponse, err error) {
	menus, err := l.svcCtx.MenuRpc.ListMenuByRole(l.ctx, &core.StringRequest{Value: req.ID})
	if err != nil {
		return nil, err
	}
	// 处理所有Menu,按照层级关系进行组装
	resp = &types.MenuListResponse{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "success",
		},
		Data: menu.ConvertToMenuInfo(menus.List),
	}

	return
}
