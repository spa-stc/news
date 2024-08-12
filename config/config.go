package config

import "errors"

type Config struct {
	Host       string
	Port       int
	DatbaseURL string

	SheetID  string
	SheetGID string
	IcalURL  string
}

func Validate(c Config) error {
	if c.Host == "" {
		return errors.New("missing host")
	}

	if c.Port == 0 {
		return errors.New("missing port")
	}

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
