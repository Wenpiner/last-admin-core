package public_user

import (
	"context"
	"time"

	lastHttp "github.com/wenpiner/last-admin-common/utils/http"
	"github.com/wenpiner/last-admin-common/utils/pointer"
	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/client/oauthproviderservice"
	"github.com/wenpiner/last-admin-core/rpc/client/tokenservice"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type OauthCallbackLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Oauth 回调
func NewOauthCallbackLogic(r *http.Request, svcCtx *svc.ServiceContext) *OauthCallbackLogic {
	return &OauthCallbackLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *OauthCallbackLogic) OauthCallback() (resp *types.CallbackResponse, err error) {
	result, err := l.svcCtx.OauthRpc.OauthCallback(l.ctx, &oauthproviderservice.OauthCallbackRequest{
		State: l.r.FormValue("state"),
		Code:  l.r.FormValue("code"),
	})
	if err != nil {
		return nil, err
	}

	// 生成Token
	accessToken, err := generateToken(result, l.svcCtx.Config.Auth.AccessExpire, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(time.Second * time.Duration(l.svcCtx.Config.Auth.AccessExpire)).Unix()

	_, err = l.svcCtx.TokenRpc.CreateToken(l.ctx, &tokenservice.CreateTokenRequest{
		TokenValue: accessToken,
		TokenType:  "access_token",
		UserId:     result.Id,
		ExpiresAt:  expiresAt,
		DeviceInfo: pointer.ToStringPtr(l.r.UserAgent()),
		IpAddress:  pointer.ToStringPtr(lastHttp.GetIP(l.r)),
		UserAgent:  pointer.ToStringPtr(l.r.UserAgent()),
	})
	if err != nil {
		return nil, err
	}

	resp = &types.CallbackResponse{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "oauth.callbackSuccess",
		},
		Data: types.CallbackInfo{
			UserID:    *result.Id,
			Token:     accessToken, // 访问令牌 / Access token
			ExpiresAt: expiresAt,
		},
	}
	return
}
