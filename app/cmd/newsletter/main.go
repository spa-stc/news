package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spa-stc/newsletter/public"
	"github.com/spa-stc/newsletter/server"
	"github.com/spa-stc/newsletter/server/profile"
	"github.com/spa-stc/newsletter/server/render"
	"github.com/spa-stc/newsletter/store"
	"github.com/spa-stc/newsletter/store/db"
	"github.com/spa-stc/newsletter/workers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func exit() {
	os.Exit(1)
}

var rootCmd = &cobra.Command{
	Use:   "newsletter",
	Short: "Provides Day Related Information to SPA Students",
	Run: func(_ *cobra.Command, _ []string) {
		viper.AutomaticEnv()

		config, err := profile.Get()
		if err != nil {
			slog.Error("failed to get profile", "error", err)
			exit()
		}

		tmpl, err := render.NewTemplates(config, public.Templates)
		if err != nil {
			slog.Error("failed to parse template fs", "error", err)
			exit()
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := setupDataDir(config.Dir); err != nil {
			slog.Error("failed to setup data dir", "error", err)
			exit()
		}

		db, err := db.GetDriver(config)
		if err != nil {
			slog.Error("error initializing database", "error", err)
			exit()
		}

		if err := db.Migrate(ctx); err != nil {
			slog.Error("error migrating database", "error", err)
			exit()
		}

		store := store.New(config, db)

		server := server.New(ctx, config, tmpl, store)

		worker, err := workers.New(config, store)
		if err != nil {
			slog.Error("error intializing workers", "error", err)
			exit()
		}

		server.Start(ctx)
		slog.Info("server started", "port", config.Port)
		worker.Start()

		// Listen for the graceful shutdown signal.
		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			<-c

			server.Shutdown(ctx)
			worker.Shutdown()
			cancel()
		}()

		<-ctx.Done()
	},
}

func init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	viper.SetDefault("port", "3000")
	viper.SetDefault("env", "production")
	viper.SetDefault("dsn", "data.db")
	viper.SetDefault("driver", "sqlite")
	viper.SetDefault("ical_url", "")
	viper.SetDefault("sheet_id", "")
	viper.SetDefault("sheet_name", "")
	viper.SetEnvPrefix("NEWSLETTER")
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func setupDataDir(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}

		if err := os.Mkdir(path, 0o770); err != nil {
			return err
		}
	}

	return nil
}
