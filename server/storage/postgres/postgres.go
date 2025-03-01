package postgresStorage

import (
	"fmt"
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
	fmt.Print("Connected to Postgres\n")
	return &PostgresStorage{Config: c, DB: db}
}

func (p *PostgresStorage) Close() {
	if p.DB != nil {
		db, err := p.DB.DB()
		if err != nil {
			log.Fatalf("failed to retrieve db")
			return
		}
		db.Close()
		fmt.Println("Database connection closed.")
	}
}
