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
		`INSERT INTO users ("chat_id", "uuid")
		 VALUES ($1, $2)
		 ON CONFLICT ("chat_id") DO
		 UPDATE SET "uuid" = $3`,
		user.ChatId, user.UUID, user.UUID,
	)
}

func SelectUser(
	tx TransactionOps,
	UUID uuid.UUID,
) (users []model.User, err error) {
	rows, err := tx.Query(`SELECT "chat_id", "uuid" FROM users WHERE "uuid" = $1`, UUID)
	if err != nil {
		return
	}

	for rows.Next() {
		user := model.User{}
		err = rows.Scan(&user.ChatId, &user.UUID)
		if err != nil {
			return
		}

		users = append(users, user)
	}

	return
}
