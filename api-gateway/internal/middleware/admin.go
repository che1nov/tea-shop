package middleware

import (
	"net/http"

	"github.com/che1nov/tea-shop/api-gateway/internal/requestctx"
)

const RoleAdmin = "admin"

func AdminMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(requestctx.RoleKey).(string)
			if !ok {
				writeError(w, http.StatusForbidden, "role not found in token")
				return
			}

			if role != RoleAdmin {
				writeError(w, http.StatusForbidden, "access denied: admin role required")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
