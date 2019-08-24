// Platestalker - scraper of openalpr debug endpoints
//
// Copyright (c) 2019 G. de Nijs

// Run command; runs main scraper and storage loop

package commands

import (
	"time"

	"github.com/gi-el/platestalker/internal/app/platestalker"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var config platestalker.Config

var runCmd = &cobra.Command{
	Use: "run",
	Short: "Run the main scraper and storage functionality",
	Run: func(cmd *cobra.Command, args []string) {
		onRun()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVar(&config.AgentAddress, "address", "", "URI to ALPR agent debug webserver")
	runCmd.Flags().StringVar(&config.DBFile, "db", "", "Path to the SQLite database file")
	runCmd.Flags().DurationVar(&config.Interval, "interval", time.Second * 30, "Interval between checking for new plates")

	runCmd.MarkFlagRequired("address")
	runCmd.MarkFlagRequired("db")
}

func onRun() {
	service := platestalker.NewService(&config)

	err := service.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
}
