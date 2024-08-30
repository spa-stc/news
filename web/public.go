package web

import (
	"fmt"
	"path/filepath"

	"stpaulacademy.tech/newsletter/web/assets"
	"stpaulacademy.tech/newsletter/web/templates"
)

type Public struct {
	templates templates.TemplateRenderer
	assets    assets.Assets
}

func NewPublic(dir string) (*Public, error) {
	return render(dir)
}

func render(dir string) (*Public, error) {
	dir = filepath.Clean(dir)

	assets := assets.NewAssets()

	err := assets.AddDir(filepath.Join(dir, "assets"))
	if err != nil {
		return nil, fmt.Errorf("error building assets: %w", err)
	}

	partials, err := templates.NewPartials(filepath.Join(dir, "templates/partials"), assets)
	if err != nil {
		return nil, fmt.Errorf("error building partials: %w", err)
	}

	templatesdir := filepath.Join(dir, "templates")
	templ, err := templates.New(templatesdir, filepath.Join(templatesdir, "root.html"), assets, partials)
	if err != nil {
		return nil, fmt.Errorf("error building partials: %w", err)
	}

	return &Public{
		templates: *templ,
		assets:    *assets,
	}, nil
}

func (p *Public) Templates() *templates.TemplateRenderer {
	return &p.templates
}

func (p *Public) Assets() *assets.Assets {
	return &p.assets
}
