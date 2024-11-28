package database

import (
	"testing"
	"time"

	"github.com/oybek/choguuket/model"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestQueries(t *testing.T) {
	tx := testdb.Conn
	chatId := int64(108349719)
	apteka := model.Apteka{
		Name:    "Фармамир",
		Phone:   "0559171775",
		Address: "Токтоналиева 61",
	}
	user := model.User{
		ChatId:   chatId,
		AptekaId: 1,
		Reader:   "test",
	}

	t.Run("AptekaInsert", func(t *testing.T) {
		id, err := AptekaInsert(tx, &apteka)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)
	})

	t.Run("SelectUser", func(t *testing.T) {
		err := UserInsert(tx, &user)
		assert.NoError(t, err)
	})

	t.Run("SelectUser", func(t *testing.T) {
		user, err := UserSelect(tx, chatId)
		assert.NoError(t, err)
		assert.Equal(t, model.User{ChatId: chatId, AptekaId: 1, Reader: "test"}, user)
	})

	t.Run("MedicineInsert", func(t *testing.T) {
		medicine := model.Medicine{
			Name:   "Анальгин",
			Amount: 10,
		}
		id, err := MedicineInsert(tx, &medicine)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)
	})

	t.Run("MedicineSearch", func(t *testing.T) {
		ids, err := MedicineSearch(tx, []string{"Аналгин", "Парацетамол"})
		assert.NoError(t, err)
		assert.Equal(t, []int64{1}, ids)
	})

	t.Run("AptekaMedicineInsert", func(t *testing.T) {
		err := AptekaMedicineInsert(tx, 1, 1, 10, time.Now())
		assert.NoError(t, err)

		err = AptekaMedicineInsert(tx, 1, 1, 10, time.Now())
		assert.NoError(t, err)
	})

	t.Run("AptekaSearch", func(t *testing.T) {
		rs, err := AptekaSearch(tx, []int64{1})
		expected := []lo.Tuple2[model.Apteka, []string]{
			lo.T2(apteka, []string{"Анальгин"}),
		}
		assert.NoError(t, err)
		assert.Equal(t, 1, len(rs))
		assert.Equal(t, expected, rs)
	})
}
