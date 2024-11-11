package database

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/oybek/choguuket/model"
)

func UpsertUser(
	tx TransactionOps,
	user *model.User,
) (sql.Result, error) {
	return tx.Exec(
		`INSERT INTO users ("chat_id", "uuid", "nick")
		 VALUES ($1, $2, $3)
		 ON CONFLICT ("chat_id") DO
		 UPDATE SET "uuid" = $4, "nick" = $5`,
		user.ChatId, user.UUID, user.Nick, user.UUID, user.Nick,
	)
}

func SelectUser(
	tx TransactionOps,
	UUID uuid.UUID,
) (users []model.User, err error) {
	rows, err := tx.Query(`SELECT "chat_id", "uuid", "nick" FROM users WHERE "uuid" = $1`, UUID)
	if err != nil {
		return
	}

	for rows.Next() {
		user := model.User{}
		err = rows.Scan(&user.ChatId, &user.UUID, &user.Nick)
		if err != nil {
			return
		}

		users = append(users, user)
	}

	return
}
