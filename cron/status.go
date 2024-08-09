package cron

import (
	"context"
	"log/slog"
	"time"
)

// Allow a cron job to send it's end status to an external service.
type StatusNotifer interface {
	// Send a sucess message, with duration as the job's total runtime.
	Success(ctx context.Context, duration time.Duration) error
	// Send a failure message, with a provided error value.
	Failure(ctx context.Context, err error) error
}

type SlogStatusNotifer struct {
	logger  *slog.Logger
	jobName string
}

func NewSlogStatusNotifer(logger *slog.Logger, jobName string) *SlogStatusNotifer {
	if logger == nil {
		logger = slog.Default()
	}

	return &SlogStatusNotifer{
		logger:  logger,
		jobName: jobName,
	}
}

func (s *SlogStatusNotifer) Success(ctx context.Context, duration time.Duration) error {
	s.logger.InfoContext(ctx, "cron job completed", "job", s.jobName, "ms", duration.Milliseconds())

	return nil
}

func (s *SlogStatusNotifer) Failure(ctx context.Context, err error) error {
	s.logger.InfoContext(ctx, "cron job failed", "job", s.jobName, "error", err)

	return nil
}
