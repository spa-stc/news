package render

// Base Template Render Context, for use inside other, template specific render contexts.
type BaseContext struct {
	Title string

	Info SiteInfo
}

type SiteInfo struct {
	Env string
}
