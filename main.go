package main

import (
	"github.com/Basu008/EasySplit.git/server/config"
	postgresStorage "github.com/Basu008/EasySplit.git/server/storage/postgres"
)

func main() {
	c := config.GetConfigFromFile()
	postgresStorage.NewPostgresStorage(&c.DatabaseConfig)
}
