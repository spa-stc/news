package web

import (
	"fmt"
	"path/filepath"

	"stpaulacademy.tech/newsletter/web/assets"
	"stpaulacademy.tech/newsletter/web/templates"
)

type Public struct {
	templates  templates.TemplateRenderer
	assets     assets.Assets
	rootAssets assets.Assets
}

func NewPublic(dir string) (*Public, error) {
	return render(dir)
}

func render(dir string) (*Public, error) {
	dir = filepath.Clean(dir)

	a := assets.NewAssets()

	err := a.AddDir(filepath.Join(dir, "assets"))
	if err != nil {
		return nil, fmt.Errorf("error building assets: %w", err)
	}

	rootAssets := assets.NewAssets()

	err = rootAssets.AddDir(filepath.Join(dir, "assets/root"))
	if err != nil {
		return nil, fmt.Errorf("error building root assets: %w", err)
	}

	partials, err := templates.NewPartials(filepath.Join(dir, "templates/partials"), a)
	if err != nil {
		return nil, fmt.Errorf("error building partials: %w", err)
	}

	templatesdir := filepath.Join(dir, "templates")
	templ, err := templates.New(templatesdir, filepath.Join(templatesdir, "root.html"), a, partials)
	if err != nil {
		return nil, fmt.Errorf("error building partials: %w", err)
	}

	return &Public{
		templates:  *templ,
		assets:     *a,
		rootAssets: *rootAssets,
	}, nil
}

func (p *Public) Templates() *templates.TemplateRenderer {
	return &p.templates
}

func (p *Public) Assets() *assets.Assets {
	return &p.assets
}

func (p *Public) RootAssets() *assets.Assets {
	return &p.rootAssets
}
