package model

import "time"

type Trip struct {
	ChatId    int64     `json:"chat_id"`
	Path      []string  `json:"path"`
	StartTime time.Time `json:"start_time"`
	Phone     string    `json:"phone"`
	Comment   string    `json:"comment"`
}
