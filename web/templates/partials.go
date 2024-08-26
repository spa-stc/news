package templates

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"stpaulacademy.tech/newsletter/web/assets"
)

var ErrPartialNotFound = errors.New("partial not found")

type Partials struct {
	partials map[string]string
	assets   assets.Assets
}

func NewPartials(basedir string, assets assets.Assets) (*Partials, error) {
	partials := make(map[string]string)

	basedir = filepath.Clean(basedir)
	err := filepath.Walk(basedir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("filepath walk error: %w", err)
		}

		filename := strings.TrimPrefix(path, basedir+"/")

		if info.IsDir() {
			return nil
		}

		file, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", filename, err)
		}

		partials[filename] = string(file)

		return nil
	})
	if err != nil {
		return nil, err
	}

	p := &Partials{
		partials: partials,
		assets:   assets,
	}

	return p, nil
}

func (p *Partials) Render(name string, data interface{}) (string, error) {
	partial, ok := p.partials[name]
	if !ok {
		return "", ErrPartialNotFound
	}

	t, err := template.New(name).Parse(partial)
	if err != nil {
		return "", fmt.Errorf("error parsing partial %s: %w", name, err)
	}

	t.Funcs(map[string]any{
		"asset": p.assets.GetLink,
	})

	buf := bytes.NewBuffer(nil)
	if err := t.Execute(buf, data); err != nil {
		return "", fmt.Errorf("error executing partial %s: %w", name, err)
	}

	return buf.String(), nil
}