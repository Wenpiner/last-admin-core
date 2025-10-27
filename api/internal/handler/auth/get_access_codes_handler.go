package auth

import (
	"net/http"

	"github.com/wenpiner/last-admin-core/api/internal/logic/auth"
	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取用户权限码(通过Menu获取按钮级别的权限)
func GetAccessCodesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := auth.NewGetAccessCodesLogic(r, svcCtx)
		resp, err := l.GetAccessCodes()
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
