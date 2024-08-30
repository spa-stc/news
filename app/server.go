package app

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimd "github.com/go-chi/chi/v5/middleware"
	"stpaulacademy.tech/newsletter/db"
	"stpaulacademy.tech/newsletter/web"
	"stpaulacademy.tech/newsletter/web/assets"
	"stpaulacademy.tech/newsletter/web/templates"
)

func NewServer(
	logger *slog.Logger,
	a *assets.Assets,
	rootAssets *assets.Assets,
	t *templates.TemplateRenderer,
	_ db.Executor,
) http.Handler {
	w := web.NewHandlerWrapper(logger)
	r := chi.NewMux()

	r.Use(chimd.RealIP)
	r.Use(chimd.Compress(5))

	r.Method(http.MethodGet, "/healthz", w.Wrap(handleHealthz))
	r.Method(http.MethodGet, "/assets/{hash}", w.Wrap(web.ServeStatics(a)))
	r.Method(http.MethodGet, "/", w.Wrap(handleIndex(t)))

	r.NotFound(func(writer http.ResponseWriter, r *http.Request) {
		w.Wrap(web.ServeRootStatics(rootAssets)).ServeHTTP(writer, r)
	})

	return r
}

func handleHealthz(w http.ResponseWriter, _ *http.Request) error {
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte("Service Ready."))
	return err
}

func handleIndex(t *templates.TemplateRenderer) web.Handler {
	return func(w http.ResponseWriter, _ *http.Request) error {
		return web.RenderTemplate(w, t, "index.html", web.TemplateCachePolicyPublic, templates.RenderData{
			Data: nil,
		})
	}
}
