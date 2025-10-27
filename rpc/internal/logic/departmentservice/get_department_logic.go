package departmentservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDepartmentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDepartmentLogic {
	return &GetDepartmentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取部门
func (l *GetDepartmentLogic) GetDepartment(in *core.ID32Request) (*core.DepartmentInfo, error) {
	dept, err := l.svcCtx.DBEnt.Department.Get(l.ctx, in.Id)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	return l.convertDepartmentToDepartmentInfo(dept), nil
}

// 将 Department 实体转换为 DepartmentInfo
func (l *GetDepartmentLogic) convertDepartmentToDepartmentInfo(dept *ent.Department) *core.DepartmentInfo {
	return &core.DepartmentInfo{
		Id:           &dept.ID,
		CreatedAt:    pointer.ToInt64Ptr(dept.CreatedAt.UnixMilli()),
		UpdatedAt:    pointer.ToInt64Ptr(dept.UpdatedAt.UnixMilli()),
		DeptName:     &dept.DeptName,
		DeptCode:     &dept.DeptCode,
		ParentId:     dept.ParentID,
		SortOrder:    &dept.Sort,
		LeaderUserId: dept.LeaderUserID,
		State:        &dept.State,
		Description:  dept.Description,
	}
}
