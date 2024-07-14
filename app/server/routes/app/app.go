package app

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/spa-stc/newsletter/server/content"
	"github.com/spa-stc/newsletter/server/profile"
	"github.com/spa-stc/newsletter/store"
)

// App Service defining the routes for the main MPA.
type Service struct {
	profile *profile.Profile
	store   *store.Store
	content *content.Content
}

func NewService(p *profile.Profile, store *store.Store, content *content.Content) *Service {
	return &Service{
		profile: p,
		store:   store,
		content: content,
	}
}

// Register the app service onto the router.
func (s *Service) Register(ctx context.Context, echoServer *echo.Echo) {
	echoServer.GET("/", s.index, getPublicMiddleware()...)
	echoServer.GET("/about", s.about, getPublicMiddleware()...)
}
