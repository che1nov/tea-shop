package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestAdminMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("admin", func(t *testing.T) {
		r := gin.New()
		r.Use(func(c *gin.Context) { c.Set("role", "admin") })
		r.Use(AdminMiddleware())
		r.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("not-admin", func(t *testing.T) {
		r := gin.New()
		r.Use(func(c *gin.Context) { c.Set("role", "user") })
		r.Use(AdminMiddleware())
		r.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
		require.Equal(t, http.StatusForbidden, w.Code)
	})
}
