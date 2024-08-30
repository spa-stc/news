package assets

import (
	"crypto/md5" //nolint:gosec // Not used in cryptography.
	"encoding/base32"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var ErrAssetNotFound = errors.New("asset not found")

type Asset struct {
	MimeType string
	Content  []byte
}

type Assets struct {
	byName map[string]string
	byHash map[string]Asset
}

func NewAssets() *Assets {
	return &Assets{
		byName: make(map[string]string),
		byHash: make(map[string]Asset),
	}
}

func (a *Assets) ByName(name string) (string, bool) {
	hash, ok := a.byName[name]

	return hash, ok
}

func (a *Assets) ByHash(hash string) (Asset, bool) {
	asset, ok := a.byHash[hash]

	return asset, ok
}

func (a *Assets) AddBytes(name, mimeType string, content []byte) {
	hash := md5.Sum(content) //nolint:gosec // Not used in cryptography.

	h := base32.HexEncoding.EncodeToString(hash[:])

	asset := Asset{
		MimeType: mimeType,
		Content:  content,
	}

	a.byName[name] = h
	a.byHash[h] = asset
}

func (a *Assets) AddDir(basepath string) error {
	basepath = filepath.Clean(basepath)

	err := filepath.WalkDir(basepath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("dir walk error: %w", err)
		}

		// Skip Directories.
		if d.IsDir() {
			return nil
		}

		filename := strings.TrimPrefix(path, basepath+"/")

		_, ext, found := strings.Cut(filename, ".")
		if !found {
			return fmt.Errorf("filename %s in invalid format (missing .type suffix)", filename)
		}

		bytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file: %w", err)
		}

		a.AddBytes(filename, getFiletype(ext, bytes), bytes)

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *Assets) GetLink(name string) (string, error) {
	hash, ok := a.byName[name]
	if !ok {
		return "", ErrAssetNotFound
	}

	return fmt.Sprintf("/assets/%s", hash), nil
}

func getFiletype(extension string, _ []byte) string {
	switch extension {
	case "txt":
		return "text/plain"
	case "css":
		return "text/css"
	case "js":
		return "application/javascript"
	}

	return "text/plain"
}
