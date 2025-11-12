package api

import (
	"net/http"

	"github.com/wenpiner/last-admin-core/api/internal/logic/api"
	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取所有API
func GetAllApiHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := api.NewGetAllApiLogic(r, svcCtx)
		resp, err := l.GetAllApi()
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
