package model

import "time"

type Event struct {
	tableName struct{} `pg:"tc_events"`
	ID        uint
	CreatedAt time.Time `pg:"servertime"`

	Type     string
	Device   uint   `pg:"deviceid"`
	Position uint   `pg:"position"`
	Payload  string `pg:"attributes"`

	GeoFence    uint `pg:"geofenceid"`
	Maintenance uint `pg:"maintenanceid"`
}
