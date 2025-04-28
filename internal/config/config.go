// Package config provides configuration management using AWS Systems Manager Parameter Store
package config

import (
	"context"
	"crypto/sha1"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/google/uuid"
)

var (
	instance *Config
	once     sync.Once
)

type Config struct {
	DbDsn         string
	Port          string
	GLOBAL_NS     uuid.UUID
	Salt          string
	Cors          map[string]bool
	GeoApi        string
	AdminPassword string
}

type ConfigDefinition struct {
	Path         string
	Type         string
	DefaultValue string
}

var configDefinitions = map[string]ConfigDefinition{
	"DB_DSN": {
		Path: "/backend/internal/db_dsn",
		Type: "SecureString",
	},
	"PORT": {
		Path:         "/backend/ports/analytics",
		Type:         "String",
		DefaultValue: ":8300",
	},
	"CORS_ORIGIN": {
		Path: "/backend/internal/analytics-cors-origin",
		Type: "StringList",
	},
	"GLOBAL_UUID_NAMESPACE_SEED": {
		Path: "/backend/internal/global-uuid-namespace-seed",
		Type: "SecureString",
	},
	"HASH_SALT": {
		Path: "/backend/internal/hash-salt",
		Type: "SecureString",
	},
	"GEO_API": {
		Path: "/backend/internal/geo-api",
		Type: "String",
	},
	"ADMIN_PASSWORD": {
		Path: "/backend/internal/admin-password",
		Type: "SecureString",
	},
}

func stringToNamespaceUUID(s string) uuid.UUID {
	hash := sha1.Sum([]byte(s)) // create a SHA-1 hash of the string

	var namespace uuid.UUID
	copy(namespace[:], hash[:16]) // use the first 16 bytes of the hash to create a UUID
	return namespace
}

func getSystemsManagerParameter(paramName string, ssmClient *ssm.Client) string {
	paramInfo, exists := configDefinitions[paramName]
	if !exists {
		log.Fatalf("***ERROR (config): Parameter '%s' not found in configDefinitions", paramName)
	}
	isEncrypted := paramInfo.Type == "SecureString"

	log.Printf("Attempting to retrieve parameter: %s (Path: %s)", paramName, paramInfo.Path)

	param, err := ssmClient.GetParameter(context.TODO(), &ssm.GetParameterInput{
		Name:           aws.String(paramInfo.Path),
		WithDecryption: aws.Bool(isEncrypted),
	})

	if err != nil {
		log.Printf("ERROR retrieving parameter %s: %v", paramName, err)

		if paramInfo.DefaultValue != "" {
			log.Printf("Using default value for %s", paramName)
			return paramInfo.DefaultValue
		}
		errorMsg := fmt.Sprintf("***ERROR (config): Failed to retrieve parameter '%s' from Systems Manager: %v", paramName, err)
		log.Fatal(errorMsg)
	}
	log.Printf("Successfully retrieved parameter: %s", paramName)

	return *param.Parameter.Value
}

func LoadConfig() *Config {
	once.Do(func() {
		log.Println("Loading configuration from AWS Systems Manager Parameter Store...")

		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion("us-east-1"), // Specify your AWS region
		)
		if err != nil {
			log.Fatal("***ERROR (config): Unable to load AWS SDK config: ", err)
		}
		log.Println("AWS SDK Config loaded successfully")

		ssmClient := ssm.NewFromConfig(cfg)

		corsString := getSystemsManagerParameter("CORS_ORIGIN", ssmClient)
		if corsString == "" {
			log.Fatal("CORS_ORIGIN parameter is empty")
		}
		corsUrls := strings.Split(corsString, ",")

		corsOrigin := make(map[string]bool)
		for _, url := range corsUrls {
			trimmedURL := strings.TrimSpace(url)
			if trimmedURL != "" {
				corsOrigin[trimmedURL] = true
			}
		}

		dsn := getSystemsManagerParameter("DB_DSN", ssmClient)
		port := getSystemsManagerParameter("PORT", ssmClient)
		global_seed := getSystemsManagerParameter("GLOBAL_UUID_NAMESPACE_SEED", ssmClient)
		salt := getSystemsManagerParameter("HASH_SALT", ssmClient)
		geo_api := getSystemsManagerParameter("GEO_API", ssmClient)
		adminPassword := getSystemsManagerParameter("ADMIN_PASSWORD", ssmClient)

		if global_seed == "" || salt == "" {
			log.Fatal("GLOBAL_UUID_NAMESPACE_SEED or HASH_SALT parameter is empty")
		}
		global_ns := stringToNamespaceUUID(global_seed)

		instance = &Config{
			DbDsn:         dsn,
			Port:          port,
			GLOBAL_NS:     global_ns,
			Salt:          salt,
			Cors:          corsOrigin,
			GeoApi:        geo_api,
			AdminPassword: adminPassword,
		}

		log.Println("Configuration loaded successfully")
	})

	return instance
}
