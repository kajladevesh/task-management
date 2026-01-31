package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DB_HOST     string `mapstructure:"DB_HOST"`
	DB_PORT     string `mapstructure:"DB_PORT"`
	DB_USER     string `mapstructure:"DB_USER"`
	DB_PASSWORD string `mapstructure:"DB_PASSWORD"`
	DB_NAME     string `mapstructure:"DB_NAME"`
	DB_SSLMODE  string `mapstructure:"DB_SSLMODE"`
	APP_ENV     string `mapstructure:"DB_ENV"`
	APP_PORT    string `mapstructure:"DB_PORT"`
	JWT_SECRET  string `mapstructure:"JWT_SECRET"`
}

func LoadConfig() (*Config, error) {

	config := &Config{}

	envConfigFileName := ".env.local"

	viper.AutomaticEnv()
	viper.AddConfigPath("../.secrets")
	viper.SetConfigName(envConfigFileName)
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found. ")
		} else {
			return nil, fmt.Errorf("Failed to read config file :- %v", err)
		}
	}

	err = viper.Unmarshal(&config)

	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal config :- %w", err)
	}

	return config, nil
}
