package app

import (
	"fmt"

	"github.com/Basu008/EasySplit.git/server/config"
	postgresStorage "github.com/Basu008/EasySplit.git/server/storage/postgres"
	"gorm.io/gorm"
)

func NewTestApp(c *config.Config) *App {
	p := postgresStorage.NewPostgresStorage(&c.DatabaseConfig)
	a := &App{
		Postgres: p,
		Config:   &c.APPConfig,
	}
	return a
}

func getTestConfig() *config.Config {
	return config.GetConfigFromFile("test")
}

func cleanUpDB(db *gorm.DB) error {
	if err := db.Exec("SET session_replication_role = 'replica';").Error; err != nil {
		return err
	}

	// List all your tables here
	tables := []string{
		"friends",
		"group_members",
		"groups",
		"expenses",
		"users",
		// "payments",
	}

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", table)).Error; err != nil {
			return err
		}
	}

	// Turn referential integrity back on
	if err := db.Exec("SET session_replication_role = 'origin';").Error; err != nil {
		return err
	}

	return nil
}
