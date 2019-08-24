// Platestalker - scraper of openalpr agent debug endpoints
//
// Copyright (c) 2019 G. de Nijs

// Service initialization and control

package platestalker

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/gi-el/platestalker/pkg/alpr"
	"github.com/gi-el/platestalker/pkg/storage"

	"github.com/rs/zerolog/log"
)

type Service struct {
	config *Config
}

func NewService(config *Config) *Service {
	s := &Service{
		config: config,
	}

	return s
}

func (s *Service) Run() error {
	scraper, err := alpr.NewScraper(s.config.AgentAddress, s.config.Interval)
	if err != nil {
		return fmt.Errorf("error creating scraper: %s", err)
	}

	db, err := storage.NewStorage(s.config.DBFile)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := signalContext(context.Background())
	defer cancel()

	go scraper.Run(ctx)

	for {
		log.Debug().Msg("waiting for plate")
		plate, err := scraper.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Warn().Err(err).Msg("error; continuing")
			continue
		}

		log.Info().Time("timestamp", plate.Start()).Str("plate", plate.PlateNumber).Float64("confidence", plate.PlateNumberConfidence).Msg("")

		// Store plate
		err = db.StorePlate(plate)
		if err != nil {
			log.Warn().Err(err).Msg("error storing plate; continuing")
		}

		// Handle any actions
		// TBD
	}

	return nil
}

func signalContext(parent context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-sigChan
		log.Info().Str("signal", s.String()).Msg("received signal; gracefully terminating")
		cancel()
	}()

	return ctx, cancel
}
