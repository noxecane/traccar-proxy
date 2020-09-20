package model

type Device struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	ExternalID   string `json:"external_id"`
	LastPosition string `json:"last_position"`
}
