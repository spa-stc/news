package api_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"stpaulacademy.tech/newsletter/api"
	"stpaulacademy.tech/newsletter/util/testutil"
)

func NewTestAPI(ctx context.Context, t *testing.T, now time.Time) http.Handler {
	tx := testutil.TestTx(ctx, t)

	timegen := testutil.NewTestTimeGen(now)

	server := api.NewServer(timegen, tx)

	return server
}
