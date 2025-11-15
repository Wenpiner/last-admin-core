package menu

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMenuLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新菜单
func NewUpdateMenuLogic(r *http.Request, svcCtx *svc.ServiceContext) *UpdateMenuLogic {
	return &UpdateMenuLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *UpdateMenuLogic) UpdateMenu(req *types.MenuInfo) (resp *types.BaseResponse, err error) {
	_, err = l.svcCtx.MenuRpc.CreateOrUpdateMenu(l.ctx, &core.MenuInfo{
		Id:          req.ID,
		MenuCode:    &req.Name,
		MenuName:    &req.Name,
		ParentId:    req.ParentId,
		MenuPath:    &req.Path,
		State:       req.State,
		Sort:        req.Meta.Order,
		MenuType:    &req.Type,
		Description: req.Description,
		Component:   req.Component,
		Redirect:    req.Redirect,
		ServiceName: &req.Service,
		Permission:  req.Permission,
		Meta: &core.MenuMeta{
			Title:        &req.Meta.Title,
			Icon:         &req.Meta.Icon,
			IsHidden:     req.Meta.HideInMenu,
			IsBreadcrumb: req.Meta.HideInMenu,
			IsCache:      req.Meta.KeepAlive,
			IsTab:        req.Meta.AffixTab,
			IsAffix:      req.Meta.AffixTab,
			Link:         req.Meta.Link,
			FrameSrc:     req.Meta.IframeSrc,
		},
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
