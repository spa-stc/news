package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spa-stc/newsletter/server/render"
	"github.com/spa-stc/newsletter/store"
)

func (a *Service) index(c echo.Context) error {
	days, err := a.store.GetWeek(c.Request().Context())
	if err != nil {
		return err
	}

	data := struct {
		render.BaseContext

		Days []store.Day
	}{
		render.BaseContext{
			Title: "Home",
			Info:  a.info,
		},

		days,
	}

	return c.Render(http.StatusOK, "pages/index.tmpl.html", &data)
}
