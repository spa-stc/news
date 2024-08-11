package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"stpaulacademy.tech/newsletter/util/testutil"
	"stpaulacademy.tech/newsletter/web/middleware"
)

func TestRecoverer(t *testing.T) {
	t.Parallel()
	_ = testutil.Setup(t)

	t.Run("test_recover_panic", func(t *testing.T) {
		r := chi.NewMux()

		r.Use(middleware.Recover)

		r.Get("/panic", func(_ http.ResponseWriter, _ *http.Request) {
			panic("oh no")
		})

		res, err := http.NewRequest(http.MethodGet, "/panic", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, res)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("test_abort_handler_norecover", func(t *testing.T) {
		defer func() {
			if rvr := recover(); rvr == nil {
				require.Fail(t, "panic must not be recovered")
			}
		}()

		r := chi.NewMux()

		r.Use(middleware.Recover)

		r.Get("/panic", func(_ http.ResponseWriter, _ *http.Request) {
			panic(http.ErrAbortHandler)
		})

		res, err := http.NewRequest(http.MethodGet, "/panic", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, res)
	})
}
