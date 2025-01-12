package main

import (
	"auth-app/cmd"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := cmd.NewRootCommand().Execute(); err != nil {
		log.Fatal().Msgf("Failed runn app: %s", err.Error())
	}
}
