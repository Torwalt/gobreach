package config

import (
	"fmt"
	"os"
	"strconv"
)

// HTTP Server config type
type Config struct {
	HIBPKey  string
	HIBPHost string
	HTTPPort string

	PostgresHost   string
	PostgresPort   int
	PostgresUser   string
	PostgresSecret string
	PostgresName   string
}

// Create the Config struct from existing environment variable keys
func FromEnv() Config {
	hibpKey := os.Getenv("HIBP_KEY")
	hibpHost := os.Getenv("HIBP_HOST")
	serverPort := os.Getenv("HTTP_PORT")

	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresSecret := os.Getenv("POSTGRES_SECRET")
	postgresName := os.Getenv("POSTGRES_NAME")

	postgresPortI, err := strconv.Atoi(postgresPort)
	if err != nil {
		panic(fmt.Sprintf("could not load config: %v", err))
	}

	return Config{
		HIBPKey:        hibpKey,
		HIBPHost:       hibpHost,
		HTTPPort:       serverPort,
		PostgresHost:   postgresHost,
		PostgresPort:   postgresPortI,
		PostgresUser:   postgresUser,
		PostgresSecret: postgresSecret,
		PostgresName:   postgresName,
	}
}
