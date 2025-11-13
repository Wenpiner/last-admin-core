package role

import (
	"net/http"

	"github.com/wenpiner/last-admin-core/api/internal/logic/role"
	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取当前系统中的所有分组列表
func GetConfigurationGroupListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := role.NewGetConfigurationGroupListLogic(r, svcCtx)
		resp, err := l.GetConfigurationGroupList()
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
