package departmentservicelogic

import (
	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/types/core"
)

// ConvertDepartmentToDepartmentInfo 将 Department 实体转换为 DepartmentInfo
func ConvertDepartmentToDepartmentInfo(dept *ent.Department) *core.DepartmentInfo {
	info := core.DepartmentInfo{
		Id:          &dept.ID,
		CreatedAt:   pointer.ToInt64Ptr(dept.CreatedAt.UnixMilli()),
		UpdatedAt:   pointer.ToInt64Ptr(dept.UpdatedAt.UnixMilli()),
		DeptName:    &dept.DeptName,
		DeptCode:    &dept.DeptCode,
		ParentId:    dept.ParentID,
		SortOrder:   &dept.Sort,
		State:       &dept.State,
		Description: dept.Description,
	}

	if dept.LeaderUserID != nil {
		info.LeaderUserId = pointer.ToStringPtrIfNotEmpty(dept.LeaderUserID.String())
	}

	if dept.Edges.Leader != nil {
		info.LaderUsername = pointer.ToStringPtrIfNotEmpty(dept.Edges.Leader.Username)
		info.LaderPhone = pointer.ToStringPtrIfNotEmpty(dept.Edges.Leader.Mobile)
		info.LaderEmail = pointer.ToStringPtrIfNotEmpty(dept.Edges.Leader.Email)
	}

	return &info
}
