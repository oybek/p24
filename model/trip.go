package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Trip struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	StartTime time.Time          `bson:"start_time,omitempty" json:"start_time"`
	SeatCount int                `bson:"seat_count,omitempty" json:"seat_count"`
	CityA     string             `bson:"city_a,omitempty" json:"city_a"`
	CityB     string             `bson:"city_b,omitempty" json:"city_b"`
	Phone     string             `bson:"phone,omitempty" json:"phone"`
	UserName  string             `bson:"user_name" json:"user_name"`
	ChatID    int64              `bson:"chat_id,omitempty" json:"chat_id"`
	UserType  string             `bson:"user_type,omitempty" json:"user_type"`
	Meta      Meta               `bson:"-" json:"meta"`
	MessageId int64              `bson:"message_id,omitempty" json:"message_id"`
}

func (t Trip) IsValid() bool {
	return t.CityA != "" && t.CityB != "" && t.SeatCount > 0 && !t.StartTime.IsZero()
}

// StartDate implements sort.Interface for sorting by the StartDate field
type ByStartTime []Trip

func (a ByStartTime) Len() int           { return len(a) }
func (a ByStartTime) Less(i, j int) bool { return a[i].StartTime.Before(a[j].StartTime) }
func (a ByStartTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
