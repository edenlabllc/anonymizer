package application

import (
	"ehealth-migration/pkg/config"
	"ehealth-migration/pkg/database"
	"flag"
	"github.com/rs/zerolog/log"
)

var (
	DEFAULT_LIMIT       = 100
	DEFAULT_CONCURRENCY = 10
)

type Flags struct {
	User  string
	Pass  string
	Name  string
	Host  string
	Port  string
	Tname string
	Limit int
	//Count       int
	Concurrency int
}

type Application struct {
	Config   config.Config
	DbClient *database.Database
	Flags    Flags
}

func (a *Application) SetFlags(flags Flags) {
	a.Flags = flags
}

var app Application

func init() {

	flags := Flags{}

	flag.StringVar(&flags.User, "user", "", "database user")
	flag.StringVar(&flags.Pass, "pass", "", "database pass")
	flag.StringVar(&flags.Name, "name", "", "database name")
	flag.StringVar(&flags.Host, "host", "localhost", "database host")
	flag.StringVar(&flags.Port, "port", "5432", "database port")
	flag.StringVar(&flags.Tname, "tname", "", "table name")
	flag.IntVar(&flags.Limit, "limit", DEFAULT_LIMIT, "limit items")

	//flag.IntVar(&flags.Count, "count", 0, "count items")
	flag.IntVar(&flags.Concurrency, "concurrency", DEFAULT_CONCURRENCY, "numbers goroutine")

	flag.Parse()

	// init Flags
	app.SetFlags(flags)
}

func NewApplication() *Application {
	// Connected DB
	db, err := database.GetDbClient(app.Flags.Name, app.Flags.User, app.Flags.Pass, app.Flags.Host, app.Flags.Port)
	if err != nil {
		log.Error().Msgf("Database errors: %s", err)
		return nil
	}
	//defer db.DB.Close()

	app.DbClient = db

	// Load Config
	config, err := config.InitYaml()
	if err != nil {
		log.Error().Msgf("Yaml parse errors: %s", err)
		return nil
	}
	app.Config = config

	return &app
}
