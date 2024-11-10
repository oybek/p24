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
) (*model.User, error) {
	user := model.User{}
	err := tx.QueryRow(`
		SELECT "chat_id", "uuid" FROM users
		WHERE "uuid" = $1
	`, UUID).Scan(&user.ChatId, &user.UUID)
	return &user, err
}
