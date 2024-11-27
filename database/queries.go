package database

import (
	"time"

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
		`INSERT INTO users ("chat_id", "apteka_id", "reader")
		 VALUES ($1, $2, $3)`,
		user.ChatId, user.AptekaId, user.Reader,
	)
	return err
}

func UserSelect(
	tx TransactionOps,
	chatId int64,
) (user model.User, err error) {
	err = tx.QueryRow(
		`SELECT "chat_id", "apteka_id", "reader"
		 FROM users WHERE "chat_id" = $1`,
		chatId,
	).Scan(&user.ChatId, &user.AptekaId, &user.Reader)
	return
}

func MedicineInsert(
	tx TransactionOps,
	medicine *model.Medicine,
) (id int64, err error) {
	err = tx.QueryRow(
		`INSERT INTO medicine ("name") VALUES ($1)
		 ON CONFLICT ("name")
		 DO UPDATE SET "name" = EXCLUDED."name"
		 RETURNING "id"`,
		medicine.Name,
	).Scan(&id)
	return
}

func AptekaMedicineInsert(
	tx TransactionOps,
	aptekaId int64,
	medicineId int64,
	amount int64,
	updated time.Time,
) error {
	_, err := tx.Exec(
		`INSERT INTO apteka_medicine ("apteka_id", "medicine_id", "amount", "updated")
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT ("apteka_id", "medicine_id")
		 DO UPDATE SET "amount" = $5, "updated" = $6`,
		aptekaId, medicineId, amount, updated, amount, updated,
	)
	return err
}
