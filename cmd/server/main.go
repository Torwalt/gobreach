package main

import (
	"fmt"
	"gobreach/cmd/server/config"
	"gobreach/internal/domains/breach"
	"gobreach/internal/ports/datasources/hibp"
	http_routes "gobreach/internal/ports/http"
	"log"
	"net/http"
	"time"
)

func main() {
	logger := log.Default()
	logger.Print("Starting Server")

	config := config.FromEnv()

	err := run(logger, &config)
	if err != nil {
		log.Fatalf("An error occured while running the server: %v", err)
	}

}

func run(l *log.Logger, c *config.Config) error {
	hconf := hibp.NewhibpConfig(c.HIBPHost, c.HIBPKey, 2)
	hclient := hibp.NewClient(http.DefaultClient, hconf, time.Sleep)
	bS := breach.NewService(nil, hclient)
	r := http_routes.NewRouter(bS, l)

	l.Printf("Server listening on port %v", c.HTTPPort)
	err := http.ListenAndServe(":" + c.HTTPPort, r.Router)
	if err != nil {
		return fmt.Errorf("could not start server: %v", err)
	}
	return nil
}
