package menuservicelogic

import (
	"context"
	"strings"

	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/ent/menu"
	"github.com/wenpiner/last-admin-core/rpc/ent/role"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPagePermissionByRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPagePermissionByRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPagePermissionByRoleLogic {
	return &ListPagePermissionByRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取角色对应的所有页面权限
func (l *ListPagePermissionByRoleLogic) ListPagePermissionByRole(in *core.StringRequest) (*core.StringListResponse, error) {
	// 解析角色编码，通过逗号分隔
	roleCodes := strings.Split(in.Value, ",")
	if len(roleCodes) == 0 || (len(roleCodes) == 1 && roleCodes[0] == "") {
		// 如果没有角色编码，返回空列表
		return &core.StringListResponse{
			List: []string{},
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
		return &core.StringListResponse{
			List: []string{},
		}, nil
	}

	// 先查询角色，然后通过角色的 Edge.Menus 进行查询
	roles, err := l.svcCtx.DBEnt.Role.Query().
		Where(role.RoleCodeIn(validRoleCodes...)).
		WithMenus(func(q *ent.MenuQuery) {
			q.Select(menu.FieldPermission).
				Order(ent.Asc(menu.FieldSort), ent.Asc(menu.FieldID)). // 按排序字段和ID排序
				Where(menu.State(true), menu.PermissionNotNil())
		}).
		All(l.ctx)

	if err != nil {
		return nil, err
	}

	// 使用 Map 对 Permission 进行去重处理
	permissionSet := make(map[string]bool)
	var permissions []string

	for _, r := range roles {
		for _, m := range r.Edges.Menus {
			if m.Permission != nil && *m.Permission != "" {
				permission := *m.Permission
				if !permissionSet[permission] {
					permissionSet[permission] = true
					permissions = append(permissions, permission)
				}
			}
		}
	}

	return &core.StringListResponse{
		List: permissions,
	}, nil
}
