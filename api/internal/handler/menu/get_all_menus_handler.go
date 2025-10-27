package menu

import (
	"net/http"

	"github.com/wenpiner/last-admin-core/api/internal/logic/menu"
	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取所有菜单
func GetAllMenusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := menu.NewGetAllMenusLogic(r, svcCtx)
		resp, err := l.GetAllMenus()
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
