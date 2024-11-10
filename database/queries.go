package database

import (
	"database/sql"

	"github.com/oybek/choguuket/model"
)

func UpsertUser(
	tx TransactionOps,
	user *model.User,
) (sql.Result, error) {
	return tx.Exec(
		`INSERT INTO users ("chat_id", "uuid")
		 VALUES ($1, $2)
		 ON CONFLICT ("chat_id") DO
		 UPDATE SET "uuid" = $3`,
		user.ChatId, user.UUID, user.UUID,
	)
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
