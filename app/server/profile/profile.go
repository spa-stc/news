package profile

import (
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
}

// Get the profile, from viper, and validate.
func Get() (*Profile, error) {
	profile := Profile{}
	if err := viper.Unmarshal(&profile); err != nil {
		return nil, err
	}

	if profile.Dir == "" {
		if profile.Env == "production" {
			profile.Dir = "/var/opt/newsletter"
		} else {
			profile.Dir = "./tmp"
		}
	}

	return &profile, nil
}
