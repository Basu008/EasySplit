package server

import (
	"net/http"

	"github.com/Basu008/EasySplit.git/api"
	"github.com/Basu008/EasySplit.git/app"
	"github.com/Basu008/EasySplit.git/server/auth"
	"github.com/Basu008/EasySplit.git/server/config"
	postgresStorage "github.com/Basu008/EasySplit.git/server/storage/postgres"
	"github.com/Basu008/EasySplit.git/server/validator"
	"github.com/gorilla/mux"
)

type Server struct {
	httpServer *http.Server
	Router     *mux.Router
	Config     *config.Config
	Postgres   postgresStorage.PostgresStorage

	API *api.API
}

func NewServer() *Server {
	c := config.GetConfig()
	ps := postgresStorage.NewPostgresStorage(&c.DatabaseConfig)
	r := mux.NewRouter()
	server := &Server{
		httpServer: &http.Server{},
		Router:     r,
		Config:     c,
		Postgres:   *ps,
	}
	server.API = api.NewAPI(&api.Options{
		MainRouter: r,
		Config:     &c.APIConfig,
		TokenAuth:  auth.NewTokenAuthentication(&c.TokenAuthConfig),
		Validator:  validator.NewValidation(),
	})
	server.API.App = app.NewApp(&app.Options{
		Postgres: ps,
		Config:   &c.APPConfig,
	})
	app.InitService(server.API.App)
	return server
}
