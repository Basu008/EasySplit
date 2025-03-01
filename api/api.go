package api

import (
	"github.com/Basu008/EasySplit.git/server/auth"
	"github.com/Basu008/EasySplit.git/server/config"
	"github.com/Basu008/EasySplit.git/server/validator"
	"github.com/gorilla/mux"
)

type API struct {
	Router     *Router
	MainRouter *mux.Router
	Config     *config.APIConfig
	TokenAuth  auth.TokenAuth
	Validator  *validator.Validator

	// App *app.App
}

type Options struct {
	MainRouter *mux.Router
	Config     *config.APIConfig
	TokenAuth  auth.TokenAuth
	Validator  *validator.Validator
}

// Router stores all the endpoints available for the server to respond.
type Router struct {
	Root    *mux.Router
	APIRoot *mux.Router
}

func NewAPI(opts *Options) *API {
	api := API{
		MainRouter: opts.MainRouter,
		Router:     &Router{},
		Config:     opts.Config,
		TokenAuth:  opts.TokenAuth,
		Validator:  opts.Validator,
	}
	api.setupRoutes()
	return &api
}

func (a *API) setupRoutes() {
	a.Router.Root = a.MainRouter
	a.Router.APIRoot = a.MainRouter.PathPrefix("/api").Subrouter()
	// a.InitRoutes()
}
