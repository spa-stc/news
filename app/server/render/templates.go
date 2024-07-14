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
	templates map[string]*template.Template
	info      SiteInfo
}

func NewTemplates(profile *profile.Profile, files embed.FS) (*Templates, error) {
	templ, err := getTemplates(files)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing templates")
	}

	info := getSiteInfo(profile)

	return &Templates{
		templates: templ,
		info:      info,
	}, nil
}

func getSiteInfo(p *profile.Profile) SiteInfo {
	return SiteInfo{
		Env: p.Env,
	}
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	injectedData := BaseContext{
		Info: t.info,

		Data: data,
	}

	if templ, ok := t.templates[name]; ok {
		return templ.Execute(w, &injectedData)
	}

	return errors.Errorf("template does not exist: %s", name)
}

// Get templates, while preserving directory structure.
func getTemplates(templatesfs embed.FS) (map[string]*template.Template, error) {
	const fsPath = "templates"

	layoutData, err := fs.ReadFile(templatesfs, "templates/layouts/main.tmpl.html")
	if err != nil {
		return nil, err
	}

	templs := make(map[string]*template.Template)
	err = fs.WalkDir(templatesfs, fsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			data, err := fs.ReadFile(templatesfs, path)
			if err != nil {
				return err
			}

			templName := strings.ReplaceAll(path, fmt.Sprintf("%s/", fsPath), "")
			templates := template.New(templName).Funcs(
				template.FuncMap{
					"newkv": NewKv,
				},
			)

			templates, err = templates.Parse(string(layoutData))
			if err != nil {
				return err
			}

			templates, err = templates.Parse(string(data))
			if err != nil {
				return err
			}

			templs[templName] = templates

		}

		return nil
	})

	return templs, err
}
