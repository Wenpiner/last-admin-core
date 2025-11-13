package configuration

import (
	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/client/configurationservice"
)

// ConvertRpcConfigurationInfoToApiConfigurationInfo 将 RPC ConfigurationInfo 转换为 API ConfigurationInfo
func ConvertRpcConfigurationInfoToApiConfigurationInfo(rpcConfig *configurationservice.ConfigurationInfo) *types.ConfigurationInfo {
	if rpcConfig == nil {
		return nil
	}

	return &types.ConfigurationInfo{
		Key:         rpcConfig.Key,
		Value:       rpcConfig.Value,
		Name:        rpcConfig.Name,
		Group:       rpcConfig.Group,
		Description: rpcConfig.Description,
	}
}

// ConvertApiConfigurationInfoToRpcConfigurationInfo 将 API ConfigurationInfo 转换为 RPC ConfigurationInfo
func ConvertApiConfigurationInfoToRpcConfigurationInfo(apiConfig *types.ConfigurationInfo) *configurationservice.ConfigurationInfo {
	if apiConfig == nil {
		return nil
	}

	return &configurationservice.ConfigurationInfo{
		Key:         apiConfig.Key,
		Value:       apiConfig.Value,
		Name:        apiConfig.Name,
		Group:       apiConfig.Group,
		Description: apiConfig.Description,
	}
}

// ConvertRpcConfigurationListResponseToApiConfigurationListResponse 将 RPC ConfigurationListResponse 转换为 API ConfigurationListResponse
func ConvertRpcConfigurationListResponseToApiConfigurationListResponse(rpcResp *configurationservice.ConfigurationListResponse) *types.ConfigurationListResponse {
	if rpcResp == nil {
		return nil
	}

	apiList := make([]types.ConfigurationInfo, 0, len(rpcResp.List))
	for _, config := range rpcResp.List {
		apiList = append(apiList, *ConvertRpcConfigurationInfoToApiConfigurationInfo(config))
	}

	return &types.ConfigurationListResponse{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "success",
		},
		Data: types.ConfigurationListInfo{
			BaseListInfo: types.BaseListInfo{
				Total: rpcResp.Page.Total,
			},
			List: apiList,
		},
	}
}

// ConvertApiConfigurationListRequestToRpcConfigurationListRequest 将 API ConfigurationListRequest 转换为 RPC ConfigurationListRequest
func ConvertApiConfigurationListRequestToRpcConfigurationListRequest(apiReq *types.ConfigurationListRequest) *configurationservice.ConfigurationListRequest {
	if apiReq == nil {
		return nil
	}

	return &configurationservice.ConfigurationListRequest{
		Page: &configurationservice.BasePageRequest{
			PageNumber: apiReq.Page.CurrentPage,
			PageSize:   apiReq.Page.PageSize,
		},
		Key:   pointer.ToStringPtrIfNotEmpty(apiReq.Key),
		Name:  pointer.ToStringPtrIfNotEmpty(apiReq.Name),
		Group: pointer.ToStringPtrIfNotEmpty(apiReq.Group),
	}
}

