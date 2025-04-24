package app

import (
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

func cleanUpDB(db *gorm.DB, m any) error {
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(m).Error; err != nil {
		return err
	}
	return nil
}
