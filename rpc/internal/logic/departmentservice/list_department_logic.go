package departmentservicelogic

import (
	"context"

	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/ent/department"
	"github.com/wenpiner/last-admin-core/rpc/ent/predicate"
	"github.com/wenpiner/last-admin-core/rpc/internal/utils/errorhandler"

	"github.com/wenpiner/last-admin-core/rpc/internal/svc"
	"github.com/wenpiner/last-admin-core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDepartmentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDepartmentLogic {
	return &ListDepartmentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取部门列表
func (l *ListDepartmentLogic) ListDepartment(in *core.DepartmentListRequest) (*core.DepartmentListResponse, error) {
	// 构建查询条件
	var predicates []predicate.Department

	// 根据部门名称模糊搜索
	if in.DeptName != nil && *in.DeptName != "" {
		predicates = append(predicates, department.DeptNameContains(*in.DeptName))
	}

	// 根据部门编码模糊搜索
	if in.DeptCode != nil && *in.DeptCode != "" {
		predicates = append(predicates, department.DeptCodeContains(*in.DeptCode))
	}

	page, err := l.svcCtx.DBEnt.Department.Query().Where(predicates...).Page(l.ctx, in.Page.PageNumber, in.Page.PageSize)
	if err != nil {
		return nil, errorhandler.DBEntError(l.Logger, err, in)
	}
	resp := &core.DepartmentListResponse{
		Page: &core.BasePageResp{
			Total:      page.PageDetails.Total,
			PageNumber: in.Page.PageNumber,
			PageSize:   in.Page.PageSize,
		},
	}

	for _, dept := range page.List {
		resp.List = append(resp.List, l.convertDepartmentToDepartmentInfo(dept))
	}

	return resp, nil
}

// 将 Department 实体转换为 DepartmentInfo
func (l *ListDepartmentLogic) convertDepartmentToDepartmentInfo(dept *ent.Department) *core.DepartmentInfo {
	return &core.DepartmentInfo{
		Id:           &dept.ID,
		CreatedAt:    pointer.ToInt64Ptr(dept.CreatedAt.UnixMilli()),
		UpdatedAt:    pointer.ToInt64Ptr(dept.UpdatedAt.UnixMilli()),
		DeptName:     &dept.DeptName,
		DeptCode:     &dept.DeptCode,
		ParentId:     dept.ParentID,
		SortOrder:    pointer.ToInt32Ptr(dept.Sort),
		LeaderUserId: dept.LeaderUserID,
		State:        pointer.ToBoolPtr(dept.State),
		Description:  dept.Description,
	}
}