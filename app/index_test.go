package app_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/neilotoole/slogt"
	"github.com/stretchr/testify/require"
	"stpaulacademy.tech/newsletter/app"
	"stpaulacademy.tech/newsletter/web"
)

func TestIndex(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	p, err := web.NewPublic("../public")
	require.NoError(t, err)

	s := app.NewServer(slogt.New(t), p.Assets(), p.RootAssets(), p.Templates())

	rec := httptest.NewRecorder()

	res, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
	require.NoError(t, err)

	s.ServeHTTP(rec, res)

	require.Equal(t, http.StatusOK, rec.Result().StatusCode)
}
