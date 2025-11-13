package role

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleConfigurationGroupLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取角色配置项分组权限
func NewGetRoleConfigurationGroupLogic(r *http.Request, svcCtx *svc.ServiceContext) *GetRoleConfigurationGroupLogic {
	return &GetRoleConfigurationGroupLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *GetRoleConfigurationGroupLogic) GetRoleConfigurationGroup(req *types.StringIDRequest) (resp *types.RoleConfigurationGroupListResponse, err error) {
	response, err := l.svcCtx.RoleRpc.GetConfigurationGroup(l.ctx, &core.StringRequest{
		Value: req.ID,
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
