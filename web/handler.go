package web

import (
	"errors"
	"log/slog"
	"net/http"

	"stpaulacademy.tech/newsletter/db"
)

type Error struct {
	Message    string `json:"message"`
	Data       any    `json:"data"`
	StatusCode int    `json:"-"`
}

func (e Error) Error() string {
	return e.Message
}

type Handler func(w http.ResponseWriter, r *http.Request) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h(w, r)
	if err != nil {
		webErr := mapWebError(err)
		if err := Encode(w, webErr.StatusCode, &webErr); err != nil {
			slog.Error("failed to send web error value", "err", err)
		}
		return
	}
}

func mapWebError(err error) Error {
	if errors.Is(err, db.ErrNotFound) {
		return Error{
			Message:    "Not Found",
			StatusCode: http.StatusNotFound,
		}
	}

	var webErr Error
	if errors.As(err, &webErr) {
		return webErr
	}

	return Error{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal Server Error",
	}
}
