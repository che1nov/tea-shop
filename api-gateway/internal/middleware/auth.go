package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/che1nov/tea-shop/api-gateway/internal/requestctx"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				writeError(w, http.StatusUnauthorized, "invalid authorization header")
				return
			}

			tokenString := parts[1]
			token, err := jwt.ParseWithClaims(
				tokenString,
				jwt.MapClaims{},
				func(token *jwt.Token) (interface{}, error) {
					return []byte(jwtSecret), nil
				},
			)

			if err != nil || !token.Valid {
				writeError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				writeError(w, http.StatusUnauthorized, "invalid token claims")
				return
			}

			ctx := context.WithValue(r.Context(), requestctx.UserIDKey, int64(claims["user_id"].(float64)))
			ctx = context.WithValue(ctx, requestctx.EmailKey, claims["email"].(string))

			if role, ok := claims["role"].(string); ok {
				ctx = context.WithValue(ctx, requestctx.RoleKey, role)
			} else {
				ctx = context.WithValue(ctx, requestctx.RoleKey, "user")
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
