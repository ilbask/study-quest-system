package config

import (
	"os"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	DSN string // Data Source Name
}

func LoadConfig() (*Config, error) {
	// Set defaults
	viper.SetDefault("server.port", "8080")
	
	// Get DSN from environment or use default
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "root:root@tcp(127.0.0.1:3306)/study_quest?charset=utf8mb4&parseTime=True&loc=Local"
	}
	viper.SetDefault("database.dsn", dsn)

	// Try to read from environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("STUDY_QUEST")
	
	// Try to read from config file (optional)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.ReadInConfig() // Ignore error if file doesn't exist

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

