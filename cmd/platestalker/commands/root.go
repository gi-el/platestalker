// Platestalker - scraper of openalpr debug endpoints
//
// Copyright (c) 2019 G. de Nijs

// Root command

package commands

import (
	"io/ioutil"
	stdlog "log"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "platestalker",
	Short: "Platestalker is a scraper of openalpr debug endpoints",
	Long:
`Platestalker extracts recognized license plates directly
from the alpr agent endpoints and stores it locally in a
SQLite database. This allows for much more flexible queries
and long term storage.`,
}

var debug bool = false

func init() {
	cobra.OnInitialize(initApp)

	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Debug logging")
}

func initApp() {
	// Set log level
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func Execute() {
	// Create colorful logger
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out: os.Stderr,
	})

	// Have standard libraries not output logging
	// (looking at you net/http)
	stdlog.SetOutput(ioutil.Discard)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("")
	}
}
