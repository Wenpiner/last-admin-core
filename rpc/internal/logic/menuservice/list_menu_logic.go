package menuservicelogic

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/wenpiner/last-admin-core/rpc/ent/menu"
	"github.com/wenpiner/last-admin-core/rpc/ent/predicate"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMenuLogic {
	return &ListMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取菜单列表
func (l *ListMenuLogic) ListMenu(in *core.MenuListRequest) (*core.MenuListResponse, error) {
	var predicates []predicate.Menu
	if in.MenuName != nil {
		predicates = append(predicates, menu.MenuNameContains(*in.MenuName))
	}
	if in.MenuCode != nil {
		predicates = append(predicates, menu.MenuCodeContains(*in.MenuCode))
	}

	page, err := l.svcCtx.DBEnt.Menu.Query().Where(predicates...).Order(menu.ByMenuLevel(), menu.BySort(sql.OrderDesc())).Page(l.ctx, in.Page.PageNumber, in.Page.PageSize)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}
	resp := &core.MenuListResponse{
		Page: &core.BasePageResp{
			Total:      page.PageDetails.Total,
			PageNumber: page.PageDetails.Page,
			PageSize:   page.PageDetails.Pages,
		},
	}
	for _, v := range page.List {
		resp.List = append(resp.List, ConvertMenuToMenuInfo(v))
	}
	return resp, nil
}
