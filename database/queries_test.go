package database

import (
	"testing"

	"github.com/google/uuid"
	"github.com/oybek/choguuket/model"
	"github.com/stretchr/testify/assert"
)

func TestQueries(t *testing.T) {
	tx := testdb.Conn
	chatId := int64(123)
	uuid0 := uuid.New()
	user := model.User{
		ChatId: chatId,
		UUID:   uuid0,
	}

	t.Run("UpsertUser", func(t *testing.T) {
		_, err := UpsertUser(tx, &user)
		assert.NoError(t, err)
	})

	t.Run("SelectUser", func(t *testing.T) {
		user, err := SelectUser(tx, uuid0)
		assert.NoError(t, err)
		assert.Equal(t, model.User{ChatId: chatId, UUID: uuid0}, *user)
	})
}
