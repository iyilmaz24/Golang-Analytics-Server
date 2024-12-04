package config

import (
	"crypto/sha1"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
)

var (
	instance *Config
	once     sync.Once
)

type Config struct {
	DSN  string
	Port string
	FL_NS uuid.UUID
	NM_NS uuid.UUID
	DEFAULT_NS uuid.UUID
	Salt  string
}

func stringToNamespaceUUID(s string) uuid.UUID {
	hash := sha1.Sum([]byte(s)) // create a SHA-1 hash of the string
	
	var namespace uuid.UUID 
	copy(namespace[:], hash[:16]) // use the first 16 bytes of the hash to create a UUID
	return namespace
}

func LoadConfig() *Config {
	once.Do(func() { // ensure that the config is only loaded once

		dsn, ok := os.LookupEnv("DB_DSN")
		if !ok {
			log.Fatal("DB_DSN is not set in environment variables")
		}

		port, ok := os.LookupEnv("PORT")
		if !ok {
			port = ":8100"
		}

		fl_seed, ok := os.LookupEnv("FL_UUID_NAMESPACE_SEED")
		if !ok {
			log.Fatal("FL_UUID_NAMESPACE_SEED is not set in environment variables")
		}

		nm_seed, ok := os.LookupEnv("NM_UUID_NAMESPACE_SEED")
		if !ok {
			log.Fatal("NM_UUID_NAMESPACE_SEED is not set in environment variables")
		}

		default_seed, ok := os.LookupEnv("DEFAULT_UUID_NAMESPACE_SEED")
		if !ok {
			log.Fatal("DEFAULT_UUID_NAMESPACE_SEED is not set in environment variables")
		}

		salt, ok := os.LookupEnv("SALT")
		if !ok {
			log.Fatal("SALT is not set in environment variables")
		}

		fl_ns := stringToNamespaceUUID(fl_seed)
		nm_ns := stringToNamespaceUUID(nm_seed)
		default_ns := stringToNamespaceUUID(default_seed)

		instance = &Config{
			DSN:  dsn,
			Port: port,
			FL_NS: fl_ns,
			NM_NS: nm_ns,
			DEFAULT_NS: default_ns,
			Salt: salt,
		}
	})

	return instance
}

