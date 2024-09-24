package config

import "errors"

type Config struct {
	DatabaseURL string

	SheetID  string
	SheetGID string
	IcalURL  string

	PublicDir string

	Port int

	Development bool

	AdminUsername string
	AdminPassword string
}

func Validate(c Config) error {
	if c.DatabaseURL == "" {
		return errors.New("missing database_url")
	}

	if c.SheetID == "" {
		return errors.New("missing sheet_id")
	}

	if c.SheetGID == "" {
		return errors.New("missing sheet_gid")
	}

	if c.IcalURL == "" {
		return errors.New("missing ical_url")
	}
	if c.PublicDir == "" {
		return errors.New("missing public_dir")
	}

	if c.AdminUsername == "" {
		return errors.New("missing admin_username")
	}

	if c.AdminPassword == "" {
		return errors.New("missing admin_password")
	}

	return nil
}
