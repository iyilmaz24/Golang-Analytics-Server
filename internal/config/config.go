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
	GLOBAL_NS uuid.UUID
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

		global_seed, ok := os.LookupEnv("GLOBAL_UUID_NAMESPACE_SEED")
		if !ok {
			log.Fatal("GLOBAL_UUID_NAMESPACE_SEED is not set in environment variables")
		}

		salt, ok := os.LookupEnv("HASH_SALT")
		if !ok {
			log.Fatal("SALT is not set in environment variables")
		}

		global_ns := stringToNamespaceUUID(global_seed)

		instance = &Config{
			DSN:  dsn,
			Port: port,
			GLOBAL_NS: global_ns,
			Salt: salt,
		}
	})

	return instance
}

