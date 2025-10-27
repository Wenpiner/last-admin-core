package menu

import (
	"context"

	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllMenusByRoleLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户角色当前所有菜单
func NewGetAllMenusByRoleLogic(r *http.Request, svcCtx *svc.ServiceContext) *GetAllMenusByRoleLogic {
	return &GetAllMenusByRoleLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *GetAllMenusByRoleLogic) GetAllMenusByRole() (resp *types.MenuListResponse, err error) {
	menus, err := l.svcCtx.MenuRpc.ListMenuByRole(l.ctx, &core.StringRequest{Value: l.ctx.Value("roleId").(string)})
	if err != nil {
		return nil, err
	}
	// 处理所有Menu,按照层级关系进行组装
	resp = &types.MenuListResponse{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "success",
		},
		Data: ConvertToMenuInfo(menus.List),
	}

	return
}

func ConvertToMenuInfo(menus []*core.MenuInfo) []types.MenuInfo {

	var rMenus []types.MenuInfo

	for _, m := range menus {
		var meta types.Meta
		if m.Meta != nil {
			meta = types.Meta{
				Title:         pointer.GetString(m.Meta.Title),
				Icon:          pointer.GetString(m.Meta.Icon),
				Order:         m.Sort,
				HideInMenu:    m.Meta.IsHidden,
				AffixTab:      m.Meta.IsAffix,
				AffixTabOrder: pointer.ToInt32Ptr(0),
				Link:          m.Meta.Link,
				IframeSrc:     m.Meta.FrameSrc,
				KeepAlive:     m.Meta.IsCache,
			}
		}
		menu := types.MenuInfo{
			Path:      pointer.GetString(m.MenuPath),
			Name:      pointer.GetString(m.MenuCode),
			Redirect:  m.Redirect,
			Component: m.Component,
			Meta:      meta,
			ParentId:  m.ParentId,
			State:     m.State,
			ID:        *m.Id,
			Service:   pointer.GetString(m.ServiceName),
			Permission: m.Permission,
			CreatedAt:  m.CreatedAt,
		}

		rMenus = append(rMenus, menu)
	}
	return rMenus
}
