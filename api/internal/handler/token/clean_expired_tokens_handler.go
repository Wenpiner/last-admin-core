package token

import (
	"net/http"

	"github.com/wenpiner/last-admin-core/api/internal/logic/token"
	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 清理过期令牌
func CleanExpiredTokensHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CleanExpiredTokensRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := token.NewCleanExpiredTokensLogic(r, svcCtx)
		resp, err := l.CleanExpiredTokens(&req)
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
