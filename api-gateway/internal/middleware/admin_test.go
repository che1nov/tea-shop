package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"context"

	"github.com/che1nov/tea-shop/api-gateway/internal/requestctx"
	"github.com/stretchr/testify/require"
)

func TestAdminMiddleware(t *testing.T) {
	t.Run("admin", func(t *testing.T) {
		handler := AdminMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req = req.WithContext(context.WithValue(req.Context(), requestctx.RoleKey, "admin"))
		handler.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("not-admin", func(t *testing.T) {
		handler := AdminMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req = req.WithContext(context.WithValue(req.Context(), requestctx.RoleKey, "user"))
		handler.ServeHTTP(w, req)
		require.Equal(t, http.StatusForbidden, w.Code)
	})
}
