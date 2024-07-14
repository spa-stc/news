package public

import "embed"

//go:embed templates/**/*.tmpl.html
var Templates embed.FS

//go:embed static/*.*
var Static embed.FS

//go:embed content/*.md
var Content embed.FS
