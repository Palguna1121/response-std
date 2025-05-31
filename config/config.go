package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	APP_PORT    string `json:"port"`
	DB_HOST     string `json:"db_host"`
	DB_PORT     string `json:"db_port"`
	DB_USER     string `json:"db_user"`
	DB_PASSWORD string `json:"db_password"`
	DB_NAME     string `json:"db_name"`
}

var ENV *Config

func InitConfig() {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err := viper.Unmarshal(&ENV); err != nil {
		panic(fmt.Errorf("unable to decode into struct: %w", err))
	}
}
