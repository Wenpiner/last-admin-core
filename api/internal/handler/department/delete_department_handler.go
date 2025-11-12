package department

import (
	"net/http"

	"github.com/wenpiner/last-admin-core/api/internal/logic/department"
	"github.com/wenpiner/last-admin-core/api/internal/svc"
	"github.com/wenpiner/last-admin-core/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 删除部门
func DeleteDepartmentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ID32Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := department.NewDeleteDepartmentLogic(r, svcCtx)
		resp, err := l.DeleteDepartment(&req)
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
