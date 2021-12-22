package config

import (
	"os"
)

// HTTP Server config type
type Config struct {
	HIBPKey    string
	HIBPHost   string
	ServerPort string
}

// Create the Config struct from existing environment variable keys
func FromEnv() Config {
	hibpKey := os.Getenv("HIBP_KEY")
	hibpHost := os.Getenv("HIBP_HOST")
	serverPort := os.Getenv("SERVER_PORT")

	return Config{
		HIBPKey:    hibpKey,
		HIBPHost:   hibpHost,
		ServerPort: serverPort,
	}
}
