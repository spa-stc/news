package workers

import (
	"log/slog"

	"github.com/go-co-op/gocron/v2"
	"github.com/spa-stc/newsletter/server/profile"
	"github.com/spa-stc/newsletter/store"
	"github.com/spa-stc/newsletter/workers/daily"
)

type Service struct {
	p     *profile.Profile
	s     gocron.Scheduler
	store *store.Store
}

func New(p *profile.Profile, db *store.Store) (*Service, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	time := gocron.NewAtTime(0, 0, 0)

	s.NewJob(gocron.DailyJob(1, gocron.NewAtTimes(time, gocron.NewTask(daily.Run, p, db))

	service := &Service{
		p:     p,
		s:     s,
		store: db,
	}

	return service, nil
}

func (s *Service) Start() {
	s.s.Start()
}

func (s *Service) Shutdown() {
	if err := s.s.Shutdown(); err != nil {
		slog.Error("failed to stop cron scheduler", "error", err)
	}
}
