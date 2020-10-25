package traccar

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"tsaron.com/traccar-proxy/pkg/model"
)

func TransformPosition(p model.TraccarPosition) (model.Position, error) {
	pos := model.Position{
		ID:         p.ID,
		CreatedAt:  time.Time(p.CreatedAt),
		RecordedAt: time.Time(p.RecordedAt),
		Valid:      p.Valid,
		Device:     p.Device,
		Latitude:   p.Latitude,
		Longitude:  p.Longitude,
		Altitude:   p.Altitude,
		Speed:      p.Speed,
		Course:     p.Course,
	}

	var attr model.TraccarAttributes
	if err := json.Unmarshal([]byte(p.Payload), &attr); err != nil {
		return pos, errors.Wrap(err, "could not decode attributes")
	}

	pos.Meta = model.Attributes{
		FuelConsumption:     attr.FuelConsumption,
		Raw:                 attr.Raw,
		GSensor:             attr.GSensor,
		Motion:              attr.Motion,
		TotalDistance:       attr.TotalDistance,
		RPM:                 attr.RPM,
		Alarm:               attr.Alarm,
		Ignition:            attr.Ignition,
		DTC:                 attr.DTC,
		EngineLoad:          attr.EngineLoad,
		CoolantTemperature:  attr.CoolantTemperature,
		TripOdometer:        attr.TripOdometer,
		IntakeTemperature:   attr.IntakeTemperature,
		Odometer:            attr.Odometer,
		MapIntake:           attr.MapIntake,
		Throttle:            attr.Throttle,
		MilDistance:         attr.MilDistance,
		Satellites:          attr.Satellites,
		TripFuelConsumption: attr.TripFuelConsumption,
	}

	return pos, nil
}
