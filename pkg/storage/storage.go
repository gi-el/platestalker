package storage

import (
	"github.com/gi-el/platestalker/pkg/alpr"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type Storage struct {
	db *sqlx.DB
}

func NewStorage(dbfile string) (*Storage, error) {
	log.Info().Str("file", dbfile).Msg("opening database")
	db, err := dbConnect(dbfile)
	if err != nil {
		return nil, err
	}

	err = createSchema(db)
	if err != nil {
		return nil, err
	}

	s := &Storage{
		db: db,
	}

	return s, nil
}

func (s *Storage) StorePlate(plate *alpr.ALPRGroup) error {
	return insertPlate(s.db, plate)
}

func (s *Storage) Close() error {
	err := s.db.Close()
	if err != nil {
		return err
	}
	log.Info().Msg("database closed")
	return nil
}
