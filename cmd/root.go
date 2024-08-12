package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"stpaulacademy.tech/newsletter/config"
	"stpaulacademy.tech/newsletter/cron"
	"stpaulacademy.tech/newsletter/cron/jobs/daysfetch"
	"stpaulacademy.tech/newsletter/server"
	"stpaulacademy.tech/newsletter/util/service"
)

var RootCMD = &cobra.Command{ //nolint:gochecknoglobals // Not state
	Use: "newsletter",
	RunE: func(_ *cobra.Command, _ []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		c := config.Config{
			Host:       viper.GetString("host"),
			Port:       viper.GetInt("port"),
			DatbaseURL: viper.GetString("database_url"),
			IcalURL:    viper.GetString("ical_url"),
			SheetID:    viper.GetString("sheet_id"),
			SheetGID:   viper.GetString("sheet_gid"),
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

		cronservice := cron.NewService()
		daygetter := daysfetch.NewGetter(c.SheetID, c.SheetGID, c.IcalURL)
		err = cronservice.AddJob(daysfetch.New(db, timegen, daygetter, cron.NewSlogStatusNotifer(nil, "days_fetch")))
		if err != nil {
			return fmt.Errorf("error adding cron job to runner: %w", err)
		}

		cronservice.Start()
		defer cronservice.Stop()

		s := server.New(c, db)
		s.Run()
		defer s.Stop(ctx)

		sc := make(chan os.Signal, 1)
		signal.Notify(sc, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		<-sc

		return nil
	},
}

func init() { //nolint:gochecknoinits // Not state related
	viper.SetDefault("host", "localhost")
	viper.SetDefault("port", "3000")

	RootCMD.PersistentFlags().String("host", "localhost", "Where to serve app")
	RootCMD.PersistentFlags().String("database-url", "", "Location of postgres database")
	RootCMD.PersistentFlags().Int("port", 3000, "Where to serve app") //nolint:mnd // ok
	RootCMD.PersistentFlags().String("sheet-id", "", "ID of Google Sheet with XPeriod Information")
	RootCMD.PersistentFlags().String("sheet-gid", "", "Name of Google Sheet with XPeriod Info")
	RootCMD.PersistentFlags().String("ical-url", "", "Location of lunch calendar")

	err := viper.BindPFlag("host", RootCMD.PersistentFlags().Lookup("host"))
	if err != nil {
		panic(err)
	}

	err = viper.BindPFlag("port", RootCMD.PersistentFlags().Lookup("port"))
	if err != nil {
		panic(err)
	}

	err = viper.BindPFlag("database_url", RootCMD.PersistentFlags().Lookup("database-url"))
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

	viper.SetEnvPrefix("NEWSLETTER")
	viper.AutomaticEnv()
}
