package roleservicelogic

import (
	"context"
	"strings"

	"github.com/wenpiner/last-admin-core/rpc/ent/role"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type AssignConfigurationGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAssignConfigurationGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignConfigurationGroupLogic {
	return &AssignConfigurationGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 为角色分配配置项分组权限
func (l *AssignConfigurationGroupLogic) AssignConfigurationGroup(in *core.RoleConfigurationGroupRequest) (*core.RoleConfigurationGroupListResponse, error) {
	// 验证角色是否存在
	role, err := l.svcCtx.DBEnt.Role.Query().Where(role.RoleCodeEQ(in.RoleValue)).First(l.ctx)
	if err != nil {
		return nil, err
	}
	var configurationGroups [][]string
	var resp []string
	for _, v := range in.ConfigurationGroups {
		result := strings.Split(v, ":")
		if len(result) != 2 {
			continue
		}
		resp = append(resp, v)
		configurationGroups = append(configurationGroups, []string{role.RoleCode, "configuration", result[0], result[1]})
	}
	if len(configurationGroups) == 0 {
		return nil,errorx.NewInvalidArgumentError("common.invalidArgument")
	}
	b, err := l.svcCtx.Casbin.UpdateFilteredPolicies(configurationGroups, 0, role.RoleCode, "configuration")
	if err != nil {
		logx.Errorw("更新 Casbin 策略失败", logx.Field("detail", err.Error()))
		return nil, err
	}

	if !b {
		logx.Errorw("更新 Casbin 策略失败", logx.Field("detail", "更新失败"))
		return nil, err
	}
	return &core.RoleConfigurationGroupListResponse{
		List: resp,
	}, nil
}
