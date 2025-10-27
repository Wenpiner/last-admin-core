package menu

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllMenusLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取所有菜单
func NewGetAllMenusLogic(r *http.Request, svcCtx *svc.ServiceContext) *GetAllMenusLogic {
	return &GetAllMenusLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *GetAllMenusLogic) GetAllMenus() (resp *types.MenuListResponse, err error) {
	menus, err := l.svcCtx.MenuRpc.ListMenu(l.ctx, &core.MenuListRequest{
		Page: &core.BasePageRequest{
			PageNumber: 1,
			PageSize:   1000,
		},
	})
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
