package user

import (
	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/client/userservice"
)

// ConvertRpcUserInfoToApiUserInfo 将 RPC UserInfo 转换为 API UserInfo
func ConvertRpcUserInfoToApiUserInfo(rpcUser *userservice.UserInfo) *types.UserInfo {
	if rpcUser == nil {
		return nil
	}

	userInfo := &types.UserInfo{
		Avatar:         pointer.GetString(rpcUser.Avatar),
		RealName:       pointer.GetString(rpcUser.FullName),
		Roles:          rpcUser.RoleValues,
		UserId:         pointer.GetString(rpcUser.Id),
		Username:       pointer.GetString(rpcUser.Username),
		Desc:           pointer.GetString(rpcUser.UserDescription),
		HomePath:       pointer.GetString(rpcUser.HomePath),
		Email:          pointer.GetString(rpcUser.Email),
		RoleNames:      rpcUser.RoleNames,
		DepartmentName: pointer.GetString(rpcUser.DepartmentName),
		Mobile:         pointer.GetString(rpcUser.Mobile),
		DepartmentId:   pointer.GetUint32(rpcUser.DepartmentId),
		PositionNames:  rpcUser.PositionNames,
		PositionIds:    rpcUser.PositionIds,
		State:          pointer.GetBool(rpcUser.State),
		CreatedAt:      pointer.GetInt64(rpcUser.CreatedAt),
		UpdatedAt:      pointer.GetInt64(rpcUser.UpdatedAt),
		LastLoginAt:    pointer.GetInt64(rpcUser.LastLoginAt),
		LastLoginIp:    pointer.GetString(rpcUser.LastLoginIp),
		RoleIds:        rpcUser.RoleIds,
	}
	if rpcUser.TotpInfo != nil {
		userInfo.TotpInfo = &types.TotpInfo{
			Id:           rpcUser.TotpInfo.Id,
			CreatedAt:    rpcUser.TotpInfo.CreatedAt,
			UpdatedAt:    rpcUser.TotpInfo.UpdatedAt,
			State:        rpcUser.TotpInfo.State,
			IsVerified:   rpcUser.TotpInfo.IsVerified,
			LastUsedAt:   rpcUser.TotpInfo.LastUsedAt,
			LastUsedCode: rpcUser.TotpInfo.LastUsedCode,
			DeviceName:   rpcUser.TotpInfo.DeviceName,
			Issuer:       rpcUser.TotpInfo.Issuer,
		}
	}

	return userInfo
}

// ConvertApiUserInfoToRpcUserInfo 将 API UserInfo 转换为 RPC UserInfo
func ConvertApiUserInfoToRpcUserInfo(apiUser *types.UserInfo) *userservice.UserInfo {
	if apiUser == nil {
		return nil
	}

	return &userservice.UserInfo{
		Id:              pointer.ToStringPtrIfNotEmpty(apiUser.UserId),
		CreatedAt:       pointer.ToInt64PtrIfNotNil(apiUser.CreatedAt),
		UpdatedAt:       pointer.ToInt64PtrIfNotNil(apiUser.UpdatedAt),
		Username:        pointer.ToStringPtrIfNotEmpty(apiUser.Username),
		Email:           pointer.ToStringPtrIfNotEmpty(apiUser.Email),
		FullName:        pointer.ToStringPtrIfNotEmpty(apiUser.RealName),
		Mobile:          pointer.ToStringPtrIfNotEmpty(apiUser.Mobile),
		Avatar:          pointer.ToStringPtrIfNotEmpty(apiUser.Avatar),
		UserDescription: pointer.ToStringPtrIfNotEmpty(apiUser.Desc),
		LastLoginAt:     pointer.ToInt64PtrIfNotNil(apiUser.LastLoginAt),
		LastLoginIp:     pointer.ToStringPtrIfNotEmpty(apiUser.LastLoginIp),
		State:           toBoolPtrIfNotFalse(apiUser.State),
		RoleIds:         apiUser.RoleIds,
		RoleValues:      apiUser.Roles,
		DepartmentId:    pointer.ToUint32PtrIfNotZero(apiUser.DepartmentId),
		PositionIds:     apiUser.PositionIds,
		HomePath:        pointer.ToStringPtrIfNotEmpty(apiUser.HomePath),
		RoleNames:       apiUser.RoleNames,
		DepartmentName:  pointer.ToStringPtrIfNotEmpty(apiUser.DepartmentName),
		PositionNames:   apiUser.PositionNames,
		PasswordHash:    apiUser.Password,
	}
}

// ConvertRpcUserListResponseToApiUserListResponse 将 RPC UserListResponse 转换为 API UserListResponse
func ConvertRpcUserListResponseToApiUserListResponse(rpcResp *userservice.UserListResponse) *types.UserListResponse {
	if rpcResp == nil {
		return nil
	}

	apiList := make([]types.UserInfo, 0, len(rpcResp.List))
	for _, user := range rpcResp.List {
		apiList = append(apiList, *ConvertRpcUserInfoToApiUserInfo(user))
	}

	return &types.UserListResponse{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "success",
		},
		Data: types.UserListInfo{
			BaseListInfo: types.BaseListInfo{
				Total: rpcResp.Page.Total,
			},
			List: apiList,
		},
	}
}

// toBoolPtrIfNotFalse 将 bool 转换为指针，如果为 false 则返回 nil
func toBoolPtrIfNotFalse(v bool) *bool {
	if !v {
		return nil
	}
	return &v
}
