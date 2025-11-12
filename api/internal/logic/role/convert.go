package role

import (
	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/client/roleservice"
)

// ConvertRpcRoleInfoToApiRoleInfo 将 RPC RoleInfo 转换为 API RoleInfo
func ConvertRpcRoleInfoToApiRoleInfo(rpcRole *roleservice.RoleInfo) *types.RoleInfo {
	if rpcRole == nil {
		return nil
	}

	return &types.RoleInfo{
		ID:          rpcRole.Id,
		CreatedAt:   rpcRole.CreatedAt,
		UpdatedAt:   rpcRole.UpdatedAt,
		RoleName:    pointer.GetString(rpcRole.RoleName),
		RoleCode:    pointer.GetString(rpcRole.RoleCode),
		Description: pointer.GetString(rpcRole.Description),
		State:       rpcRole.State,
	}
}

// ConvertApiRoleInfoToRpcRoleInfo 将 API RoleInfo 转换为 RPC RoleInfo
func ConvertApiRoleInfoToRpcRoleInfo(apiRole *types.RoleInfo) *roleservice.RoleInfo {
	if apiRole == nil {
		return nil
	}

	return &roleservice.RoleInfo{
		Id:          apiRole.ID,
		CreatedAt:   apiRole.CreatedAt,
		UpdatedAt:   apiRole.UpdatedAt,
		RoleName:    pointer.ToStringPtrIfNotEmpty(apiRole.RoleName),
		RoleCode:    pointer.ToStringPtrIfNotEmpty(apiRole.RoleCode),
		Description: pointer.ToStringPtrIfNotEmpty(apiRole.Description),
		State:       apiRole.State,
	}
}

