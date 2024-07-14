package static

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spa-stc/newsletter/public"
)

func Register(router *echo.Echo) {
	router.Use(
		middleware.StaticWithConfig(middleware.StaticConfig{
			HTML5:      false,
			Root:       "static",
			Filesystem: http.FS(public.Static),
		}),
	)
}
