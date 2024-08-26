package templates_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"stpaulacademy.tech/newsletter/web/assets"
	"stpaulacademy.tech/newsletter/web/templates"
)

func TestPartials(t *testing.T) {
	t.Parallel()
	assets := assets.NewAssets()

	partials, err := templates.NewPartials("./fixtures/partials", assets)
	require.NoError(t, err)

	_, err = partials.Render("p.html", struct{ Hi string }{Hi: "Hi"})
	require.NoError(t, err)
}
