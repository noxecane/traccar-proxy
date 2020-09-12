package model

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
