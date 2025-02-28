package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	ServerConfig ServerConfig `mapstructure:"server"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

func GetConfigFromFile() {
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
}
