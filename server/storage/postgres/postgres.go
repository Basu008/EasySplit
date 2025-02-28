package postgresStorage

import (
	"log"

	"github.com/Basu008/EasySplit.git/server/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresStorage struct {
	Config *config.DatabaseConfig
	DB     *gorm.DB
}

func NewPostgresStorage(c *config.DatabaseConfig) *PostgresStorage {
	dsn := c.ConnectionURL()
	var db *gorm.DB
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to establish connection with postgres: %s", err.Error())
		return nil
	}
	return &PostgresStorage{Config: c, DB: db}
}
