package departmentservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-core/rpc/ent/department"
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
	dept, err := l.svcCtx.DBEnt.Department.Query().WithLeader().Where(department.ID(in.Id)).Only(l.ctx)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}

	return ConvertDepartmentToDepartmentInfo(dept), nil
}
