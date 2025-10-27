package menu

import (
	"net/http"

	"github.com/wenpiner/last-admin-core/api/internal/logic/menu"
	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取用户角色当前所有菜单
func GetAllMenusByRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := menu.NewGetAllMenusByRoleLogic(r, svcCtx)
		resp, err := l.GetAllMenusByRole()
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
