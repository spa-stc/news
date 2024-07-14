package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spa-stc/newsletter/store"
)

func getPublicMiddleware() []echo.MiddlewareFunc {
	cachecontrol := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", "public, max-age=3600, stale-if-error=60")

			return next(c)
		}
	}

	return []echo.MiddlewareFunc{
		cachecontrol,
	}
}

func (a *Service) index(c echo.Context) error {
	days, err := a.store.GetWeek(c.Request().Context())
	if err != nil {
		return err
	}

	data := struct {
		Title string

		Days []store.Day
	}{
		"Home",
		days,
	}

	return c.Render(http.StatusOK, "pages/index.tmpl.html", &data)
}
