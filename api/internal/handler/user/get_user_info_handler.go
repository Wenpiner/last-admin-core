package user

import (
	"net/http"

	"github.com/wenpiner/last-admin-core/api/internal/logic/user"
	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取用户信息
func GetUserInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetUserInfoLogic(r, svcCtx)
		resp, err := l.GetUserInfo()
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
