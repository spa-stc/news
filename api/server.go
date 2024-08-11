package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"stpaulacademy.tech/newsletter/db"
	"stpaulacademy.tech/newsletter/util/service"
	"stpaulacademy.tech/newsletter/web"
)

func NewServer(timegetter service.TimeGenerator, db db.Executor) http.Handler {
	r := chi.NewMux()

	r.Method(http.MethodGet, "/healthz", handleHealthZ())
	r.Method(http.MethodGet, "/week", handleGetWeek(db, timegetter))

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
