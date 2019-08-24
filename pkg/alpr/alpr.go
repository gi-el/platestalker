package alpr

import (
	"time"
)

type ALPRGroup struct {
	DataType string `json:"data_type"`
	Version int `json:"version"`
	EpochStart int64 `json:"epoch_start"`
	EpochEnd int64 `json:"epoch_end"`
	CameraID int `json:"camera_id"`
	Country string `json:"country,omitempty"`
	VehicleCropJPEG []byte `json:"vehicle_crop_jpeg,omitempty"`

	BestPlate Plate `json:"best_plate"`

	UUID string `json:"best_uuid"`
	PlateNumber string `json:"best_plate_number"`
	PlateNumberConfidence float64 `json:"best_confidence"`
	Region string `json:"best_region,omitempty"`
	RegionConfidence float64 `json:"best_region_confidence,omitempty"`
	TravelDirection float64 `json:"travel_direction,omitempty"`
	IsParked bool `json:"is_parked"`

	// Current properties: color, make, make_model, body_type, year, orientation
	Vehicle map[string][]Recognition `json:"vehicle,omitempty"`
}

func (a *ALPRGroup) Start() time.Time {
	timestamp := time.Unix(0, a.EpochStart * 1000000)
	return timestamp
}

func (a *ALPRGroup) End() time.Time {
	timestamp := time.Unix(0, a.EpochStart * 1000000)
	return timestamp
}

type Plate struct {
	PlateNumber string `json:"plate"`
	PlateNumberConfidence float64 `json:"confidence"`
	Region string `json:"region,omitempty"`
	RegionConfidence float64 `json:"region_confidence,omitempty"`
	PlateCropJPEG []byte `json:"plate_crop_jpeg,omitempty"`
}

type Recognition struct {
	Name string `json:"name"`
	Confidence float64 `json:"confidence"`
}
