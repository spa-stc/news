package app

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimd "github.com/go-chi/chi/v5/middleware"
	"stpaulacademy.tech/newsletter/db"
	"stpaulacademy.tech/newsletter/resource"
	"stpaulacademy.tech/newsletter/util/service"
	"stpaulacademy.tech/newsletter/util/timeutil"
	"stpaulacademy.tech/newsletter/web"
	"stpaulacademy.tech/newsletter/web/assets"
	"stpaulacademy.tech/newsletter/web/templates"
)

func NewServer(
	logger *slog.Logger,
	a *assets.Assets,
	rootAssets *assets.Assets,
	t *templates.TemplateRenderer,
	e db.Executor,
	timeGetter service.TimeGenerator,
) http.Handler {
	w := web.NewHandlerWrapper(logger)
	r := chi.NewMux()

	r.Use(chimd.RealIP)
	r.Use(chimd.Compress(5))

	r.Method(http.MethodGet, "/healthz", w.Wrap(handleHealthz))
	r.Method(http.MethodGet, "/assets/{hash}", w.Wrap(web.ServeStatics(a)))
	r.Method(http.MethodGet, "/", w.Wrap(handleIndex(t, e, timeGetter)))

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

func handleIndex(t *templates.TemplateRenderer, e db.Executor, timeGen service.TimeGenerator) web.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		week := timeutil.GetWeek(time.Sunday, timeGen.NowUTC())
		days, err := resource.GetManyDays(r.Context(), e, week[1:6])
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				return web.RespondError("Not Found.", http.StatusNotFound, err)
			}

			return err
		}

		announcements, err := resource.GetManyAnnouncementsByCurrentDay(r.Context(), e, timeGen.NowUTC())
		if err != nil {
			return err
		}

		data := struct {
			Announcements []resource.Announcement
			Days          []resource.Day
			DayUpdatedTS  time.Time
		}{
			Announcements: announcements,
			Days:          days,
			DayUpdatedTS:  days[0].UpdatedTS,
		}

		return web.RenderTemplate(w, t, "index.html", web.TemplateCachePolicyPublic, templates.RenderData{
			Data: data,
		})
	}
}
