package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Basu008/EasySplit.git/api"
	"github.com/Basu008/EasySplit.git/app"
	"github.com/Basu008/EasySplit.git/server/auth"
	"github.com/Basu008/EasySplit.git/server/config"
	postgresStorage "github.com/Basu008/EasySplit.git/server/storage/postgres"
	"github.com/Basu008/EasySplit.git/server/validator"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
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

func (s *Server) StartServer() {
	fmt.Println("Setting up Server....")
	n := negroni.New()
	recovery := negroni.NewRecovery()
	n.Use(recovery)
	n.UseHandler(s.Router)
	s.httpServer = &http.Server{
		Handler:      n,
		Addr:         fmt.Sprintf("%s:%s", s.Config.ServerConfig.ListenAddr, s.Config.ServerConfig.Port),
		ReadTimeout:  s.Config.ServerConfig.ReadTimeout * time.Second,
		WriteTimeout: s.Config.ServerConfig.WriteTimeout * time.Second,
	}
	fmt.Printf("Server Started listening at %s:%s", s.Config.ServerConfig.ListenAddr, s.Config.ServerConfig.Port)
	go func() {
		err := s.httpServer.ListenAndServe()
		if err != nil {
			log.Fatal(err.Error())
			return
		}
	}()
}

func (s *Server) StopServer() {
	fmt.Println("Closing Postgres...")
	s.Postgres.Close()
	fmt.Println("Closed Postgres")
}
