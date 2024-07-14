package render

// Base Template Render Context, for use inside other, template specific render contexts.
type BaseContext struct {
	Info SiteInfo

	Data any
}

type SiteInfo struct {
	Env string
}
