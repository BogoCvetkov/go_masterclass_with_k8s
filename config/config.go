package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DB_URL                 string        `mapstructure:"DB_URL"`
	DB_DRIVER              string        `mapstructure:"DB_DRIVER"`
	PORT                   string        `mapstructure:"PORT"`
	GRPC_PORT              string        `mapstructure:"GRPC_PORT"`
	GRPC_GATEWAY_PORT      string        `mapstructure:"GRPC_GATEWAY_PORT"`
	ENV                    string        `mapstructure:"ENV"`
	TOKEN_DURATION         time.Duration `mapstructure:"TOKEN_DURATION"`
	TOKEN_SECRET           string        `mapstructure:"TOKEN_SECRET"`
	REFRESH_TOKEN_DURATION time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	SMTP_HOST              string        `mapstructure:"SMTP_HOST"`
	SMTP_PORT              string        `mapstructure:"SMTP_PORT"`
	SMTP_USER              string        `mapstructure:"SMTP_USER"`
	SMTP_PASS              string        `mapstructure:"SMTP_PASS"`
	SMTP_FROM              string        `mapstructure:"SMTP_FROM"`
	REDIS                  string        `mapstructure:"REDIS"`
}

func LoadConfig() *Config {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	viper.SetDefault("ENV", "DEV")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("GRPC_PORT", "9000")
	viper.SetDefault("GRPC_GATEWAY_PORT", "5000")
	viper.SetDefault("SMTP_FROM", "bogo@test.com")

	// Allow override by ENV vars
	viper.SetDefault("DB_URL", "")
	viper.SetDefault("DB_DRIVER", "")
	viper.SetDefault("TOKEN_DURATION", "")
	viper.SetDefault("TOKEN_SECRET", "")
	viper.SetDefault("REFRESH_TOKEN_DURATION", "")
	viper.SetDefault("SMTP_HOST", "")
	viper.SetDefault("SMTP_PORT", "")
	viper.SetDefault("SMTP_USER", "")
	viper.SetDefault("SMTP_PASS", "")
	viper.SetDefault("SMTP_FROM", "")
	viper.SetDefault("REDIS", "")

	duration, err := time.ParseDuration("30m")
	if err != nil {
		panic(fmt.Errorf("error parsing jwt duration: %w", err))
	}
	viper.SetDefault("JWTDuration", duration)

	// Override config with env variables outside of it
	viper.AutomaticEnv()

	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("error unmarschaling config file: %w", err))
	}

	return &config
}
