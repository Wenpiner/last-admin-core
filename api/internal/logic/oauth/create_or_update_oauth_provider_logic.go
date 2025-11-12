package oauth

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/client/oauthproviderservice"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrUpdateOauthProviderLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建或更新Oauth
func NewCreateOrUpdateOauthProviderLogic(r *http.Request, svcCtx *svc.ServiceContext) *CreateOrUpdateOauthProviderLogic {
	return &CreateOrUpdateOauthProviderLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *CreateOrUpdateOauthProviderLogic) CreateOrUpdateOauthProvider(req *types.OauthProviderInfo) (resp *types.OauthModifyResponse, err error) {
	// 构建 RPC 请求
	var authStyle *uint32
	if req.AuthStyle != nil {
		authStyleVal := uint32(*req.AuthStyle)
		authStyle = &authStyleVal
	}

	rpcReq := &oauthproviderservice.OauthProviderInfo{
		Id:               req.ID,
		CreatedAt:        req.CreatedAt,
		UpdatedAt:        req.UpdatedAt,
		State:            req.State,
		ProviderName:     req.ProviderName,
		ProviderCode:     req.ProviderCode,
		ClientId:         req.ClientId,
		ClientSecret:     req.ClientSecret,
		RedirectUri:      req.RedirectUri,
		Scopes:           req.Scopes,
		AuthorizationUrl: req.AuthorizationUrl,
		TokenUrl:         req.TokenUrl,
		UserinfoUrl:      req.UserinfoUrl,
		LogoutUrl:        req.LogoutUrl,
		AuthStyle:        authStyle,
	}

	// 调用 RPC 服务创建或更新 OAuth 提供商
	rpcResp, err := l.svcCtx.OauthRpc.CreateOrUpdateOauthProvider(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 将 RPC 响应转换为 API 响应
	info := convertRpcOauthProviderInfoToApiOauthProviderInfo(rpcResp)
	resp = &types.OauthModifyResponse{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "success",
		},
		Data: *info,
	}
	return resp, nil
}
