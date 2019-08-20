package main

import (
	"ehealth-migration/application"
	"ehealth-migration/lib/encoding"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Debug().Msg("Run scripts")
	app := application.NewApplication()
	err := encoding.Run(app)
	if err != nil {
		log.Error().Msgf("encoding.Run error: %s", err)
		return
	}
}
