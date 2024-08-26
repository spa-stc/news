package templates_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	"stpaulacademy.tech/newsletter/web/templates"
)

func TestTemplateRender(t *testing.T) {
	t.Parallel()
	renderer, err := templates.New("./fixtures/templates", "./fixtures/templates/root.html", nil, nil)

	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = renderer.Render("yo.html", templates.RenderData{}, buf)

	require.NoError(t, err)
}
