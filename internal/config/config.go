package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

// Config структура конфига для сервера.
type Config struct {
	Env         string        `yaml:"env" env:"ENV" env-default:"local"`
	DBPath      string        `yaml:"db_path" env:"DB" env-default:"user=postgres password=postgres dbname=postgres sslmode=disable"`
	StoragePath string        `yaml:"storage_path" env:"STORE" env-default:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
	SyncTimeout time.Duration `yaml:"sync_timeout" env:"SYNCTIME" env-default:"300s"`
	Listener    struct {
		Addr        string        `yaml:"addr" env:"ADDR" env-default:"localhost:8080"`
		Timeout     time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"4s"`
		IdleTimeout time.Duration `yaml:"idle_timeout" env:"ITIMEOUT" env-default:"60s"`
	} `yaml:"listener"`
}

func NewConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config.yaml"
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
