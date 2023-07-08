package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DB              string        `mapstructure:"DB_URL"`
	DBDriver        string        `mapstructure:"DB_DRIVER"`
	Port            string        `mapstructure:"PORT"`
	GRPCPort        string        `mapstructure:"GRPC_PORT"`
	GRPCGatewayPort string        `mapstructure:"GRPC_GATEWAY_PORT"`
	Env             string        `mapstructure:"ENV"`
	TokenDuration   time.Duration `mapstructure:"TOKEN_DURATION"`
	TokenSecret     string        `mapstructure:"TOKEN_SECRET"`
	RTokenDuration  time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	SmtpHost        string        `mapstructure:"SMTP_HOST"`
	SmtpPort        string        `mapstructure:"SMTP_PORT"`
	SmtpUser        string        `mapstructure:"SMTP_USER"`
	SmtpPass        string        `mapstructure:"SMTP_PASS"`
	SmtpFrom        string        `mapstructure:"SMTP_FROM"`
	Redis           string        `mapstructure:"REDIS"`
}

func LoadConfig() *Config {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	viper.SetDefault("Env", "DEV")
	viper.SetDefault("Port", "8080")
	viper.SetDefault("GRPCPort", "9000")
	viper.SetDefault("GRPCGatewayPort", "5000")
	viper.SetDefault("SmtpFrom", "bogo@test.com")

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
