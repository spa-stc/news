package profile

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Main App Configuration, Usually Parsed From Environment.
type Profile struct {
	// Port where the app is to be served.
	Port string

	// One of development, production.
	Env string

	// Database location.
	DSN string

	// Data Dir.
	Dir string

	// Driver type.
	Driver string

	ICALURL string `mapstructure:"ical_url"`

	SheetName string `mapstructure:"sheet_name"`

	SheedID string `mapstructure:"sheet_id"`
}

// Get the profile, from viper, and validate.
func Get() (*Profile, error) {
	profile := Profile{}
	if err := viper.Unmarshal(&profile); err != nil {
		return nil, err
	}

	if profile.Dir == "" {
		if profile.Env == "production" {
			profile.Dir = "/var/lib/newsletter"
		} else {
			profile.Dir = "./tmp"
		}
	}

	if profile.ICALURL == "" {
		return nil, errors.New("missing required field ical url")
	}

	if profile.SheedID == "" {
		return nil, errors.New("missing required field sheet id")
	}

	if profile.SheetName == "" {
		return nil, errors.New("missing required field sheet name")
	}

	return &profile, nil
}
