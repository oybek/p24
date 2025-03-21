package telegram

import (
	"fmt"
	"strings"
	"time"

	"github.com/oybek/p24/model"
)

type TripView struct {
	CityA     string
	CityB     string
	UserName  string
	Date      string
	Time      string
	SeatCount string
}

func FormatTrip(trip TripView, userType string) string {
	lineWidth := 26
	separator := strings.Repeat("-", lineWidth)

	// Helper function to format rows with left/right justification
	formatRow := func(label, value string) string {
		return fmt.Sprintf("%-8s%18s", label, value)
	}

	userTypeText := "Пассажир"
	if userType == "driver" {
		userTypeText = "Водитель"
	}
	return fmt.Sprintf(
		"%s - %s\n%s\n%s\n%s\n%s\n%s",
		trip.CityA, trip.CityB,
		separator,
		formatRow(userTypeText, trip.UserName),
		formatRow("Дата", trip.Date),
		formatRow("Время", trip.Time),
		formatRow("Мест", trip.SeatCount),
	)
}

func (bot *Bot) MapToTripView(trip *model.Trip, user *model.User) TripView {
	utcPlus6 := time.FixedZone("UTC+6", 6*60*60)
	localTime := trip.StartDate.In(utcPlus6)
	date := fmt.Sprintf("%d %s", localTime.Day(), monthsRU[localTime.Month()])

	name := user.Name
	if len(trip.UserName) > 0 {
		name = trip.UserName
	}
	return TripView{
		CityA:     bot.CityName(trip.CityA),
		CityB:     bot.CityName(trip.CityB),
		UserName:  name,
		Date:      date,                      // Example: "16 March 2025"
		Time:      localTime.Format("15:04"), // Example: "09:30"
		SeatCount: fmt.Sprintf("%d", trip.SeatCount),
	}
}

var monthsRU = map[time.Month]string{
	time.January:   "января",
	time.February:  "февраля",
	time.March:     "марта",
	time.April:     "апреля",
	time.May:       "мая",
	time.June:      "июня",
	time.July:      "июля",
	time.August:    "августа",
	time.September: "сентября",
	time.October:   "октября",
	time.November:  "ноября",
	time.December:  "декабря",
}
