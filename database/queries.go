package database

import (
	"database/sql"
	"strings"

	"github.com/oybek/choguuket/model"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func InsertTrip(
	tx TransactionOps,
	trip *model.Trip,
) (sql.Result, error) {
	path := strings.Join(trip.Path, ",")
	return tx.Exec(
		`INSERT INTO trips (chat_id, "path", phone, comment, start_at)
		 VALUES ($1, LOWER($2), $3, $4, $5)`,
		trip.ChatId, path, trip.Phone, trip.Comment, trip.StartTime,
	)
}

func InsertTripReq(
	tx TransactionOps,
	tripReq *model.TripReq,
) (int64, error) {
	var id int64
	err := tx.QueryRow(
		`INSERT INTO trip_reqs (chat_id, "from", "to", start_date)
		 VALUES ($1, LOWER($2), LOWER($3), $4)
		 RETURNING id`,
		tripReq.ChatId, tripReq.From, tripReq.To, tripReq.StartDate,
	).Scan(&id)
	return id, err
}

func SearchTrip(
	tx TransactionOps,
	tripReq *model.TripReq,
) ([]model.Trip, error) {
	rows, err := tx.Query(
		`SELECT chat_id, "path", phone, comment, start_at FROM trips
		 WHERE start_at > NOW() AND start_at::date = $1::date AND path LIKE '%'||LOWER($2)||'%'||LOWER($3)||'%'`,
		tripReq.StartDate, tripReq.From, tripReq.To,
	)
	if err != nil {
		return nil, err
	}

	return ItToSlice(rows, func(it Iterator) (t model.Trip, e error) {
		var path string
		e = it.Scan(
			&t.ChatId,
			&path,
			&t.Phone,
			&t.Comment,
			&t.StartTime,
		)
		t.Path = strings.Split(path, ",")
		for i := range t.Path {
			t.Path[i] = cases.Title(language.Und).String(t.Path[i])
		}
		return t, e
	})
}

func SearchTripReq(
	tx TransactionOps,
	trip model.Trip,
) ([]model.TripReq, error) {
	path := strings.Join(trip.Path, ",")
	rows, err := tx.Query(
		`SELECT chat_id, "from", "to", "start_date" FROM trip_reqs
		 WHERE "start_date" > NOW() AND "start_date"::date = $1::date AND LOWER($2) LIKE '%'||"from"||'%'||"to"||'%'`,
		trip.StartTime, path,
	)
	if err != nil {
		return nil, err
	}

	return ItToSlice(rows, func(it Iterator) (t model.TripReq, e error) {
		e = it.Scan(
			&t.ChatId,
			&t.From,
			&t.To,
			&t.StartDate,
		)
		t.From = cases.Title(language.Und).String(t.From)
		t.To = cases.Title(language.Und).String(t.To)
		return t, e
	})
}

// Iterator - description of the iterator interface
type Iterator interface {
	Next() bool
	Scan(dest ...any) error
}

// ItToSlice convert Iterator to Slice
func ItToSlice[T any](it Iterator, scan func(Iterator) (T, error)) ([]T, error) {
	ts := []T{}

	for it.Next() {
		t, err := scan(it)
		if err != nil {
			return ts, err
		}
		ts = append(ts, t)
	}

	return ts, nil
}
