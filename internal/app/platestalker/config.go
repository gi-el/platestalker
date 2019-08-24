// Platestalker - scraper of openalpr agent debug endpoints
//
// Copyright (c) 2019 G. de Nijs

// Application configuration

package platestalker

import (
	"time"
)

type Config struct {
	AgentAddress string
	Interval time.Duration
	DBFile string
}
