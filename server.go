package main

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	conf "github.com/messageparser/config"
	handlers "github.com/messageparser/http"
)

func main() {

	// Set up logging
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	// Load configuration
	config := conf.GlobalConfig
	log.Debug("Configuration: ", config)

	if 0 == config.Port {
		log.Info("Port not provided, defaulting to 8080.")
		config.Port = 8080
	}

	// TODO: Create tests.
	// TODO: Update Readme.
	// TODO: pprof results.

	http.HandleFunc("/v1/message", handlers.ParseMessageHandler)
	http.HandleFunc("/health", handlers.HealthCheck)

	// Get to business
	log.Info("Starting Message Parsing server on port ", config.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)

}
