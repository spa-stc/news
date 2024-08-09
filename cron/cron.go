package cron

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"
)

type Job interface {
	Run(ctx context.Context) error
	Notifer() StatusNotifer
	Spec() string
}

type Service struct {
	cron *cron.Cron
}

func NewService() *Service {
	return &Service{
		cron: cron.New(cron.WithLocation(time.Local)),
	}
}

func (s *Service) Start() {
	go s.cron.Run()
}

func (s *Service) Stop() {
	ctx := s.cron.Stop()

	<-ctx.Done()
}

func (s *Service) AddJob(j Job) error {
	fn := func() {
		ctx := context.Background()
		start := time.Now()

		if err := j.Run(ctx); err != nil {
			err = j.Notifer().Failure(ctx, err)
			if err != nil {
				slog.Error("error sending cron notifer failure", "error", err)
			}
			return
		}

		err := j.Notifer().Success(ctx, time.Since(start))
		if err != nil {
			slog.Error("error sending cron notifer sucess", "error", err)
		}
	}

	if _, err := s.cron.AddFunc(j.Spec(), fn); err != nil {
		return fmt.Errorf("error adding job to cron instance: %w", err)
	}

	return nil
}
