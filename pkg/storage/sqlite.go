package storage

import (
	"database/sql"
	"fmt"

	"github.com/gi-el/platestalker/pkg/alpr"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

const (
	TABLE_PREFIX = "alpr" // alpr_plates alpr_plate_images alpr_car_images alpr_properties

	CREATE_TABLE_TEMPLATE = "CREATE TABLE %s (%s)"
	SELECT_TABLE_QUERY = "SELECT sql FROM sqlite_master WHERE type='table' AND name=?"
)

func dbConnect(dbFile string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createSchema(db *sqlx.DB) error {
	for table, ddl := range(Tables) {
		tableName := fmt.Sprintf("%s_%s", TABLE_PREFIX, table)
		err := createTable(db, tableName, ddl)
		if err != nil {
			return err
		}
	}

	return nil
}

func createTable(db *sqlx.DB, table, ddl string) error {
	query := fmt.Sprintf(CREATE_TABLE_TEMPLATE, table, ddl)

	var schemaRes string
	err := db.Get(&schemaRes, SELECT_TABLE_QUERY, table)
	if err == sql.ErrNoRows {
		log.Info().Str("query", query).Msg("creating schema")
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
		return nil
	}
	if err != nil {
		return err
	}
	if schemaRes != query {
		return fmt.Errorf("non-matching table definition: %s != %s", schemaRes, query)
	}

	return nil
}

func insertPlate(db *sqlx.DB, group *alpr.ALPRGroup) error {
	log.Debug().Interface("plate", group).Msg("inserting plate")

	plate := &TablePlates{
		ID: group.UUID,
		Start: group.Start(),
		End: group.End(),
		CameraID: group.CameraID,
		Plate: group.PlateNumber,
		Confidence: group.PlateNumberConfidence,
		Direction: group.TravelDirection,
		IsParked: group.IsParked,
	}

	plateImage := &TablePlateImages{
		ID: group.UUID,
		Image: group.BestPlate.PlateCropJPEG,
	}

	vehicleImage := &TableVehicleImages{
		ID: group.UUID,
		Image: group.VehicleCropJPEG,
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.NamedExec("INSERT INTO alpr_plates VALUES (:id, :start, :end, :camera_id, :plate, :confidence, :direction, :is_parked)", plate)
	if err != nil {
		return err
	}

	_, err = tx.NamedExec("INSERT INTO alpr_plate_images VALUES (:id, :image)", plateImage)
	if err != nil {
		return err
	}

	_, err = tx.NamedExec("INSERT INTO alpr_vehicle_images VALUES (:id, :image)", vehicleImage)
	if err != nil {
		return err
	}

	for property, recognitions := range(group.Vehicle) {
		for _, recognition := range(recognitions) {
			property := &TableProperties{
				ID: group.UUID,
				Key: property,
				Value: recognition.Name,
				Confidence: recognition.Confidence,
			}
			_, err = tx.NamedExec("INSERT INTO alpr_properties VALUES (:id, :key, :value, :confidence)", property)
			if err != nil {
				return err
			}

		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
