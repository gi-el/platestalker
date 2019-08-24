package storage

import (
	"time"
)

var Tables = map[string]string{
	"plates": "id text PRIMARY KEY, start timestamp, end timestamp, camera_id int, plate string, confidence float, direction float, is_parked integer",
	"plate_images": "id text PRIMARY KEY, image blob",
	"vehicle_images": "id text PRIMARY KEY, image blob",
	"properties": "id text, key text, value text, confidence float",
}

type TablePlates struct {
	ID string `db:"id"`
	Start time.Time `db:"start"`
	End time.Time `db:"end"`
	CameraID int `db:"camera_id"`
	Plate string `db:"plate"`
	Confidence float64 `db:"confidence"`
	Direction float64 `db:"direction"`
	IsParked bool `db:"is_parked"`
}

type TablePlateImages struct {
	ID string `db:"id"`
	Image []byte `db:"image"`
}

type TableVehicleImages struct {
	ID string `db:"id"`
	Image []byte `db:"image"`
}

type TableProperties struct {
	ID string `db:"id"`
	Key string `db:"key"`
	Value string `db:"value"`
	Confidence float64 `db:"confidence"`
}
