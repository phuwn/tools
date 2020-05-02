package middleware

import (
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

// CorsConfig - CORS middleware
func CorsConfig() echo.MiddlewareFunc {
	return mw.CORSWithConfig(mw.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	})
}
