package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			c.Abort()
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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		c.Set("user_id", int64(claims["user_id"].(float64)))
		c.Set("email", claims["email"].(string))
		
		// Добавляем роль в контекст (если есть в токене)
		if role, ok := claims["role"].(string); ok {
			c.Set("role", role)
		} else {
			// Если роль не указана, устанавливаем "user" по умолчанию
			c.Set("role", "user")
		}

		c.Next()
	}
}
