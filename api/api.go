package api

import (
	"net/http"

	"github.com/Basu008/EasySplit.git/app"
	"github.com/Basu008/EasySplit.git/server/auth"
	"github.com/Basu008/EasySplit.git/server/config"
	"github.com/Basu008/EasySplit.git/server/handler"
	"github.com/Basu008/EasySplit.git/server/validator"
	"github.com/gorilla/mux"
)

type API struct {
	Router     *Router
	MainRouter *mux.Router
	Config     *config.APIConfig
	TokenAuth  auth.TokenAuthentication
	Validator  *validator.Validator

	App *app.App
}

type Options struct {
	MainRouter *mux.Router
	Config     *config.APIConfig
	TokenAuth  auth.TokenAuthentication
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
	a.InitRoutes()
}

func (a *API) requestHandler(h func(c *handler.RequestContext, w http.ResponseWriter, r *http.Request)) http.Handler {
	return &handler.Request{
		HandlerFunc: h,
		AuthFunc:    a.TokenAuth,
		IsLoggedIn:  false,
		IsSudoUser:  false,
	}
}

func (a *API) requestWithAuthHandler(h func(c *handler.RequestContext, w http.ResponseWriter, r *http.Request)) http.Handler {
	return &handler.Request{
		HandlerFunc: h,
		AuthFunc:    a.TokenAuth,
		IsLoggedIn:  true,
		IsSudoUser:  false,
	}
}
