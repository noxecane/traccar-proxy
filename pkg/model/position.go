package model

import (
	"encoding/json"
	"strings"
	"time"
)

type ISOWithoutTZ time.Time

// imeplement Marshaler und Unmarshalere interface
func (i *ISOWithoutTZ) UnmarshalJSON(b []byte) error {
	// remove quotes
	tStr := strings.Trim(string(b), "\"")

	t, err := time.Parse("2006-01-02T15:04:05.999", tStr)
	if err != nil {
		return err
	}

	// update time in place
	*i = ISOWithoutTZ(t)

	return nil
}

func (i ISOWithoutTZ) MarshalJSON() ([]byte, error) {
	return json.Marshal(i)
}

type TraccarPosition struct {
	ID         uint
	CreatedAt  ISOWithoutTZ `json:"servertime"`
	RecordedAt ISOWithoutTZ `json:"devicetime"`
	Valid      bool         `json:"valid"`
	Device     uint         `json:"deviceid"`
	Latitude   float64      `json:"latitude"`
	Longitude  float64      `json:"longitude"`
	Altitude   float64      `json:"altitude"`
	Speed      float64      `json:"speed"`
	Course     float64      `json:"course"`
	Payload    string       `json:"attributes"`
	Accuracy   uint         `json:"accuracy"`
	Address    string       `json:"address"`
	Protocol   string       `json:"protocol"`
	Network    string       `json:"network"`
	FixedAt    ISOWithoutTZ `json:"fixtime"`
}

type TraccarAttributes struct {
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

type Position struct {
	ID         uint       `json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	RecordedAt time.Time  `json:"recorded_at"`
	Valid      bool       `json:"valid"`
	Device     uint       `json:"device_id"`
	Latitude   float64    `json:"latitude"`
	Longitude  float64    `json:"longitude"`
	Altitude   float64    `json:"altitude"`
	Speed      float64    `json:"speed"`
	Course     float64    `json:"course"`
	Meta       Attributes `json:"metadata"`
}

type Attributes struct {
	FuelConsumption     float32 `json:"fuel_used,omitempty"`
	Raw                 string  `json:"raw_code,omitempty"`
	GSensor             string  `json:"accelerometer,omitempty"`
	Motion              bool    `json:"motion,omitempty"`
	TotalDistance       float32 `json:"total_distance,omitempty"`
	RPM                 uint    `json:"rpm,omitempty"`
	Alarm               string  `json:"alarm,omitempty"`
	Ignition            bool    `json:"ignition,omitempty"`
	DTC                 string  `json:"dtcs,omitempty"`
	EngineLoad          int     `json:"engine_load,omitempty"`
	CoolantTemperature  int     `json:"coolant_temparature,omitempty"`
	TripOdometer        uint    `json:"trip_odometer,omitempty"`
	IntakeTemperature   int     `json:"intake_temperature,omitempty"`
	Odometer            uint64  `json:"odometer,omitempty"`
	MapIntake           int     `json:"map_intake,omitempty"`
	Throttle            float32 `json:"throttle,omitempty"`
	MilDistance         float32 `json:"mil_distance,omitempty"`
	Satellites          uint    `json:"satellites,omitempty"`
	TripFuelConsumption float32 `json:"trip_fuel_used,omitempty"`
}
