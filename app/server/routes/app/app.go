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
	info    render.SiteInfo
}

func NewService(p *profile.Profile) *Service {
	return &Service{
		profile: p,
		info:    getSiteInfo(p),
	}
}

// Register the app service onto the router.
func (s *Service) Register(ctx context.Context, echoServer *echo.Echo) {
	echoServer.GET("/", func(c echo.Context) error {
		data := struct {
			render.BaseContext

			Message string
		}{
			BaseContext: render.BaseContext{
				Title: "Hello",
				Info:  s.info,
			},
			Message: "World",
		}

		return c.Render(200, "pages/index.tmpl.html", &data)
	})
}

func getSiteInfo(p *profile.Profile) render.SiteInfo {
	return render.SiteInfo{
		Env: p.Env,
	}
}
