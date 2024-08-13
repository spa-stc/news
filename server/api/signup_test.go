package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"stpaulacademy.tech/newsletter/server/api"
	"stpaulacademy.tech/newsletter/util/testutil"
)

func sendSignupRequest(t *testing.T, srv http.Handler, r api.SignupRequest) *httptest.ResponseRecorder {
	req, err := http.NewRequest(http.MethodPost, "/signup", testutil.SerializeJSONReader(t, r))
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	srv.ServeHTTP(rr, req)

	return rr
}

func TestSignup(t *testing.T) {
	t.Parallel()

	ctx := testutil.Setup(t)

	t.Run("test_success", func(t *testing.T) {
		srv := NewTestAPI(ctx, t, time.Now())
		request := api.SignupRequest{
			Name:     "Bingo Ting",
			Password: "12345532",
			Email:    "bingo@e.edu",
		}

		rr := sendSignupRequest(t, srv, request)

		require.Equal(t, http.StatusCreated, rr.Result().StatusCode)
	})

	t.Run("test_conflict", func(t *testing.T) {
		srv := NewTestAPI(ctx, t, time.Now())
		request := api.SignupRequest{
			Name:     "Baaaaaaa",
			Password: "12345532w",
			Email:    "email@spa.edu",
		}

		rr := sendSignupRequest(t, srv, request)

		require.Equal(t, http.StatusConflict, rr.Result().StatusCode)
	})

	t.Run("test_nameshort", func(t *testing.T) {
		srv := NewTestAPI(ctx, t, time.Now())
		request := api.SignupRequest{
			Name:     "B",
			Password: "12345532",
			Email:    "bingo@e.edu",
		}

		rr := sendSignupRequest(t, srv, request)

		require.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
	})

	t.Run("test_namelong", func(t *testing.T) {
		srv := NewTestAPI(ctx, t, time.Now())
		request := api.SignupRequest{
			Name:     "Baaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			Password: "12345532",
			Email:    "bingo@e.edu",
		}

		rr := sendSignupRequest(t, srv, request)

		require.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
	})

	t.Run("test_notemail", func(t *testing.T) {
		srv := NewTestAPI(ctx, t, time.Now())
		request := api.SignupRequest{
			Name:     "Baaaaa",
			Password: "12345532",
			Email:    "bingo@",
		}

		rr := sendSignupRequest(t, srv, request)

		require.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
	})

	t.Run("test_passwordshort", func(t *testing.T) {
		srv := NewTestAPI(ctx, t, time.Now())
		request := api.SignupRequest{
			Name:     "Baaaaaa",
			Password: "123455",
			Email:    "bingo@e2.edu",
		}

		rr := sendSignupRequest(t, srv, request)

		require.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
	})

	t.Run("test_passwordlong", func(t *testing.T) {
		srv := NewTestAPI(ctx, t, time.Now())
		request := api.SignupRequest{
			Name:     "Baaaaaaa",
			Password: "12345532w23asdsdersdsdseeeeeeeeeersddsdsersdssd",
			Email:    "bingo@e3.edu",
		}

		rr := sendSignupRequest(t, srv, request)

		require.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
	})
}
