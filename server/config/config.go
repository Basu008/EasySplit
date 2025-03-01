package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	ServerConfig    ServerConfig    `mapstructure:"server"`
	DatabaseConfig  DatabaseConfig  `mapstructure:"database"`
	TokenAuthConfig TokenAuthConfig `mapstructure:"token"`
	APIConfig       APIConfig       `mapstructure:"api"`
	APPConfig       APPConfig       `mapstructure:"app"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"sslMode"`
	Name     string `mapstructure:"name"`
}

func (d *DatabaseConfig) ConnectionURL() string {
	return fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=%s", d.Host, d.Username, d.Name, d.Password, d.SSLMode)
}

type TokenAuthConfig struct {
	JWTSignKey   string `mapstructure:"jwtSignKey"`
	JWTExpiresAt int64  `mapstructure:"expiresAt"`
}

type APIConfig struct {
}

type APPConfig struct {
}

func GetConfig() *Config {
	fileName := "default"
	viper.SetConfigName(fileName)
	viper.SetConfigType("toml")
	viper.AddConfigPath("../conf/")
	viper.AddConfigPath("../../conf/")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./conf/")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("couldn't load config: %s", err)
		os.Exit(1)
	}
	config := &Config{}
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Printf("couldn't read config: %s", err)
		os.Exit(1)
	}
	return config
}
