package main

import (
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog"
)

var (
	branch         string
	buildTimestamp string
	commit         string
	version        string
)

type config struct {
	Role string `env:"ROLE,required"`
	Port string `env:"PORT"`
}

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger.Info().
		Str("commit", commit).
		Str("branch", branch).
		Str("version", version).
		Str("buildTimestamp", buildTimestamp).
		Msg("starting up")

	config := config{}
	if err := env.Parse(&config); err != nil {
		logger.Fatal().Err(err).Msg("failed to parse config")
	}

	logger.Info().Any("config", config).Msg("read config")
}
