package roleservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type ResetSuperApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetSuperApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetSuperApiLogic {
	return &ResetSuperApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 重置超级管理员API权限
func (l *ResetSuperApiLogic) ResetSuperApi(in *core.ID32Request) (*core.BaseResponse, error) {
	// 1. 查询角色
	roleEntity, err := l.svcCtx.DBEnt.Role.Get(l.ctx, in.Id)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	// 2. 检查角色编码是否为超级管理员 "super"
	if roleEntity.RoleCode != "super" {
		return nil, errorx.NewInvalidArgumentError("只能重置超级管理员角色")
	}

	// 3. 查询系统中所有的 API
	apis, err := l.svcCtx.DBEnt.API.Query().All(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	// 4. 为超级管理员分配所有 API
	var policies [][]string
	for _, apiItem := range apis {
		policies = append(policies, []string{roleEntity.RoleCode, apiItem.Path, apiItem.Method})
	}

	// 5. 查询并清理超级管理员旧策略
	oldPolicies, err := l.svcCtx.Casbin.GetFilteredPolicy(0, roleEntity.RoleCode)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	if len(oldPolicies) != 0 {
		removeResult, err := l.svcCtx.Casbin.RemoveFilteredPolicy(0, roleEntity.RoleCode)
		if err != nil {
			return nil, errorhandler.DBEntError(l.Logger, err, in)
		}
		if !removeResult {
			return nil, errorhandler.DBEntError(l.Logger, err, in)
		}
	}

	// 6. 添加新策略
	if len(policies) > 0 {
		if result, err := l.svcCtx.Casbin.AddPolicies(policies); err != nil || !result {
			return nil, errorhandler.DBEntError(l.Logger, err, in)
		}
	}

	return &core.BaseResponse{
		Message: "重置成功",
	}, nil
}
