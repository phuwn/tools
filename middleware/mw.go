package middleware

import (
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

// RemoveTrailingSlash - fork mw with the same name of labstack
var RemoveTrailingSlash echo.MiddlewareFunc = mw.RemoveTrailingSlash()
