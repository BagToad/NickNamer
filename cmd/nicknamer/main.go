package main

import (
	"os"

	"nicknamer/internal/bot"
	"nicknamer/internal/config"

	"github.com/charmbracelet/log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		logger := log.New(os.Stderr)
		logger.Fatal("Failed to load configuration:", err)
	}

	nickBot, err := bot.New(cfg)
	if err != nil {
		cfg.Logger.Fatal("Failed to create bot:", err)
	}

	if err := nickBot.Start(); err != nil {
		cfg.Logger.Fatal("Failed to start bot:", err)
	}
}
