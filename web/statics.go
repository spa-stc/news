package web

import (
	"errors"
	"net/http"
	"strings"

	"stpaulacademy.tech/newsletter/web/assets"
)

const (
	immutableCacheHeader = "public, max-age=31536000, immutable"
	wellKnownCacheHeader = "public, max-age=86400, immutable"
)

// Serve assets to /assets/{hash}.
func ServeStatics(a *assets.Assets) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		hash := r.PathValue("hash")

		asset, ok := a.ByHash(hash)
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

// Custom not found handler that serves static paths like favicon.ico.
func ServeRootStatics(a *assets.Assets) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		hash, ok := a.ByName(strings.TrimPrefix(r.URL.Path, "/"))
		if !ok {
			return Error{
				Code:    http.StatusNotFound,
				Message: "Not Found.",
			}
		}

		if r.Header.Get("If-None-Match") == hash {
			w.Header().Set("Cache-Control", wellKnownCacheHeader)
			w.WriteHeader(http.StatusNotModified)
			return nil
		}

		asset, ok := a.ByHash(hash)
		if !ok {
			return errors.New("missing asset from byHash that exists in byName")
		}

		w.Header().Set("Content-Type", asset.MimeType)
		w.Header().Set("Cache-Control", wellKnownCacheHeader)
		w.Header().Set("Etag", hash)
		if _, err := w.Write(asset.Content); err != nil {
			return err
		}
		return nil
	}
}
