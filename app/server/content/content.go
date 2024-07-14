package content

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/spa-stc/newsletter/public"
)

type Content struct {
	data map[string]string
}

func New() (*Content, error) {
	extensions := parser.Tables | parser.AutoHeadingIDs

	p := parser.NewWithExtensions(extensions)
	r := html.NewRenderer(html.RendererOptions{})

	data := make(map[string]string)
	err := fs.WalkDir(public.Content, "content", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			bytes, err := fs.ReadFile(public.Content, path)
			if err != nil {
				return err
			}

			html := markdown.ToHTML(bytes, p, r)
			name := strings.ReplaceAll(path, fmt.Sprintf("%s/", "content"), "")

			data[name] = string(html)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Content{data}, nil
}

func (c *Content) Get(title string) (string, error) {
	if data, ok := c.data[title]; ok {
		return data, nil
	}

	return "", fmt.Errorf("content: %s does not exist")
}
