package roleservicelogic

import (
	"context"
	"fmt"

	"github.com/wenpiner/last-admin-core/rpc/ent/configuration"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConfigurationGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConfigurationGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConfigurationGroupLogic {
	return &GetConfigurationGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取角色配置项分组权限
func (l *GetConfigurationGroupLogic) GetConfigurationGroup(in *core.StringRequest) (*core.RoleConfigurationGroupListResponse, error) {

	if in.Value == "" {
		// 查询当前所有分组列表
		var list []struct{
			Group string `json:"group,omitempty"`
		}
		err := l.svcCtx.DBEnt.Configuration.Query().GroupBy(configuration.FieldGroup).Scan(l.ctx,&list)
		if err != nil {
			return nil, err
		}
		var resp []string
		for _, v := range list {
			resp = append(resp, v.Group)
		}
		return &core.RoleConfigurationGroupListResponse{
			List: resp,
		}, nil
	}

	policies, err := l.svcCtx.Casbin.GetFilteredPolicy(0, in.Value, "configuration")
	if err != nil {
		return nil, err
	}
	var configurationGroups []string
	for _, v := range policies {
		configurationGroups = append(configurationGroups, fmt.Sprintf("%s:%s", v[2], v[3]))
	}
	return &core.RoleConfigurationGroupListResponse{
		List: configurationGroups,
	}, nil
}
