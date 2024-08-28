package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"stpaulacademy.tech/newsletter/app"
	"stpaulacademy.tech/newsletter/config"
	"stpaulacademy.tech/newsletter/cron"
	"stpaulacademy.tech/newsletter/cron/jobs/daysfetch"
	"stpaulacademy.tech/newsletter/util/service"
)

var RootCMD = &cobra.Command{ //nolint:gochecknoglobals // Not state
	Use: "newsletter",
	RunE: func(_ *cobra.Command, _ []string) error {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		c := config.Config{
			DatbaseURL: viper.GetString("database_url"),
			IcalURL:    viper.GetString("ical_url"),
			SheetID:    viper.GetString("sheet_id"),
			SheetGID:   viper.GetString("sheet_gid"),
			Port:       viper.GetInt("port"),
		}

		err := config.Validate(c)
		if err != nil {
			return fmt.Errorf("error getting configuration: %w", err)
		}

		db, err := pgxpool.New(ctx, c.DatbaseURL)
		if err != nil {
			return fmt.Errorf("error connecting to database: %w", err)
		}

		timegen := &service.TimeGen{}

		cronservice := cron.NewService(logger)
		daygetter := daysfetch.NewGetter(c.SheetID, c.SheetGID, c.IcalURL)
		err = cronservice.AddJob(daysfetch.New(db, timegen, daygetter, cron.NewSlogStatusNotifer(logger, "days_fetch")))
		if err != nil {
			return fmt.Errorf("error adding cron job to runner: %w", err)
		}

		cronservice.Start()
		defer cronservice.Stop()

		app := app.NewServer()
		server := runServer(logger, app, fmt.Sprintf("0.0.0.0:%d", c.Port))
		defer func() {
			if err := server.Shutdown(ctx); err != nil {
				logger.Error("failed to stop http server", "error", err)
			}
		}()

		logger.InfoContext(ctx, "startup complete", "http_port", c.Port)

		sc := make(chan os.Signal, 1)
		signal.Notify(sc, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		<-sc

		return nil
	},
}

func init() { //nolint:gochecknoinits // Not state related
	RootCMD.PersistentFlags().String("database-url", "", "Location of postgres database")
	RootCMD.PersistentFlags().String("sheet-id", "", "ID of Google Sheet with XPeriod Information")
	RootCMD.PersistentFlags().String("sheet-gid", "", "Name of Google Sheet with XPeriod Info")
	RootCMD.PersistentFlags().String("ical-url", "", "Location of lunch calendar")
	RootCMD.PersistentFlags().Int("port", 3000, "What port to serve http over")

	err := viper.BindPFlag("database_url", RootCMD.PersistentFlags().Lookup("database-url"))
	if err != nil {
		panic(err)
	}

	err = viper.BindPFlag("ical_url", RootCMD.PersistentFlags().Lookup("ical-url"))
	if err != nil {
		panic(err)
	}

	err = viper.BindPFlag("sheet_gid", RootCMD.PersistentFlags().Lookup("sheet-gid"))
	if err != nil {
		panic(err)
	}

	err = viper.BindPFlag("sheet_id", RootCMD.PersistentFlags().Lookup("sheet-id"))
	if err != nil {
		panic(err)
	}

	err = viper.BindPFlag("port", RootCMD.PersistentFlags().Lookup("port"))
	if err != nil {
		panic(err)
	}

	viper.SetEnvPrefix("NEWSLETTER")
	viper.AutomaticEnv()
}

func runServer(logger *slog.Logger, h http.Handler, addr string) *http.Server {
	s := &http.Server{
		Addr:              addr,
		Handler:           h,
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second, //nolint:mnd // fine
		ReadHeaderTimeout: 2 * time.Second,  //nolint:mnd // fine
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("failed to start http server", "error", err)
		}
	}()

	return s
}
