package public_user

import (
	"context"

	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/wenpiner/last-admin-core/rpc/client/oauthproviderservice"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type OauthLoginLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Oauth 登录
func NewOauthLoginLogic(r *http.Request, svcCtx *svc.ServiceContext) *OauthLoginLogic {
	return &OauthLoginLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
		ctx:    r.Context(),
	}
}

func (l *OauthLoginLogic) OauthLogin(req *types.OauthLoginRequest) (resp *types.OauthRedirectResponse, err error) {
	result, err := l.svcCtx.OauthRpc.OauthLogin(l.ctx, &oauthproviderservice.OauthLoginRequest{
		State:    req.State,
		Provider: req.Provider,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.OauthRedirectResponse{
		BaseDataInfo: types.BaseDataInfo{
			Code:    0,
			Message: "oauth.loginSuccess",
		},
		Data: result.Url,
	}

	return
}
