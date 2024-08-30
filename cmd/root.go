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
	"stpaulacademy.tech/newsletter/web"
)

var RootCMD = &cobra.Command{ //nolint:gochecknoglobals // Not state
	Use: "newsletter",
	RunE: func(_ *cobra.Command, _ []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		c := config.Config{
			DatbaseURL:  viper.GetString("database_url"),
			IcalURL:     viper.GetString("ical_url"),
			SheetID:     viper.GetString("sheet_id"),
			SheetGID:    viper.GetString("sheet_gid"),
			PublicDir:   viper.GetString("public_dir"),
			Port:        viper.GetInt("port"),
			Development: viper.GetBool("development"),
		}

		err := config.Validate(c)
		if err != nil {
			return fmt.Errorf("error getting configuration: %w", err)
		}

		logLevel := slog.LevelInfo.Level()
		if c.Development {
			logLevel = slog.LevelDebug.Level()
		}

		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		}))

		public, err := web.NewPublic(c.PublicDir)
		if c.Development {
			public, err = web.NewWatchedPublic(logger, c.PublicDir)
		}
		if err != nil {
			return err
		}

		db, err := pgxpool.New(ctx, c.DatbaseURL)
		if err != nil {
			return fmt.Errorf("error connecting to database: %w", err)
		}

		timegen := &service.TimeGen{}

		cronservice := cron.NewService(logger)
		daygetter := daysfetch.NewGetter(c.SheetID, c.SheetGID, c.IcalURL)
		err = cronservice.AddJob(daysfetch.New(db,
			timegen,
			daygetter,
			cron.NewSlogStatusNotifer(logger, "days_fetch")))
		if err != nil {
			return fmt.Errorf("error adding cron job to runner: %w", err)
		}

		cronservice.Start()
		defer cronservice.Stop()

		app := app.NewServer(logger, public.Assets(), public.Templates())
		server := runServer(logger, app, fmt.Sprintf("0.0.0.0:%d", c.Port))
		defer func() {
			if err := server.Shutdown(ctx); err != nil {
				logger.Error("failed to stop http server", "error", err)
			}
		}()

		logger.InfoContext(ctx, "startup complete", "http_port", c.Port, "development", c.Development)

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
	RootCMD.PersistentFlags().String("public-dir", "", "Location of templates and static files.")
	RootCMD.PersistentFlags().Bool("development", false, "Enable development mode.")

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

	err = viper.BindPFlag("public_dir", RootCMD.PersistentFlags().Lookup("public-dir"))
	if err != nil {
		panic(err)
	}

	err = viper.BindPFlag("development", RootCMD.PersistentFlags().Lookup("development"))
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
