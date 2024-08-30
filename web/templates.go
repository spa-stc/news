package web

import (
	"net/http"

	"stpaulacademy.tech/newsletter/web/templates"
)

type TemplateCachePolicy string

const (
	TemplateCachePolicyPublic TemplateCachePolicy = "public, max-age=3600, stale-if-error=60"

	TemplateCachePolicyPrivate = "no-store"
)

func RenderTemplate(w http.ResponseWriter,
	t *templates.TemplateRenderer,
	name string,
	cachePolicy TemplateCachePolicy,
	data templates.RenderData,
) error {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", string(cachePolicy))

	return t.Render(name, data, w)
}
