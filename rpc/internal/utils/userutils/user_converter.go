package userutils

import (
	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/rpc/ent"
	"github.com/wenpiner/last-admin-core/rpc/types/core"
)

// ConvertUserToUserInfo 将 User 实体转换为 UserInfo
// 这是一个共享的转换函数，用于避免在多个 logic 文件中重复相同的代码
func ConvertUserToUserInfo(user *ent.User) *core.UserInfo {

	userInfo := &core.UserInfo{
		Id:              pointer.ToStringPtr(user.ID.String()),
		CreatedAt:       pointer.ToInt64Ptr(user.CreatedAt.UnixMilli()),
		UpdatedAt:       pointer.ToInt64Ptr(user.UpdatedAt.UnixMilli()),
		Username:        &user.Username,
		Email:           pointer.ToStringPtrIfNotEmpty(user.Email),
		FullName:        pointer.ToStringPtrIfNotEmpty(user.FullName),
		Mobile:          pointer.ToStringPtrIfNotEmpty(user.Mobile),
		Avatar:          pointer.ToStringPtrIfNotEmpty(user.Avatar),
		UserDescription: pointer.ToStringPtrIfNotEmpty(user.UserDescription),
		LastLoginIp:     pointer.ToStringPtrIfNotEmpty(user.LastLoginIP),
		State:           &user.State,
		DepartmentId:    pointer.ToUint32PtrIfNotZero(user.DepartmentID),
		PasswordHash:    pointer.ToStringPtrIfNotEmpty(user.PasswordHash),
		HomePath:        pointer.ToStringPtrIfNotEmpty(user.HomePath),
	}

	// 设置最后登录时间
	if user.LastLoginAt != nil {
		userInfo.LastLoginAt = pointer.ToInt64Ptr(user.LastLoginAt.UnixMilli())
	}

	// 提取角色ID和角色值
	if user.Edges.Roles != nil {
		roleIDs := make([]uint32, len(user.Edges.Roles))
		roleValues := make([]string, len(user.Edges.Roles))
		roleNames := make([]string, len(user.Edges.Roles))
		for i, role := range user.Edges.Roles {
			roleIDs[i] = role.ID
			roleValues[i] = role.RoleCode
			roleNames[i] = role.RoleName
		}
		userInfo.RoleValues = roleValues
		userInfo.RoleIds = roleIDs
		userInfo.RoleNames = roleNames
	}

	// 提取部门信息
	if user.Edges.Department != nil {
		userInfo.DepartmentName = pointer.ToStringPtrIfNotEmpty(user.Edges.Department.DeptName)
	}

	// 提取职位ID
	if user.Edges.Positions != nil {
		positionIDs := make([]uint32, len(user.Edges.Positions))
		positionNames := make([]string, len(user.Edges.Positions))
		for i, position := range user.Edges.Positions {
			positionIDs[i] = position.ID
			positionNames[i] = position.PositionName
		}
		userInfo.PositionIds = positionIDs
		userInfo.PositionNames = positionNames
	}

	// 设置TOTP信息
	if user.Edges.Totp != nil {
		totpInfo := &core.TotpInfo{
			Id:           pointer.ToStringPtr(user.Edges.Totp.ID.String()),
			CreatedAt:    pointer.ToInt64Ptr(user.Edges.Totp.CreatedAt.UnixMilli()),
			UpdatedAt:    pointer.ToInt64Ptr(user.Edges.Totp.UpdatedAt.UnixMilli()),
			State:        pointer.ToBoolPtr(user.Edges.Totp.State),
			IsVerified:   pointer.ToBoolPtr(user.Edges.Totp.IsVerified),
			DeviceName:   user.Edges.Totp.DeviceName,
			Issuer:       user.Edges.Totp.Issuer,
		}
		if user.Edges.Totp.LastUsedAt != nil {
			totpInfo.LastUsedAt = pointer.ToInt64Ptr(user.Edges.Totp.LastUsedAt.UnixMilli())
		}
		if user.Edges.Totp.LastUsedCode != nil {
			totpInfo.LastUsedCode = user.Edges.Totp.LastUsedCode
		}
		userInfo.TotpInfo = totpInfo
	}

	return userInfo
}
