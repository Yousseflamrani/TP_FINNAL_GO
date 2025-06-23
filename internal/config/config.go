package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port    int    `mapstructure:"port"`
		BaseURL string `mapstructure:"base_url"`
	} `mapstructure:"server"`

	Database struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
		SSLMode  string `mapstructure:"sslmode"`
	} `mapstructure:"database"`

	Analytics struct {
		BufferSize int `mapstructure:"buffer_size"`
	} `mapstructure:"analytics"`

	Monitor struct {
		IntervalMinutes int `mapstructure:"interval_minutes"`
	} `mapstructure:"monitor"`
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath("./configs")

	viper.SetConfigName("config")

	viper.SetConfigType("yaml")

	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.base_url", "http://localhost:8080")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.name", "urlshortener")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("analytics.buffer_size", 1000)
	viper.SetDefault("monitor.interval_minutes", 5)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Fichier de configuration non trouvé, utilisation des valeurs par défaut")
		} else {
			return nil, fmt.Errorf("erreur lors de la lecture du fichier de configuration: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("erreur lors du démapping de la configuration: %w", err)
	}

	log.Printf("Configuration loaded: Server Port=%d, DB Name=%s, Analytics Buffer=%d, Monitor Interval=%dmin",
		cfg.Server.Port, cfg.Database.Name, cfg.Analytics.BufferSize, cfg.Monitor.IntervalMinutes)

	return &cfg, nil
}
