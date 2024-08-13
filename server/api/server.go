package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"stpaulacademy.tech/newsletter/db"
	"stpaulacademy.tech/newsletter/util/service"
	"stpaulacademy.tech/newsletter/web"
	"stpaulacademy.tech/newsletter/web/middleware"
)

func NewServer(timegetter service.TimeGenerator, db db.Executor) http.Handler {
	v := web.NewValidator()

	r := chi.NewMux()
	r.Use(middleware.Recover)

	r.Method(http.MethodGet, "/healthz", handleHealthZ())
	r.Method(http.MethodGet, "/week", handleGetWeek(db, timegetter))
	r.Method(http.MethodPost, "/signup", handleSignup(db, v))

	return r
}

func handleHealthZ() web.Handler {
	return func(w http.ResponseWriter, _ *http.Request) error {
		if _, err := w.Write([]byte("Service Ready.")); err != nil {
			return err
		}

		return nil
	}
}
