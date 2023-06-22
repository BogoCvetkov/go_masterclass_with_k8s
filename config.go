package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DB       string `mapstructure:"DB_URL"`
	DBDriver string `mapstructure:"DB_DRIVER"`
	Port     string `mapstructure:"PORT"`
}

func LoadConfig() *Config {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	viper.SetDefault("Port", "8080")

	// Override config with env variables outside of it
	viper.AutomaticEnv()

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("error reading config file: %w", err))
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("error unmarschaling config file: %w", err))
	}

	return &config
}
