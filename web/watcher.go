package web

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/bep/debounce"
	"github.com/fsnotify/fsnotify"
)

func NewWatchedPublic(logger *slog.Logger, dir string) (*Public, error) {
	dir = filepath.Clean(dir)

	public, err := render(dir)
	if err != nil {
		return nil, err
	}

	p := *public

	rebuildFunc := func() {
		logger.Info("rebuilding public dir")
		public, err = render(dir)
		if err != nil {
			logger.Error("failed to build publics", "error", err)
		}

		p = *public
	}

	debouncer := debounce.New(100 * time.Millisecond)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("error creating file watcher: %w", err)
	}

	paths := []string{"templates", "templates/partials", "assets"}
	for _, path := range paths {
		if err = watcher.Add(filepath.Join(dir, path)); err != nil {
			return nil, fmt.Errorf("error watching public dir: %w", err)
		}
	}

	go func() {
		for range watcher.Events {
			debouncer(rebuildFunc)
		}
	}()

	return &p, nil
}
