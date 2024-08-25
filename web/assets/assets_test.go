package assets_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"stpaulacademy.tech/newsletter/web/assets"
)

func TestAssetLoad(t *testing.T) {
	t.Parallel()
	a := assets.NewAssets()
	err := a.AddDir("./fixtures")
	require.NoError(t, err)

	hash, ok := a.ByName("hi.css")
	require.True(t, ok)

	asset, ok := a.ByHash(hash)
	require.True(t, ok)

	require.Equal(t, "text/css", asset.MimeType)
}
