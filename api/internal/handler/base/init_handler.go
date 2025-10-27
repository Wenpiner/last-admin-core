package base

import (
	"net/http"

	"github.com/wenpiner/last-admin-core/api/internal/logic/base"
	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func InitHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := base.NewInitLogic(r, svcCtx)
		resp, err := l.Init()
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
