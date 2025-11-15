package public_config

import (
	"net/http"

	"github.com/wenpiner/last-admin-core/api/internal/logic/public_config"
	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取VBen Preference配置
func GetVbenPreferenceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public_config.NewGetVbenPreferenceLogic(r, svcCtx)
		resp, err := l.GetVbenPreference()
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
