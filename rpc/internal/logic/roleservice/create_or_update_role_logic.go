package roleservicelogic

import (
	"context"

	last_i18n "github.com/wenpiner/last-admin-common/last-i18n"
	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrUpdateRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrUpdateRoleLogic {
	return &CreateOrUpdateRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建或更新角色
func (l *CreateOrUpdateRoleLogic) CreateOrUpdateRole(in *core.RoleInfo) (*core.RoleInfo, error) {
	// 开启事务，并先进行检查是否存在，如果存在则进行更新否则进行创建
	tx, err := l.svcCtx.DBEnt.Tx(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}
	defer tx.Rollback()

	var role *ent.Role

	if in.Id != nil && *in.Id != 0 {
		// 更新操作
		if err := l.validateUpdate(in); err != nil {
			return nil, err
		}

		updateQuery := tx.Role.UpdateOneID(*in.Id)

		// 设置可更新的字段
		if in.RoleName != nil {
			updateQuery.SetRoleName(*in.RoleName)
		}
		if in.RoleCode != nil {
			updateQuery.SetRoleCode(*in.RoleCode)
		}
		if in.Description != nil {
			updateQuery.SetNillableDescription(in.Description)
		}
		if in.State != nil {
			updateQuery.SetState(*in.State)
		}

		role, err = updateQuery.Save(l.ctx)
	} else {
		// 创建操作
		if err := l.validateCreate(in); err != nil {
			return nil, err
		}

		createQuery := tx.Role.Create().
			SetRoleName(*in.RoleName).
			SetRoleCode(*in.RoleCode).
			SetNillableDescription(in.Description).
			SetState(l.getStateValue(in.State))

		role, err = createQuery.Save(l.ctx)
	}

	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	err = tx.Commit()
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	return ConvertRoleToRoleInfo(role), nil
}

// 验证新增参数可用性
func (l *CreateOrUpdateRoleLogic) validateCreate(in *core.RoleInfo) error {
	if in.RoleName == nil || in.RoleCode == nil {
		return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}
	return nil
}

// 验证更新参数可用性
func (l *CreateOrUpdateRoleLogic) validateUpdate(in *core.RoleInfo) error {
	if in.Id == nil {
		return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}
	return nil
}

// 获取状态值，默认为 true
func (l *CreateOrUpdateRoleLogic) getStateValue(state *bool) bool {
	if state != nil {
		return *state
	}
	return true
}

// 获取系统角色值，默认为 false
func (l *CreateOrUpdateRoleLogic) getIsSystemValue(isSystem *bool) bool {
	if isSystem != nil {
		return *isSystem
	}
	return false
}
