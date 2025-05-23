package app

import (
	"github.com/Basu008/EasySplit.git/server/config"
	postgresStorage "github.com/Basu008/EasySplit.git/server/storage/postgres"
)

type Options struct {
	Postgres *postgresStorage.PostgresStorage
	Config   *config.APPConfig
}

type App struct {
	Postgres *postgresStorage.PostgresStorage
	Config   *config.APPConfig

	//Sevices
	User    User
	Friend  Friend
	Expense Expense
	Group   Group
	Payment Payment
}

func NewApp(opts *Options) *App {
	return &App{
		Postgres: opts.Postgres,
		Config:   opts.Config,
	}
}
