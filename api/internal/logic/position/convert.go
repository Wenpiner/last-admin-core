package position

import (
	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/client/positionservice"
)

// ConvertRpcPositionInfoToApiPositionInfo 将 RPC PositionInfo 转换为 API PositionInfo
func ConvertRpcPositionInfoToApiPositionInfo(rpcPos *positionservice.PositionInfo) *types.PositionInfo {
	if rpcPos == nil {
		return nil
	}

	return &types.PositionInfo{
		ID:           rpcPos.Id,
		CreatedAt:    rpcPos.CreatedAt,
		UpdatedAt:    rpcPos.UpdatedAt,
		PositionName: pointer.GetString(rpcPos.PositionName),
		PositionCode: pointer.GetString(rpcPos.PositionCode),
		SortOrder:    pointer.GetInt32(rpcPos.SortOrder),
		State:        rpcPos.State,
		Description:  rpcPos.Description,
	}
}

// ConvertApiPositionInfoToRpcPositionInfo 将 API PositionInfo 转换为 RPC PositionInfo
func ConvertApiPositionInfoToRpcPositionInfo(apiPos *types.PositionInfo) *positionservice.PositionInfo {
	if apiPos == nil {
		return nil
	}

	return &positionservice.PositionInfo{
		Id:           apiPos.ID,
		CreatedAt:    apiPos.CreatedAt,
		UpdatedAt:    apiPos.UpdatedAt,
		PositionName: pointer.ToStringPtrIfNotEmpty(apiPos.PositionName),
		PositionCode: pointer.ToStringPtrIfNotEmpty(apiPos.PositionCode),
		SortOrder:    pointer.ToInt32PtrIfNotZero(apiPos.SortOrder),
		State:        apiPos.State,
		Description:  apiPos.Description,
	}
}

// ConvertRpcPositionListResponseToApiPositionListResponse 将 RPC PositionListResponse 转换为 API PositionListResponse
func ConvertRpcPositionListResponseToApiPositionListResponse(rpcResp *positionservice.PositionListResponse) *types.PositionListResponse {
	if rpcResp == nil {
		return nil
	}

	apiList := make([]types.PositionInfo, 0, len(rpcResp.List))
	for _, pos := range rpcResp.List {
		apiList = append(apiList, *ConvertRpcPositionInfoToApiPositionInfo(pos))
	}

	return &types.PositionListResponse{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "success",
		},
		Data: []types.PositionListInfo{
			{
				BaseListInfo: types.BaseListInfo{
					Total: rpcResp.Page.Total,
				},
				List: apiList,
			},
		},
	}
}

