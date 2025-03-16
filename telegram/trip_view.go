package telegram

import (
	"fmt"
	"strings"
	"time"

	"github.com/oybek/p24/model"
)

type TripView struct {
	CityA          string
	CityB          string
	PassengerName  string
	Date           string
	Time           string
	PassengerCount string
}

func FormatTrip(trip TripView) string {
	lineWidth := 26
	separator := strings.Repeat("-", lineWidth)

	// Helper function to format rows with left/right justification
	formatRow := func(label, value string) string {
		return fmt.Sprintf("%-8s%18s", label, value)
	}

	// Center align function
	centerText := func(text string) string {
		padding := (lineWidth - len(text)) / 2
		return fmt.Sprintf("%s%s", strings.Repeat(" ", padding), text)
	}

	return fmt.Sprintf(
		"%s - %s\n%s\n%s\n%s\n%s\n%s\n\n%s",
		trip.CityA, trip.CityB,
		separator,
		formatRow("Пассажир", trip.PassengerName),
		formatRow("Дата", trip.Date),
		formatRow("Время", trip.Time),
		formatRow("Мест", trip.PassengerCount),
		centerText("poputka24bot"),
	)
}

func MapToTripView(trip *model.Trip, user *model.User) TripView {
	utcPlus6 := time.FixedZone("UTC+6", 6*60*60)
	localTime := trip.StartDate.In(utcPlus6)
	date := fmt.Sprintf("%d %s", localTime.Day(), monthsRU[localTime.Month()])

	name := user.Name
	if len(trip.PassengerName) > 0 {
		name = trip.PassengerName
	}
	return TripView{
		CityA:          CityName(trip.CityA),
		CityB:          CityName(trip.CityB),
		PassengerName:  name,
		Date:           date,                      // Example: "16 March 2025"
		Time:           localTime.Format("15:04"), // Example: "09:30"
		PassengerCount: fmt.Sprintf("%d", trip.PassengerCount),
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
