package model

import "time"

type RideRequest struct {
	From string    `json:"from"`
	To   string    `json:"to"`
	When time.Time `json:"when"`
}
