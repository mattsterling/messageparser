package config

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
)

// GlobalConfig is a struct that is initialized at start up
// with all the configuration supplied to the env variables.
var GlobalConfig = GetConfig()

// Config is the structure defining the configuration
// of this program.
type Config struct {
	Port    int
	MsgSize int
}

// GetConfig will parse the environment variables and return
// Config struct with the configuration.
func GetConfig() Config {
	var c Config
	err := envconfig.Process("message_parser", &c)
	if err != nil {
		log.Error("Failed to parse configuration:", err)
		os.Exit(1)
	}
	return c
}
