package roleservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-core/rpc/ent/api"
	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAssignApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignApiLogic {
	return &AssignApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 为角色分配API
func (l *AssignApiLogic) AssignApi(in *core.RoleApiRequest) (*core.BaseResponse, error) {
	// 查询角色
	role, err := l.svcCtx.DBEnt.Role.Get(l.ctx, *in.RoleId)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	// 查询API
	apis, err := l.svcCtx.DBEnt.API.Query().Where(api.IDIn(in.ApiIds...)).All(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	// 为角色分配API
	var policies [][]string
	for _, api := range apis {
		policies = append(policies, []string{role.RoleCode, api.Path, api.Method})
	}

	// 查询旧策略
	oldPolicies, err := l.svcCtx.Casbin.GetFilteredPolicy(0, role.RoleCode)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	if len(oldPolicies) != 0 {
		removeResult, err := l.svcCtx.Casbin.RemoveFilteredPolicy(0, role.RoleCode)
		if err != nil {
			return nil, errorhandler.DBEntError(l.Logger, err, in)
		}
		if !removeResult {
			return nil, errorhandler.DBEntError(l.Logger, err, in)
		}
	}
	// 添加新策略
	if result, err := l.svcCtx.Casbin.AddPolicies(policies); err != nil || !result {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}
	
	return &core.BaseResponse{}, nil
}
