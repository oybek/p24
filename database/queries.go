package database

import (
	"time"

	"github.com/lib/pq"
	"github.com/oybek/choguuket/model"
	"github.com/samber/lo"
)

func AptekaInsert(tx TransactionOps, apteka *model.Apteka) (id int64, err error) {
	err = tx.QueryRow(
		`INSERT INTO apteka ("name", "phone", "address")
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		apteka.Name, apteka.Phone, apteka.Address,
	).Scan(&id)
	return id, err
}

func UserInsert(tx TransactionOps, user *model.User) error {
	_, err := tx.Exec(
		`INSERT INTO users ("chat_id", "apteka_id", "reader")
		 VALUES ($1, $2, $3)`,
		user.ChatId, user.AptekaId, user.Reader,
	)
	return err
}

func UserSelect(tx TransactionOps, chatId int64) (user model.User, err error) {
	err = tx.QueryRow(
		`SELECT "chat_id", "apteka_id", "reader"
		 FROM users WHERE "chat_id" = $1`,
		chatId,
	).Scan(&user.ChatId, &user.AptekaId, &user.Reader)
	return
}

func MedicineInsert(tx TransactionOps, medicine *model.Medicine) (id int64, err error) {
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

func MedicineSearch(tx TransactionOps, names []string) ([]int64, error) {
	rows, err := tx.Query(
		`WITH t1 AS (SELECT UNNEST($1::varchar[]) AS name)
		 SELECT id, levenshtein(LOWER(medicine.name), LOWER(t1.name)) AS ed
		 FROM t1, medicine
		 WHERE levenshtein(LOWER(medicine.name), LOWER(t1.name)) <= LENGTH(t1.name)/2
		 ORDER BY ed ASC limit $2`,
		pq.Array(names), len(names),
	)
	if err != nil {
		return nil, err
	}

	ids, err := itToSlice(rows, func(it Iterator) (id int64, err error) {
		var d int
		err = it.Scan(&id, &d)
		return
	})

	return lo.Uniq(ids), err
}

func AptekaSearch(tx TransactionOps, mIds []int64) ([]lo.Tuple2[model.Apteka, []string], error) {
	rows, err := tx.Query(
		`WITH t1 as (
		   SELECT apteka_id, ARRAY_AGG(medicine.name)::varchar[]
		   FROM apteka_medicine LEFT JOIN medicine ON medicine_id = medicine.id
		   WHERE medicine_id = ANY($1::bigint[]) AND amount > 0
		   GROUP BY apteka_id
		 )
		 SELECT name, phone, address, array_agg FROM t1
		 LEFT JOIN apteka ON apteka_id = apteka.id
		 ORDER BY array_length(array_agg, 1) DESC`,
		pq.Array(mIds),
	)
	if err != nil {
		return nil, err
	}

	rs, err := itToSlice(rows, func(it Iterator) (r lo.Tuple2[model.Apteka, []string], err error) {
		err = it.Scan(&r.A.Name, &r.A.Phone, &r.A.Address, pq.Array(&r.B))
		return
	})

	return rs, err
}
