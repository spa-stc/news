package testutil

import (
	"context"
	"log/slog"
	"testing"

	"github.com/neilotoole/slogt"
)

func Setup(t *testing.T) context.Context {
	logger := slogt.New(t)
	slog.SetDefault(logger)

	return context.Background()
}
