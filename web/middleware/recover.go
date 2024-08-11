package middleware

import (
	"log/slog"
	"net/http"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler { //nolint:errorlint // Rvr is not an error value.
					panic(rvr)
				}

				w.WriteHeader(http.StatusInternalServerError)

				slog.ErrorContext(r.Context(),
					"recovered panic in http handler",
					"route", r.URL,
					"err", rvr)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
