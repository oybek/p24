package database

import (
	"testing"
	"time"

	"github.com/oybek/choguuket/model"
	"github.com/stretchr/testify/assert"
)

func TestQueries(t *testing.T) {
	tx := testdb.Conn
	chatId := int64(123)
	tomorrow := time.Now().UTC().Truncate(time.Second).AddDate(0, 0, 1)
	trip := model.Trip{
		ChatId:    chatId,
		Path:      []string{"Бишкек", "Балыкчы", "Бостери", "Каракол"},
		StartTime: tomorrow,
		Phone:     "123",
		Comment:   "kia k5, 3 места, Айбек",
	}
	tripReq := model.TripReq{
		ChatId:    chatId,
		From:      "Бишкек",
		To:        "Каракол",
		StartDate: tomorrow,
	}

	t.Run("InsertTrip", func(t *testing.T) {
		_, err := InsertTrip(tx, trip)
		assert.NoError(t, err)
	})

	t.Run("InsertTripReq", func(t *testing.T) {
		_, err := InsertTripReq(tx, tripReq)
		assert.NoError(t, err)
	})

	t.Run("SearchTrip", func(t *testing.T) {
		trips, err := SearchTrip(tx, tripReq)

		assert.NoError(t, err)
		assert.Equal(t, 1, len(trips))
		assert.Equal(t, trip.ChatId, trips[0].ChatId)
		assert.Equal(t, trip.Path, trips[0].Path)
		assert.Equal(t, trip.StartTime.String(), trips[0].StartTime.String())
		assert.Equal(t, trip.Phone, trips[0].Phone)
		assert.Equal(t, trip.Comment, trips[0].Comment)
	})

	t.Run("SearchTrip", func(t *testing.T) {
		trips, err := SearchTrip(tx, model.TripReq{
			From:      "Бишкек",
			To:        "Чолпон-ата",
			StartDate: tomorrow,
		})

		assert.NoError(t, err)
		assert.Equal(t, 0, len(trips))
	})

	t.Run("SearchTripReq", func(t *testing.T) {
		tripReqs, err := SearchTripReq(tx, trip)

		assert.NoError(t, err)
		assert.Equal(t, 1, len(tripReqs))
		assert.Equal(t, tripReq.From, tripReqs[0].From)
		assert.Equal(t, tripReq.To, tripReqs[0].To)
		assert.Equal(t, tripReq.StartDate.String(), tripReqs[0].StartDate.String())
	})

	t.Run("SearchTripReq", func(t *testing.T) {
		tripReqs, err := SearchTripReq(tx, model.Trip{
			ChatId:    chatId,
			Path:      []string{"Бишкек", "Балыкчы", "Бостери", "Чолпон-ата"},
			StartTime: tomorrow,
			Phone:     "123",
			Comment:   "kia k5, 3 места, Айбек",
		})

		assert.NoError(t, err)
		assert.Equal(t, 0, len(tripReqs))
	})
}
