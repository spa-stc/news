package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spa-stc/newsletter/public"
	"github.com/spa-stc/newsletter/server"
	"github.com/spa-stc/newsletter/server/profile"
	"github.com/spa-stc/newsletter/server/render"
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

		server := server.New(ctx, config, tmpl)

		server.Start(ctx)
		slog.Info("server started", "port", config.Port)

		// Listen for the graceful shutdown signal.
		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			<-c

			server.Shutdown(ctx)
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
	viper.SetEnvPrefix("NEWSLETTER")
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
