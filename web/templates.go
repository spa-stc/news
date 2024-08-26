package web

import (
	"net/http"

	"stpaulacademy.tech/newsletter/web/templates"
)

func RenderTemplate(w http.ResponseWriter, t templates.TemplateRenderer, name string, data templates.RenderData) error {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")

	return t.Render(name, data, w)
}
