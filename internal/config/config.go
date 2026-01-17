package config

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

// Config holds application configuration.
type Config struct {
	BotToken string
	DataFile string
	Logger   *log.Logger
}

// New creates a new Config by reading environment variables.
func New() (*Config, error) {
	token := os.Getenv("API_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("API_TOKEN environment variable is required")
	}

	dataFile := os.Getenv("DATA_FILE")
	if dataFile == "" {
		dataFile = "data.json"
	}

	return &Config{
		BotToken: token,
		DataFile: dataFile,
		Logger:   log.New(os.Stderr),
	}, nil
}
