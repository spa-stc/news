package config

import (
	"os"
	"sync"
)

var (
	isDevelopment   bool      //nolint:gochecknoglobals // Added as a convienence feature, and is private.
	developmentOnce sync.Once //nolint:gochecknoglobals // Syncs the above.
)

func IsDevelopment() bool {
	developmentOnce.Do(func() {
		devel := os.Getenv("NEWSLETTER_DEVELOPMENT")

		var d bool
		switch devel {
		case "1":
			d = true
		case "true":
			d = true
		default:
			d = false
		}

		isDevelopment = d
	})

	return isDevelopment
}

func SetDevelopment(b bool) bool {
	developmentOnce.Do(func() {
		isDevelopment = b
	})

	return isDevelopment
}
