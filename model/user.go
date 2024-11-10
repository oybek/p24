package model

import (
	"github.com/google/uuid"
)

type User struct {
	ChatId int64
	UUID   uuid.UUID
}
