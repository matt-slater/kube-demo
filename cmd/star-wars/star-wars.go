package main

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog"

	"github.com/matt-slater/star-wars/internal/pkg/deathstar"
	"github.com/matt-slater/star-wars/internal/pkg/rebelbase"
	"github.com/matt-slater/star-wars/internal/pkg/xwing"
)

var (
	branch         string
	buildTimestamp string
	commit         string
	version        string
)

type config struct {
	Role      string `env:"ROLE,required"`
	Port      string `env:"PORT,required"`
	Commander string `env:"COMMANDER,required"`
	Target    string `env:"TARGET,required"`
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

	switch config.Role {
	case "deathstar":
		ds := deathstar.New(config.Port, logger.With().Str("component", "deathstar").Logger())

		err := ds.Launch()
		if err != nil {
			logger.Fatal().Err(err).Msg("deathstar failure to launch")
		}
	case "rebelbase":
		rb := rebelbase.New(config.Commander, config.Target, config.Port, logger.With().Str("component", "rebelbase").Logger())

		err := rb.Launch()
		if err != nil {
			logger.Fatal().Err(err).Msg("rebelbase failure to launch")
		}
	case "xwing":
		xw := xwing.New(config.Target, config.Commander, logger.With().Str("component", "xwing").Logger())

		err := xw.Launch()
		if err != nil {
			logger.Fatal().Err(err).Msg("xwing mission aborted")
		}
	default:
		logger.Fatal().Msg(fmt.Sprintf("unrecognized role: %s", config.Role))
	}
}
