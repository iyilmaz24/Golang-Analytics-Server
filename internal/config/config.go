package config

import (
	"log"
	"os"
)

type Config struct {
	DSN  string
	Port string
}

func LoadConfig() *Config {
	dsn, ok := os.LookupEnv("DB_DSN")
	if !ok {
		log.Fatal("DB_DSN is not set in environment variables")
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = ":8100"
	}

	return &Config{
		DSN:  dsn,
		Port: port,
	}
}