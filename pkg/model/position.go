package model

import "time"

type Position struct {
	tableName  struct{} `pg:"tc_positions"`
	ID         uint
	CreatedAt  time.Time `pg:"servertime"`
	RecordedAt time.Time `pg:"devicetime"`
	Valid      bool      `pg:",use_zero"`
	Device     uint      `pg:"deviceid"`
	Latitude   float64   `pg:",use_zero"`
	Longitude  float64   `pg:",use_zero"`
	Altitude   float64   `pg:",use_zero"`
	Speed      float64   `pg:",use_zero"`
	Course     float64   `pg:",use_zero"`
	Payload    string    `pg:"attributes"`
	Accuracy   string
	Address    string
	Protocol   string
	Network    string
	FixedAt    time.Time `pg:"fixtime"`
}

type Attributes struct {
	FuelConsumption     float32 `json:"fuelConsumption,omitempty"`
	Raw                 string  `json:"raw,omitempty"`
	GSensor             string  `json:"gSensor,omitempty"`
	Result              string  `json:"result,omitempty"`
	Status              uint    `json:"status,omitempty"`
	Motion              bool    `json:"motion,omitempty"`
	ClearedDistance     float32 `json:"clearedDistance,omitempty"`
	TotalDistance       float32 `json:"totalDistance,omitempty"`
	RPM                 uint    `json:"rpm,omitempty"`
	Alarm               string  `json:"alarm,omitempty"`
	Ignition            bool    `json:"ignition,omitempty"`
	DTC                 string  `json:"dtcs,omitempty"`
	OBDSpeed            uint    `json:"obdSpeed,omitempty"`
	EngineLoad          int     `json:"engineLoad,omitempty"`
	CoolantTemperature  int     `json:"coolantTemp,omitempty"`
	Distance            float32 `json:"distance,omitempty"`
	TripOdometer        uint    `json:"tripOdometer,omitempty"`
	IntakeTemperature   int     `json:"intakeTemp,omitempty"`
	Odometer            uint64  `json:"odometer,omitempty"`
	MapIntake           int     `json:"mapIntake,omitempty"`
	Throttle            float32 `json:"throttle,omitempty"`
	MilDistance         float32 `json:"milDistance,omitempty"`
	Satellites          uint    `json:"sat,omitempty"`
	TripFuelConsumption float32 `json:"tripFuelConsumption,omitempty"`
}
