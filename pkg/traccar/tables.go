package traccar

import "time"

type Device struct {
	tableName struct{} `pg:"tc_devices"`
	ID        uint
	UpdatedAt time.Time `pg:"lastupdate"`

	Name         string
	ExternalID   string `pg:"uniqueid"`
	LastPosition uint   `pg:"positionid"`
	Payload      string `pg:"attributes"`

	Phone    string
	Model    string
	Contact  string
	Category string
	Disabled bool
	Group    uint `pg:"groupid"`
}

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
	Accuracy   uint
	Address    string
	Protocol   string
	Network    string
	FixedAt    time.Time `pg:"fixtime"`
}
