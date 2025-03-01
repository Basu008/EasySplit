package server

import (
	"net/http"

	"github.com/Basu008/EasySplit.git/server/config"
	postgresStorage "github.com/Basu008/EasySplit.git/server/storage/postgres"
	"github.com/gorilla/mux"
)

type Server struct {
	httpServer *http.Server
	Router     *mux.Router
	Config     *config.Config
	Postgres   postgresStorage.PostgresStorage
}
