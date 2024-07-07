package main

import (
	"os"

	"github.com/rs/zerolog"
)

type envVars struct {
}

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	logger.Info().Msg("hello world")

}
