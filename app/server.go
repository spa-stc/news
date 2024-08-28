package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimd "github.com/go-chi/chi/v5/middleware"
	"stpaulacademy.tech/newsletter/web"
)

func NewServer() http.Handler {
	r := chi.NewMux()

	r.Use(chimd.RealIP)
	r.Use(chimd.Compress(5))

	r.Method(http.MethodGet, "/healthz", web.Handler(handleHealthz))

	return r
}

func handleHealthz(w http.ResponseWriter, _ *http.Request) error {
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte("Service Ready."))
	return err
}
