package assets

import (
	"crypto/md5" //nolint:gosec // Not used in cryptography.
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
	byName map[string][16]byte
	byHash map[[16]byte]Asset
}

func NewAssets() Assets {
	return Assets{
		byName: make(map[string][16]byte),
		byHash: make(map[[16]byte]Asset),
	}
}

func (a *Assets) ByName(name string) ([16]byte, bool) {
	hash, ok := a.byName[name]

	return hash, ok
}

func (a *Assets) ByHash(hash [16]byte) (Asset, bool) {
	asset, ok := a.byHash[hash]

	return asset, ok
}

func (a *Assets) AddBytes(name, mimeType string, content []byte) {
	hash := md5.Sum(content) //nolint:gosec // Not used in cryptography.

	asset := Asset{
		MimeType: mimeType,
		Content:  content,
	}

	a.byName[name] = hash
	a.byHash[hash] = asset
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
