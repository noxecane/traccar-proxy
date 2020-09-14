package proxy

import "time"

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
