package telegram

import (
	"fmt"
	"time"

	"github.com/oybek/p24/model"
)

const nl = "\n"

var cityNames map[string]string

func (bot *Bot) InitCityNames() (err error) {
	cityNames, err = bot.mc.CityNamesGet()
	return err
}

func CityName(key string) string {
	value, exists := cityNames[key]
	if exists {
		return value
	}
	return key
}

func Show(t *model.Trip) string {
	localTime := t.StartDate.UTC().Add(
		time.Duration(t.Meta.TimeOffset) * time.Hour,
	).Format("02/01/2006 15:04")
	return fmt.Sprintf(
		"%s - %s"+nl+
			"Время: %s"+nl+
			"Кол-во мест: %d",
		CityName(t.CityA), CityName(t.CityB), localTime, t.PassengerCount,
	)
}
