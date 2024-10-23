package app

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	chimd "github.com/go-chi/chi/v5/middleware"
	"github.com/gomarkdown/markdown"
	"stpaulacademy.tech/newsletter/config"
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
	c config.Config,
) http.Handler {
	w := web.NewHandlerWrapper(logger)
	r := chi.NewMux()

	r.Use(chimd.RealIP)
	r.Use(chimd.Compress(5))

	r.Method(http.MethodGet, "/healthz", w.Wrap(handleHealthz))
	r.Method(http.MethodGet, "/assets/{hash}", w.Wrap(web.ServeStatics(a)))
	r.Method(http.MethodGet, "/", w.Wrap(handleIndex(t, e, timeGetter)))

	r.Route("/admin", func(r chi.Router) {
		r.Use(chimd.BasicAuth("spa-newsletter", map[string]string{
			c.AdminUsername: c.AdminPassword,
		}))
		r.Method(http.MethodGet, "/", w.Wrap(handleAdmin(t, e, timeGetter)))
		r.Method(http.MethodPost, "/announcements", w.Wrap(handleSubmit(logger, e)))
		r.Method(http.MethodDelete, "/announcements/{id}", w.Wrap(handleDelete(e)))
	})

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
		week := timeutil.GetWeek(time.Sunday, timeGen.NowUTC().Local())
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
			ShowToday     bool
			Today         time.Time
		}{
			Announcements: announcements,
			Days:          days,
			DayUpdatedTS:  days[0].UpdatedTS,
			ShowToday:     timeutil.IsWeekday(timeGen.NowUTC().Local()),
			Today:         timeGen.NowUTC().Local(),
		}

		return web.RenderTemplate(w, t, "index.html", web.TemplateCachePolicyPublic, templates.RenderData{
			Data: data,
		})
	}
}

func handleAdmin(t *templates.TemplateRenderer, e db.Executor, timeGen service.TimeGenerator) web.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		announcements, err := resource.GetUpcomingAnnouncements(r.Context(), e, timeGen.NowUTC())
		if err != nil {
			return err
		}

		data := struct {
			Announcements []resource.Announcement
		}{
			Announcements: announcements,
		}

		return web.RenderTemplate(w, t, "admin.html", web.TemplateCachePolicyPrivate, templates.RenderData{
			Data: data,
		})
	}
}

func handleDelete(e db.Executor) web.Handler {
	return func(_ http.ResponseWriter, r *http.Request) error {
		rawID := r.PathValue("id")

		id, err := strconv.Atoi(rawID)
		if err != nil {
			return web.RespondError("Invalid ID.", http.StatusBadRequest, err)
		}

		if err := resource.DeleteAnnouncement(r.Context(), e, int64(id)); err != nil {
			return err
		}

		return nil
	}
}

func handleSubmit(logger *slog.Logger, e db.Executor) web.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if err := r.ParseForm(); err != nil {
			return web.RespondError("Invalid Form Values.", http.StatusBadRequest, err)
		}

		displayStart, err := time.Parse("2006-01-02", r.FormValue("start_date"))
		if err != nil {
			return web.RespondError("Invalid Form Date Value.", http.StatusBadRequest, err)
		}

		displayEnd, err := time.Parse("2006-01-02", r.FormValue("end_date"))
		if err != nil {
			return web.RespondError("Invalid Form Date Value.", http.StatusBadRequest, err)
		}

		content := markdown.ToHTML([]byte(r.FormValue("content")), nil, nil)

		n := resource.NewAnnouncement{
			Title:        r.FormValue("title"),
			Author:       r.FormValue("author"),
			Content:      string(content),
			DisplayStart: displayStart,
			DisplayEnd:   displayEnd,
		}

		if _, err := resource.InsertAnnouncement(r.Context(), e, n); err != nil {
			return err
		}

		logger.Info("added announcement", "title", n.Title, "author", n.Author)

		w.Header().Set("Location", "/admin")
		w.WriteHeader(http.StatusFound)

		return nil
	}
}
