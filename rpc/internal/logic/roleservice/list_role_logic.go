package roleservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-core/rpc/ent/predicate"
	"github.com/wenpiner/last-admin-core/rpc/ent/role"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRoleLogic {
	return &ListRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取角色列表
func (l *ListRoleLogic) ListRole(in *core.RoleListRequest) (*core.RoleListResponse, error) {
	// 构建查询条件
	var predicates []predicate.Role

	// 根据角色名称模糊搜索
	if in.RoleName != nil && *in.RoleName != "" {
		predicates = append(predicates, role.RoleNameContains(*in.RoleName))
	}

	// 根据角色编码模糊搜索
	if in.RoleCode != nil && *in.RoleCode != "" {
		predicates = append(predicates, role.RoleCodeContains(*in.RoleCode))
	}

	page, err := l.svcCtx.DBEnt.Role.Query().Where(predicates...).Page(l.ctx, in.Page.PageNumber, in.Page.PageSize)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}
	resp := &core.RoleListResponse{
		Page: &core.BasePageResp{
			Total:      page.PageDetails.Total,
			PageNumber: in.Page.PageNumber,
			PageSize:   in.Page.PageSize,
		},
	}

	for _, role := range page.List {
		resp.List = append(resp.List, ConvertRoleToRoleInfo(role))
	}

	return resp, nil
}
