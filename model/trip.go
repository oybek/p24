package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Trip struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	CityA          string             `bson:"city_a,omitempty" json:"city_a"`
	CityB          string             `bson:"city_b,omitempty" json:"city_b"`
	StartDate      time.Time          `bson:"start_time,omitempty" json:"start_time"`
	PassengerCount int                `bson:"passenger_count,omitempty" json:"passenger_count"`
	ChatID         int64              `bson:"chat_id,omitempty" json:"chat_id"`
	State          string             `bson:"state,omitempty" json:"state"`
	Meta           Meta               `bson:"-" json:"meta"`
}

func (t Trip) IsValid() bool {
	return t.CityA != "" && t.CityB != "" && t.PassengerCount > 0 && !t.StartDate.IsZero()
}
