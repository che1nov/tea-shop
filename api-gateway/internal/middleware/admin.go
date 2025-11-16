package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const RoleAdmin = "admin"

// AdminMiddleware проверяет, что пользователь имеет роль администратора
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем роль из контекста (устанавливается в AuthMiddleware)
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "role not found in token"})
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid role type"})
			c.Abort()
			return
		}

		if roleStr != RoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied: admin role required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

