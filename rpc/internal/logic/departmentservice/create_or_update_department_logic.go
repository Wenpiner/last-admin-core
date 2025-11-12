package departmentservicelogic

import (
	"context"

	"github.com/google/uuid"
	last_i18n "github.com/wenpiner/last-admin-common/last-i18n"
	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrUpdateDepartmentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrUpdateDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrUpdateDepartmentLogic {
	return &CreateOrUpdateDepartmentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建或更新部门
func (l *CreateOrUpdateDepartmentLogic) CreateOrUpdateDepartment(in *core.DepartmentInfo) (*core.DepartmentInfo, error) {
	// 开启事务，并先进行检查是否存在，如果存在则进行更新否则进行创建
	tx, err := l.svcCtx.DBEnt.Tx(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	var leaderUserID uuid.UUID
	if in.LeaderUserId != nil {
		leaderUserID, err = uuid.Parse(*in.LeaderUserId)
		if err != nil {
			return nil, errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
		}
	}

	defer tx.Rollback()

	var department *ent.Department

	if in.Id != nil && *in.Id != 0 {
		// 更新操作
		if err := l.validateUpdate(in); err != nil {
			return nil, err
		}

		updateQuery := tx.Department.UpdateOneID(*in.Id)

		// 设置可更新的字段
		if in.DeptName != nil {
			updateQuery.SetDeptName(*in.DeptName)
		}
		if in.DeptCode != nil {
			updateQuery.SetDeptCode(*in.DeptCode)
		}
		if in.ParentId != nil {
			updateQuery.SetNillableParentID(in.ParentId)
		}
		if in.SortOrder != nil {
			updateQuery.SetNillableSort(in.SortOrder)
		}
		if in.LeaderUserId != nil {
			updateQuery.SetNillableLeaderUserID(&leaderUserID)
		}
		if in.State != nil {
			updateQuery.SetNillableState(in.State)
		}
		if in.Description != nil {
			updateQuery.SetNillableDescription(in.Description)
		}

		department, err = updateQuery.Save(l.ctx)
	} else {
		// 创建操作
		if err := l.validateCreate(in); err != nil {
			return nil, err
		}

		createQuery := tx.Department.Create().
			SetDeptName(*in.DeptName).
			SetDeptCode(*in.DeptCode).
			SetNillableParentID(in.ParentId).
			SetNillableSort(in.SortOrder).
			SetNillableLeaderUserID(&leaderUserID).
			SetState(pointer.GetBool(in.State)).
			SetNillableDescription(in.Description)

		department, err = createQuery.Save(l.ctx)
	}

	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	err = tx.Commit()
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	return l.convertDepartmentToDepartmentInfo(department), nil
}

// 验证新增参数可用性
func (l *CreateOrUpdateDepartmentLogic) validateCreate(in *core.DepartmentInfo) error {
	if in.DeptName == nil || in.DeptCode == nil {
		return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}
	return nil
}

// 验证更新参数可用性
func (l *CreateOrUpdateDepartmentLogic) validateUpdate(in *core.DepartmentInfo) error {
	if in.Id == nil {
		return errorx.NewInvalidArgumentError(last_i18n.ValidationFailed)
	}
	return nil
}

// 将 Department 实体转换为 DepartmentInfo
func (l *CreateOrUpdateDepartmentLogic) convertDepartmentToDepartmentInfo(dept *ent.Department) *core.DepartmentInfo {
	return &core.DepartmentInfo{
		Id:           &dept.ID,
		CreatedAt:    pointer.ToInt64Ptr(dept.CreatedAt.UnixMilli()),
		UpdatedAt:    pointer.ToInt64Ptr(dept.UpdatedAt.UnixMilli()),
		DeptName:     &dept.DeptName,
		DeptCode:     &dept.DeptCode,
		ParentId:     dept.ParentID,
		SortOrder:    &dept.Sort,
		LeaderUserId: pointer.ToStringPtrIfNotEmpty(dept.LeaderUserID.String()),
		State:        pointer.ToBoolPtr(dept.State),
		Description:  dept.Description,
	}
}

