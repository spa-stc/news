package main

import (
	"flag"
	"log"

	"github.com/facebookgo/flagenv"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	// Our pocketbase migrations.
	_ "github.com/spa-stc/news/migrations"
)

var (
	Production bool
)

func init() {
	flag.BoolVar(&Production, "production", false, "production mode")
}

func main() {
	flagenv.Parse()

	app := pocketbase.New()

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Admin UI while not in prod.
		Automigrate: !Production,
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
