package database

import (
	"github.com/oybek/choguuket/model"
)

func AptekaInsert(
	tx TransactionOps,
	apteka *model.Apteka,
) (id int64, err error) {
	err = tx.QueryRow(
		`INSERT INTO apteka ("name", "phone", "address")
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		apteka.Name, apteka.Phone, apteka.Address,
	).Scan(&id)
	return id, err
}

func UserInsert(
	tx TransactionOps,
	user *model.User,
) error {
	_, err := tx.Exec(
		`INSERT INTO users ("chat_id", "apteka_id")
		 VALUES ($1, $2)`,
		user.ChatId, user.AptekaId,
	)
	return err
}
