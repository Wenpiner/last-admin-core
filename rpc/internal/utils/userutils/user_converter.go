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
		for i, position := range user.Edges.Positions {
			positionIDs[i] = position.ID
		}
		userInfo.PositionIds = positionIDs
	}

	// 设置TOTP状态
	if len(user.Edges.Totp) > 0 {
		totp := user.Edges.Totp[0] // 一个用户只有一个TOTP记录
		userInfo.TotpEnabled = &totp.IsEnabled
		userInfo.TotpVerified = &totp.IsVerified
	} else {
		falseValue := false
		userInfo.TotpEnabled = &falseValue
		userInfo.TotpVerified = &falseValue
	}

	return userInfo
}
