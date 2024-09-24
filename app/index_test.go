package app_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/neilotoole/slogt"
	"github.com/stretchr/testify/require"
	"stpaulacademy.tech/newsletter/app"
	"stpaulacademy.tech/newsletter/config"
	"stpaulacademy.tech/newsletter/util/testutil"
	"stpaulacademy.tech/newsletter/web"
)

func TestIndex(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	p, err := web.NewPublic("../public")
	require.NoError(t, err)
	tx := testutil.TestTx(ctx, t)
	timegen := testutil.NewTestTimeGen(time.Date(2024, 12, 18, 0, 0, 0, 0, time.UTC))

	s := app.NewServer(slogt.New(t), p.Assets(), p.RootAssets(), p.Templates(), tx, timegen, config.Config{})

	rec := httptest.NewRecorder()

	res, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
	require.NoError(t, err)

	s.ServeHTTP(rec, res)

	require.Equal(t, http.StatusOK, rec.Result().StatusCode)
}
