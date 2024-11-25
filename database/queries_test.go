package database

import (
	"testing"

	"github.com/oybek/choguuket/model"
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
	}

	t.Run("AptekaInsert", func(t *testing.T) {
		id, err := AptekaInsert(tx, &apteka)
		assert.NoError(t, err)
		assert.Equal(t, 1, id)
	})

	t.Run("SelectUser", func(t *testing.T) {
		err := UserInsert(tx, &user)
		assert.NoError(t, err)
	})
}
