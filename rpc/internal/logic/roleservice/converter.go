package roleservicelogic

import (
	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/types/core"
)

// ConvertRoleToRoleInfo 将 Role 实体转换为 RoleInfo
func ConvertRoleToRoleInfo(role *ent.Role) *core.RoleInfo {
	return &core.RoleInfo{
		Id:          &role.ID,
		CreatedAt:   pointer.ToInt64Ptr(role.CreatedAt.UnixMilli()),
		UpdatedAt:   pointer.ToInt64Ptr(role.UpdatedAt.UnixMilli()),
		RoleName:    &role.RoleName,
		RoleCode:    &role.RoleCode,
		Description: role.Description,
		State:       &role.State,
	}
}

