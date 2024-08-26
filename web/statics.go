package web

import (
	"net/http"

	"stpaulacademy.tech/newsletter/web/assets"
)

// /assets/{hash}.
func ServeStatics(a assets.Assets) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		hash := r.PathValue("hash")

		var h [16]byte
		copy(h[:], []byte(hash)[0:16])
		asset, ok := a.ByHash(h)
		if !ok {
			return RespondError("Asset Not Found.", http.StatusNotFound, nil)
		}

		w.Header().Set("Content-Type", asset.MimeType)
		if _, err := w.Write(asset.Content); err != nil {
			return err
		}

		return nil
	}
}

// /{filename}
func RootStatics(a assets.Assets) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		filename := r.PathValue("filename")

		hash, ok := a.ByName(filename)
		if !ok {
			return RespondError("Asset Not Found.", http.StatusNotFound, nil)
		}

		asset, ok := a.ByHash(hash)
		if !ok {
			return RespondError("Asset Not Found.", http.StatusNotFound, nil)
		}

		w.Header().Set("Content-Type", asset.MimeType)
		if _, err := w.Write(asset.Content); err != nil {
			return err
		}

		return nil
	}
}
