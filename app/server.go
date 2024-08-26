package app

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"stpaulacademy.tech/newsletter/web"
)

func NewServer(logger *slog.Logger) http.Handler {
	r := chi.NewMux()

	r.Method(http.MethodGet, "/healthz", web.Handler(handleHealthz))

	return r
}

func handleHealthz(w http.ResponseWriter, _ *http.Request) error {
	_, err := w.Write([]byte("Service Ready."))
	return err
}
