package web

import (
	"net/http"

	"stpaulacademy.tech/newsletter/web/assets"
)

const (
	immutableCacheHeader = "public, max-age=31536000, immutable"
)

// Serve assets to /assets/{hash}.
func ServeStatics(a *assets.Assets) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		hash := r.PathValue("hash")
		var byteHash [16]byte
		copy(byteHash[:], hash)

		asset, ok := a.ByHash(byteHash)
		if !ok {
			return Error{
				Code:    http.StatusNotFound,
				Message: "Asset Not Found.",
			}
		}

		w.Header().Set("Content-Type", asset.MimeType)
		w.Header().Set("Cache-Control", immutableCacheHeader)
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(asset.Content); err != nil {
			return err
		}

		return nil
	}
}
