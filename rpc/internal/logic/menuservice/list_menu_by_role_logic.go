package menuservicelogic

import (
	"context"
	"strings"

	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/ent/menu"
	"github.com/wenpiner/last-admin-core/rpc/ent/role"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMenuByRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListMenuByRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMenuByRoleLogic {
	return &ListMenuByRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 通过角色查询菜单
func (l *ListMenuByRoleLogic) ListMenuByRole(in *core.StringRequest) (*core.MenuListResponse, error) {
	// 解析角色编码，通过逗号分隔
	roleCodes := strings.Split(in.Value, ",")
	if len(roleCodes) == 0 || (len(roleCodes) == 1 && roleCodes[0] == "") {
		// 如果没有角色编码，返回空列表
		return &core.MenuListResponse{
			Page: &core.BasePageResp{
				Total:      0,
				PageNumber: 1,
				PageSize:   0,
			},
			List: []*core.MenuInfo{},
		}, nil
	}

	// 去除空字符串和前后空格
	var validRoleCodes []string
	for _, code := range roleCodes {
		trimmed := strings.TrimSpace(code)
		if trimmed != "" {
			validRoleCodes = append(validRoleCodes, trimmed)
		}
	}

	if len(validRoleCodes) == 0 {
		// 如果没有有效的角色编码，返回空列表
		return &core.MenuListResponse{
			Page: &core.BasePageResp{
				Total:      0,
				PageNumber: 1,
				PageSize:   0,
			},
			List: []*core.MenuInfo{},
		}, nil
	}

	// 先查询角色，然后通过角色的 Edge.Menus 进行查询
	roles, err := l.svcCtx.DBEnt.Role.Query().
		Where(role.RoleCodeIn(validRoleCodes...)).
		WithMenus(func(q *ent.MenuQuery) {
			q.Order(ent.Asc(menu.FieldMenuLevel), ent.Asc(menu.FieldSort)) // 按排序字段和ID排序
			q.Where(menu.State(true))
		}).
		All(l.ctx)

	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	// 从角色中提取所有菜单，同时进行去重处理
	menuSet := make(map[uint32]bool)
	var menus []*ent.Menu
	// 构建响应
	resp := &core.MenuListResponse{
		Page: &core.BasePageResp{
			Total:      uint64(len(menus)),
			PageNumber: 1,
			PageSize:   uint32(len(menus)),
		},
	}

	for _, r := range roles {
		if r.Edges.Menus != nil {
			for _, menu := range r.Edges.Menus {
				if !menuSet[menu.ID] {
					menuSet[menu.ID] = true
					resp.List = append(resp.List, ConvertMenuToMenuInfo(menu))
				}
			}
		}
	}

	resp.Page.Total = uint64(len(resp.List))
	resp.Page.PageSize = uint32(len(resp.List))

	return resp, nil
}
