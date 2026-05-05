package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/che1nov/tea-shop/api-gateway/internal/requestctx"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddlewareSuccessAndRoleDefault(t *testing.T) {
	handler := AuthMiddleware("secret")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, roleExists := r.Context().Value(requestctx.RoleKey).(string)
		require.True(t, roleExists)
		w.WriteHeader(http.StatusOK)
	}))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": 1, "email": "u@test.local"})
	tokenString, err := token.SignedString([]byte("secret"))
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddlewareInvalidHeader(t *testing.T) {
	handler := AuthMiddleware("secret")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}
