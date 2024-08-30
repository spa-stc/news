package web

import (
	"errors"
	"log/slog"
	"net/http"
)

type Handler func(w http.ResponseWriter, r *http.Request) error

type HandlerWrapper struct {
	logger *slog.Logger
}

func NewHandlerWrapper(logger *slog.Logger) HandlerWrapper {
	return HandlerWrapper{
		logger: logger,
	}
}

func (w HandlerWrapper) Wrap(h Handler) WrappedHandler {
	return WrappedHandler{
		logger: w.logger,
		h:      h,
	}
}

type WrappedHandler struct {
	h      Handler
	logger *slog.Logger
}

func (h WrappedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.h(w, r); err != nil {
		var webError Error
		if errors.As(err, &webError) {
			w.WriteHeader(webError.Code)
			if _, err = w.Write([]byte(webError.Message)); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}

		h.logger.Error("unhandled http error", "error", err)

		w.WriteHeader(http.StatusInternalServerError)

		if _, err = w.Write([]byte("Internal Server Error.")); err != nil {
			h.logger.Error("failed to send bytes", "error", err)
		}
	}
}
