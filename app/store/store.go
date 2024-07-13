package store

import "github.com/spa-stc/newsletter/server/profile"

type Store struct {
	profile *profile.Profile
	db      Driver
}

func New(p *profile.Profile, driver Driver) *Store {
	return &Store{
		profile: p,
		db:      driver,
	}
}
