package app

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/spa-stc/newsletter/server/profile"
	"github.com/spa-stc/newsletter/server/render"
	"github.com/spa-stc/newsletter/store"
)

// App Service defining the routes for the main MPA.
type Service struct {
	profile *profile.Profile
	info    render.SiteInfo
	store   *store.Store
}

func NewService(p *profile.Profile, store *store.Store) *Service {
	return &Service{
		profile: p,
		info:    getSiteInfo(p),
		store:   store,
	}
}

// Register the app service onto the router.
func (s *Service) Register(ctx context.Context, echoServer *echo.Echo) {
	echoServer.GET("/", s.index)
}

func getSiteInfo(p *profile.Profile) render.SiteInfo {
	return render.SiteInfo{
		Env: p.Env,
	}
}
