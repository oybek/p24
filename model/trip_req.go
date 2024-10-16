package model

import (
	"fmt"
	"time"
)

type TripReq struct {
	ChatId    int64     `json:"chat_id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	StartDate time.Time `json:"start_date"`
}

func (t TripReq) IsValid() bool {
	return len(t.From) > 0 && len(t.To) > 0
}

func (t *TripReq) String() string {
	dateText := fmt.Sprintf("%02d/%02d/%04d", t.StartDate.Day(), t.StartDate.Month(), t.StartDate.Year())
	return fmt.Sprintf("%s ➡️ %s 🕖 %s", t.From, t.To, dateText)
}
