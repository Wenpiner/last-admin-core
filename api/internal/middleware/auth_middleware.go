package middleware

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/redis/go-redis/v9"
	"github.com/wenpiner/last-admin-common/ctx/rolectx"
	last_i18n "github.com/wenpiner/last-admin-common/last-i18n"
	last_redis "github.com/wenpiner/last-admin-common/last-redis"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type AuthMiddleware struct {
	trans *last_i18n.Translator
	cbn   *casbin.Enforcer
	rds   *redis.Client
}

func NewAuthMiddleware(trans *last_i18n.Translator, cbn *casbin.Enforcer, rds *redis.Client) *AuthMiddleware {
	return &AuthMiddleware{
		trans: trans,
		cbn:   cbn,
		rds:   rds,
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

		// 查询redis中的黑名单对应的token是否存在
		key := string(last_redis.BlacklistToken)
		token := r.Header.Get("Authorization")
		if token != "" {
			token = token[7:]
			_, err := m.rds.SIsMember(r.Context(), key, token).Result()
			if err != nil {
				httpx.Error(w, errorx.NewApiForbiddenError(m.trans.Trans(r.Context(), "common.forbidden")))
				return
			}
		}

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
		reqs = append(reqs, []any{v, "api", obj, act})
	}

	res, err := cbn.BatchEnforce(reqs)
	if err != nil {
		logx.Errorw("验证 Casbin 异常", logx.Field("error", err))
		return false
	}

	for _, v := range res {
		if v {
			return true
		}
	}

	return false
}
