package web

import (
	"errors"
	"log/slog"
	"net/http"

	"stpaulacademy.tech/newsletter/db"
)

type Error struct {
	Message    string
	StatusCode int
}

type Handler func(w http.ResponseWriter, r *http.Request) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h(w, r)
	if err != nil {
		webErr := mapWebError(err)
		w.WriteHeader(webErr.StatusCode)
		if _, err := w.Write([]byte(webErr.Message)); err != nil {
			slog.ErrorContext(r.Context(), "failed to write http error body", "err", err)
			return
		}
	}
}

func mapWebError(err error) Error {
	if errors.Is(err, db.ErrNotFound) {
		return Error{
			Message:    "Not Found",
			StatusCode: http.StatusNotFound,
		}
	}

	return Error{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal Server Error",
	}
}
