package render

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/spa-stc/newsletter/server/profile"
)

// Echo Render Wrapper.
type Templates struct {
	templ *template.Template
}

func NewTemplates(profile *profile.Profile, files embed.FS) (*Templates, error) {
	templ, err := getTemplates(files)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing templates")
	}

	return &Templates{
		templ: templ,
	}, nil
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templ.ExecuteTemplate(w, name, data)
}

// Get templates, while preserving directory structure.
func getTemplates(templatesfs embed.FS) (*template.Template, error) {
	const fsPath = "templates"

	templates := template.New("")
	templates = templates.Funcs(
		template.FuncMap{
			"newkv": NewKv,
		},
	)
	err := fs.WalkDir(templatesfs, fsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			data, err := fs.ReadFile(templatesfs, path)
			if err != nil {
				return err
			}

			templates, err = templates.New(strings.ReplaceAll(path, fmt.Sprintf("%s/", fsPath), "")).Parse(string(data))
			if err != nil {
				return err
			}
		}

		return nil
	})

	return templates, err
}
