package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spa-stc/newsletter/store"
)

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
