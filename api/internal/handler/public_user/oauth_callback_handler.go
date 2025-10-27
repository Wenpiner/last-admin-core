package public_user

import (
	"net/http"

	"github.com/wenpiner/last-admin-core/api/internal/logic/public_user"
	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Oauth 回调
func OauthCallbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public_user.NewOauthCallbackLogic(r, svcCtx)
		resp, err := l.OauthCallback()
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
