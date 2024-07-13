package db

import (
	"errors"

	"github.com/spa-stc/newsletter/server/profile"
	"github.com/spa-stc/newsletter/store"
	"github.com/spa-stc/newsletter/store/db/sqlite"
)

func GetDriver(p *profile.Profile) (store.Driver, error) {
	var driver store.Driver
	var err error

	switch p.Driver {
	case "sqlite":

		driver, err = sqlite.NewDB(p)
	default:
		return nil, errors.New("unknown db driver")
	}

	if err != nil {
		return nil, err
	}

	return driver, nil
}
