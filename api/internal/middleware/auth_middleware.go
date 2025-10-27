package middleware

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/wenpiner/last-admin-common/ctx/rolectx"
	last_i18n "github.com/wenpiner/last-admin-common/last-i18n"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type AuthMiddleware struct {
	trans *last_i18n.Translator
	cbn   *casbin.Enforcer
}

func NewAuthMiddleware(trans *last_i18n.Translator, cbn *casbin.Enforcer) *AuthMiddleware {
	return &AuthMiddleware{
		trans: trans,
		cbn:   cbn,
	}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		obj := r.URL.Path
		act := r.Method

		roles, ok := rolectx.GetRoleFromContext(r.Context())
		if !ok {
			httpx.Error(w, errorx.NewApiForbiddenError(m.trans.Trans(r.Context(), "common.forbidden")))
			return
		}

		// TODO Token 黑名单功能

		if check(m.cbn, roles, obj, act) {
			next(w, r)
		} else {
			httpx.Error(w, errorx.NewApiForbiddenError(m.trans.Trans(r.Context(), "common.api-forbidden")))
		}
	}
}

// Casbin check
func check(cbn *casbin.Enforcer, rolesIds []string, obj, act string) bool {
	var reqs [][]any
	for _, v := range rolesIds {
		reqs = append(reqs, []any{v, obj, act})
	}

	res, err := cbn.BatchEnforce(reqs)
	if err != nil {
		return false
	}

	for _, v := range res {
		if v {
			return true
		}
	}

	return false
}
