package role

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignConfigurationGroupToRoleLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 为角色分配配置项分组权限
func NewAssignConfigurationGroupToRoleLogic(r *http.Request, svcCtx *svc.ServiceContext) *AssignConfigurationGroupToRoleLogic {
	return &AssignConfigurationGroupToRoleLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *AssignConfigurationGroupToRoleLogic) AssignConfigurationGroupToRole(req *types.RoleConfigurationGroupRequest) (resp *types.RoleConfigurationGroupListResponse, err error) {
	response, err := l.svcCtx.RoleRpc.AssignConfigurationGroup(l.ctx, &core.RoleConfigurationGroupRequest{
		RoleValue:           req.RoleValue,
		ConfigurationGroups: req.ConfigurationGroups,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.RoleConfigurationGroupListResponse{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "success",
		},
		Data: response.List,
	}

	return
}
