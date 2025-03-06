package app

import (
	"net/http"

	"github.com/Basu008/EasySplit.git/server/config"
	postgresStorage "github.com/Basu008/EasySplit.git/server/storage/postgres"
)

type Options struct {
	Postgres *postgresStorage.PostgresStorage
	Config   *config.APPConfig
}

type App struct {
	Postgres   *postgresStorage.PostgresStorage
	Config     *config.APPConfig
	HttpClient http.Client

	//Sevices
	User   User
	Friend Friend
}

func NewApp(opts *Options) *App {
	return &App{
		Postgres: opts.Postgres,
		Config:   opts.Config,
	}
}
