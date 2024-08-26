package web

import (
	"errors"
	"net/http"
)

type Handler func(w http.ResponseWriter, r *http.Request) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		var webError Error
		if errors.As(err, &webError) {
			w.WriteHeader(webError.Code)
			if _, err := w.Write([]byte(webError.Message)); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}

		w.WriteHeader(http.StatusInternalServerError)

		if _, err := w.Write([]byte("Internal Server Error.")); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
