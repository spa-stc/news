package templates

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"stpaulacademy.tech/newsletter/web/assets"
)

var ErrTemplateNotFound = errors.New("template not found")

type TemplateRenderer struct {
	Assets       *assets.Assets
	Partials     *Partials
	BaseTemplate string
	Templates    map[string]string
}

type RenderData struct {
	Data interface{}
}

func New(baseDir string, rootPath string, assets *assets.Assets, partials *Partials) (*TemplateRenderer, error) {
	baseDir = filepath.Clean(baseDir)
	rootPath = filepath.Clean(rootPath)

	templates := make(map[string]string)
	err := filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walkdir error: %w", err)
		}

		if d.IsDir() {
			return nil
		}

		filename := strings.TrimPrefix(path, baseDir+"/")

		file, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("file read error: %w", err)
		}

		templates[filename] = string(file)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error generating templates: %w", err)
	}

	rootTemplate, err := os.ReadFile(rootPath)
	if err != nil {
		return nil, fmt.Errorf("error reading root template: %w", err)
	}

	t := &TemplateRenderer{
		Assets:       assets,
		Partials:     partials,
		BaseTemplate: string(rootTemplate),
		Templates:    templates,
	}

	return t, nil
}

func (t *TemplateRenderer) Render(name string, data RenderData, w io.Writer) error {
	templateString, ok := t.Templates[name]
	if !ok {
		return ErrTemplateNotFound
	}

	templ := template.New("index")
	templ.Funcs(map[string]any{
		"asset":     t.Assets.GetLink,
		"partial":   t.Partials.Render,
		"key_value": keyValue,
		"sanitize":  renderSafe(),
	})

	templ, err := templ.New("root.html").Parse(t.BaseTemplate)
	if err != nil {
		return fmt.Errorf("error parsing base template: %w", err)
	}

	templ, err = templ.New(name).Parse(templateString)
	if err != nil {
		return fmt.Errorf("error parsing template %s: %w", name, err)
	}

	if err := templ.ExecuteTemplate(w, name, data); err != nil {
		return fmt.Errorf("template render error: %w", err)
	}

	return nil
}
