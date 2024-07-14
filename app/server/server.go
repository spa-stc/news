package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echolog "github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"github.com/spa-stc/newsletter/server/profile"
	"github.com/spa-stc/newsletter/server/render"
	"github.com/spa-stc/newsletter/server/routes/app"
	"github.com/spa-stc/newsletter/server/static"
	"github.com/spa-stc/newsletter/store"
)

type Server struct {
	profile *profile.Profile

	echoServer *echo.Echo
	templ      *render.Templates
	store      *store.Store
}

func New(ctx context.Context, p *profile.Profile, templ *render.Templates, store *store.Store) *Server {
	echoServer := echo.New()
	echoServer.Logger.SetLevel(echolog.OFF)
	echoServer.HideBanner = true
	echoServer.HidePort = true
	echoServer.Debug = true
	echoServer.Use(middleware.Recover())
	echoServer.Use(middleware.GzipWithConfig(
		middleware.GzipConfig{
			Level: 5,
		},
	))
	echoServer.Renderer = templ

	static.Register(echoServer)

	echoServer.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "Service	Ready.")
	})

	app.NewService(p, store).Register(ctx, echoServer)

	s := &Server{
		profile:    p,
		echoServer: echoServer,
		store:      store,
		templ:      templ,
	}

	return s
}

func (s *Server) Start(ctx context.Context) {
	addr := fmt.Sprintf("localhost:%s", s.profile.Port)

	go func() {
		if err := s.echoServer.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start echo server", "error", err)
			os.Exit(1)
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := s.echoServer.Shutdown(ctx); err != nil {
		slog.Error("failed to stop echo server", "error", err)
	}
}
