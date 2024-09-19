package templates

import (
	"html/template"

	"github.com/microcosm-cc/bluemonday"
)

type kvReturn struct {
	Key   interface{}
	Value interface{}
}

func keyValue(key, value any) kvReturn {
	return kvReturn{
		Key:   key,
		Value: value,
	}
}

func renderSafe() func(string) template.HTML {
	p := bluemonday.UGCPolicy()
	return func(s string) template.HTML {
		return template.HTML(p.Sanitize(s)) //nolint:gosec // Escaped using p.sanitize.
	}
}
