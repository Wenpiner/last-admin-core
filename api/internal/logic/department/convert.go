package department

import (
	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/client/departmentservice"
)

// ConvertRpcDepartmentInfoToApiDepartmentInfo 将 RPC DepartmentInfo 转换为 API DepartmentInfo
func ConvertRpcDepartmentInfoToApiDepartmentInfo(rpcDept *departmentservice.DepartmentInfo) *types.DepartmentInfo {
	if rpcDept == nil {
		return nil
	}

	return &types.DepartmentInfo{
		ID:             rpcDept.Id,
		CreatedAt:      rpcDept.CreatedAt,
		UpdatedAt:      rpcDept.UpdatedAt,
		DeptName:       pointer.GetString(rpcDept.DeptName),
		DeptCode:       pointer.GetString(rpcDept.DeptCode),
		ParentId:       rpcDept.ParentId,
		SortOrder:      pointer.GetInt32(rpcDept.SortOrder),
		LeaderUserId:   pointer.GetString(rpcDept.LeaderUserId),
		State:          pointer.GetBool(rpcDept.State),
		Description:    rpcDept.Description,
		LeaderUsername: rpcDept.LaderUsername,
		LeaderPhone:    rpcDept.LaderPhone,
		LeaderEmail:    rpcDept.LaderEmail,
	}
}

// ConvertApiDepartmentInfoToRpcDepartmentInfo 将 API DepartmentInfo 转换为 RPC DepartmentInfo
func ConvertApiDepartmentInfoToRpcDepartmentInfo(apiDept *types.DepartmentInfo) *departmentservice.DepartmentInfo {
	if apiDept == nil {
		return nil
	}

	return &departmentservice.DepartmentInfo{
		Id:           apiDept.ID,
		CreatedAt:    apiDept.CreatedAt,
		UpdatedAt:    apiDept.UpdatedAt,
		DeptName:     pointer.ToStringPtrIfNotEmpty(apiDept.DeptName),
		DeptCode:     pointer.ToStringPtrIfNotEmpty(apiDept.DeptCode),
		ParentId:     apiDept.ParentId,
		SortOrder:    pointer.ToInt32PtrIfNotZero(apiDept.SortOrder),
		LeaderUserId: pointer.ToStringPtrIfNotEmpty(apiDept.LeaderUserId),
		State:        pointer.ToBoolPtrIfNotFalse(apiDept.State),
		Description:  apiDept.Description,
	}
}
// ConvertRpcDepartmentListResponseToApiDepartmentListResponse 将 RPC DepartmentListResponse 转换为 API DepartmentListResponse
func ConvertRpcDepartmentListResponseToApiDepartmentListResponse(rpcResp *departmentservice.DepartmentListResponse) *types.DepartmentListResponse {
	if rpcResp == nil {
		return nil
	}

	apiList := make([]types.DepartmentInfo, 0, len(rpcResp.List))
	for _, dept := range rpcResp.List {
		apiList = append(apiList, *ConvertRpcDepartmentInfoToApiDepartmentInfo(dept))
	}

	return &types.DepartmentListResponse{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "success",
		},
		Data: types.DepartmentListInfo{
			BaseListInfo: types.BaseListInfo{
				Total: rpcResp.Page.Total,
			},
			List: apiList,
		},
	}
}
