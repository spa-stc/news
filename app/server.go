package app

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewServer(logger *slog.Logger) http.Handler {
	r := chi.NewMux()

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte("Service Ready."))
		if err != nil {
			logger.Error("error sending healthz response", "error", err)
		}
	})

	return r
}
