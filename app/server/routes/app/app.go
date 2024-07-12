package app

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/spa-stc/newsletter/server/profile"
	"github.com/spa-stc/newsletter/server/render"
)

// App Service defining the routes for the main MPA.
type Service struct {
	profile *profile.Profile
}

func NewService(p *profile.Profile) *Service {
	return &Service{
		profile: p,
	}
}

// Register the app service onto the router.
func (*Service) Register(ctx context.Context, echoServer *echo.Echo) {
	echoServer.GET("/", func(c echo.Context) error {
		data := render.BaseContext{
			Title: "Hello",

			Info: render.SiteInfo{},
		}

		return c.Render(200, "layouts/main.tmpl.html", &data)
	})
}
