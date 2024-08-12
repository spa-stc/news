package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"stpaulacademy.tech/newsletter/config"
	"stpaulacademy.tech/newsletter/db"
	"stpaulacademy.tech/newsletter/server/api"
	"stpaulacademy.tech/newsletter/util/service"
)

type Server struct {
	s *http.Server
}

func New(c config.Config, e db.Executor) *Server {
	api := api.NewServer(&service.TimeGen{}, e)

	s := &http.Server{
		Handler:           api,
		Addr:              fmt.Sprintf("%s:%d", c.Host, c.Port),
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second, //nolint:mnd // fine
		ReadHeaderTimeout: 2 * time.Second,  //nolint:mnd // fine
	}

	return &Server{
		s: s,
	}
}

func (s *Server) Run() {
	go func() {
		if err := s.s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("error serving over http", "err", err)
			os.Exit(1)
		}
	}()
}

func (s *Server) Stop(ctx context.Context) {
	if err := s.s.Shutdown(ctx); err != nil {
		slog.ErrorContext(ctx, "error stopping http server", "err", err)
	}
}
