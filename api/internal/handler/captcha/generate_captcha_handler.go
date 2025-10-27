package captcha

import (
	"net/http"

	"github.com/wenpiner/last-admin-core/api/internal/logic/captcha"
	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 生成验证码
func GenerateCaptchaHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := captcha.NewGenerateCaptchaLogic(r, svcCtx)
		resp, err := l.GenerateCaptcha()
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
