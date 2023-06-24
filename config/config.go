package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DB            string        `mapstructure:"DB_URL"`
	DBDriver      string        `mapstructure:"DB_DRIVER"`
	Port          string        `mapstructure:"PORT"`
	TokenDuration time.Duration `mapstructure:"TOKEN_DURATION"`
	TokenSecret   string        `mapstructure:"TOKEN_SECRET"`
}

func LoadConfig() *Config {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	viper.SetDefault("Port", "8080")

	duration, err := time.ParseDuration("30m")
	if err != nil {
		panic(fmt.Errorf("error parsing jwt duration: %w", err))
	}
	viper.SetDefault("JWTDuration", duration)

	// Override config with env variables outside of it
	viper.AutomaticEnv()

	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		panic(fmt.Errorf("error reading config file: %w", err))
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("error unmarschaling config file: %w", err))
	}

	return &config
}
