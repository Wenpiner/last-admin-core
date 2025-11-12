package oauth

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/client/oauthproviderservice"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteOauthProviderLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除Oauth
func NewDeleteOauthProviderLogic(r *http.Request, svcCtx *svc.ServiceContext) *DeleteOauthProviderLogic {
	return &DeleteOauthProviderLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *DeleteOauthProviderLogic) DeleteOauthProvider(req *types.ID32Request) (resp *types.BaseResponse, err error) {
	// 调用 RPC 服务删除 OAuth 提供商
	_, err = l.svcCtx.OauthRpc.DeleteOauthProvider(l.ctx, &oauthproviderservice.ID32Request{Id: req.ID})
	if err != nil {
		return nil, err
	}

	resp = &types.BaseResponse{
		Code:    0,
		Message: "success",
	}

	return
}
