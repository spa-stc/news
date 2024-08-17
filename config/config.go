package config

import "errors"

type Config struct {
	DatbaseURL string

	SheetID  string
	SheetGID string
	IcalURL  string
}

func Validate(c Config) error {
	if c.DatbaseURL == "" {
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

	return nil
}
