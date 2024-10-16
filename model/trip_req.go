package model

import "time"

type TripReq struct {
	ChatId    int64     `json:"chat_id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	StartDate time.Time `json:"start_date"`
}
