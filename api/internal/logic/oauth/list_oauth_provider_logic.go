package oauth

import (
	"context"

	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/client/oauthproviderservice"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListOauthProviderLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取Oauth列表
func NewListOauthProviderLogic(r *http.Request, svcCtx *svc.ServiceContext) *ListOauthProviderLogic {
	return &ListOauthProviderLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *ListOauthProviderLogic) ListOauthProvider(req *types.OauthProviderListRequest) (resp *types.OauthProviderListResponse, err error) {
	// 构建 RPC 请求
	rpcReq := &oauthproviderservice.OauthProviderListRequest{
		Page: &oauthproviderservice.BasePageRequest{
			PageNumber: req.Page.CurrentPage,
			PageSize:   req.Page.PageSize,
		},
		ProviderName: pointer.ToStringPtrIfNotEmpty(req.ProviderName),
		ProviderCode: pointer.ToStringPtrIfNotEmpty(req.ProviderCode),
		State:        req.State,
	}

	// 调用 RPC 服务获取 OAuth 提供商列表
	rpcResp, err := l.svcCtx.OauthRpc.ListOauthProvider(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 将 RPC 响应转换为 API 响应
	apiList := make([]types.OauthProviderInfo, 0, len(rpcResp.List))
	for _, provider := range rpcResp.List {
		apiList = append(apiList, *convertRpcOauthProviderInfoToApiOauthProviderInfo(provider))
	}

	return &types.OauthProviderListResponse{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "success",
		},
		Data: types.OauthProviderInfoList{
			BaseListInfo: types.BaseListInfo{
				Total: rpcResp.Page.Total,
			},
			List: apiList,
		},
	},nil

	return resp, nil
}

// convertRpcOauthProviderInfoToApiOauthProviderInfo 将 RPC OauthProviderInfo 转换为 API OauthProviderInfo
func convertRpcOauthProviderInfoToApiOauthProviderInfo(rpcInfo *oauthproviderservice.OauthProviderInfo) *types.OauthProviderInfo {
	if rpcInfo == nil {
		return nil
	}

	var authStyle *uint8
	if rpcInfo.AuthStyle != nil {
		authStyleVal := uint8(*rpcInfo.AuthStyle)
		authStyle = &authStyleVal
	}

	return &types.OauthProviderInfo{
		ID:               rpcInfo.Id,
		CreatedAt:        rpcInfo.CreatedAt,
		UpdatedAt:        rpcInfo.UpdatedAt,
		State:            rpcInfo.State,
		ProviderName:     rpcInfo.ProviderName,
		ProviderCode:     rpcInfo.ProviderCode,
		ClientId:         rpcInfo.ClientId,
		ClientSecret:     rpcInfo.ClientSecret,
		RedirectUri:      rpcInfo.RedirectUri,
		Scopes:           rpcInfo.Scopes,
		AuthorizationUrl: rpcInfo.AuthorizationUrl,
		TokenUrl:         rpcInfo.TokenUrl,
		UserinfoUrl:      rpcInfo.UserinfoUrl,
		LogoutUrl:        rpcInfo.LogoutUrl,
		AuthStyle:        authStyle,
	}
}
