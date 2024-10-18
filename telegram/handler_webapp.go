package telegram

import (
	"fmt"
	"log"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/jellydator/ttlcache/v3"
	"github.com/oybek/choguuket/database"
	"github.com/oybek/choguuket/model"
)

func (lp *LongPoll) handleWebAppData(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := &ctx.EffectiveMessage.Chat
	data := ctx.EffectiveMessage.WebAppData.Data

	log.Printf("[ChatId=%d] Got webapp data: %s", chat.Id, data)

	if trip, err := parse[model.Trip](data); err == nil {
		return lp.handleNewTrip(chat, trip)
	}
	if tripReq, err := parse[model.TripReq](data); err == nil {
		return lp.handleNewTripReq(chat, tripReq)
	} else {
		log.Printf("error parsing: %s", err)
	}

	return nil
}

func (lp *LongPoll) handleNewTrip(chat *gotgbot.Chat, trip *model.Trip) error {
	_, err := database.Transact(lp.db, func(tx database.TransactionOps) (any, error) {
		return database.InsertTrip(tx, trip)
	})
	if err != nil {
		return fmt.Errorf("failed to insert trip: %w", err)
	}

	err = lp.sendText(chat.Id, "Поездка создана ✅")
	if err != nil {
		return err
	}

	time.Sleep(300 * time.Millisecond)

	go lp.notifyChatsAboutTrip(trip)

	return lp.sendText(chat.Id, trip.String())
}

func (lp *LongPoll) handleNewTripReq(chat *gotgbot.Chat, tripReq *model.TripReq) error {
	tripReqId, err := database.Transact(lp.db, func(tx database.TransactionOps) (int64, error) {
		return database.InsertTripReq(tx, tripReq)
	})
	if err != nil {
		return fmt.Errorf("failed to insert trip: %w", err)
	}

	err = lp.sendText(chat.Id, "Ищу поездки по запросу:\n"+tripReq.String())
	if err != nil {
		return err
	}

	time.Sleep(time.Second)

	trips, err := database.Transact(lp.db, func(tx database.TransactionOps) ([]model.Trip, error) {
		return database.SearchTrip(tx, tripReq)
	})
	if err != nil {
		return err
	}

	if len(trips) == 0 {
		return lp.sendText(chat.Id, "Пока нет поездок по Вашему запросу, как только появится поездка я Вам сообщу")
	}

	lp.searchCache.Set(tripReqId, trips, ttlcache.DefaultTTL)
	return lp.sendTrip(chat, tripReqId, trips)
}

func (lp *LongPoll) notifyChatsAboutTrip(trip *model.Trip) error {
	tripReqs, err := database.Transact(lp.db, func(tx database.TransactionOps) ([]model.TripReq, error) {
		return database.SearchTripReq(tx, trip)
	})
	if err != nil {
		log.Printf("error SearchTripReq %s", err)
		return err
	}

	log.Printf("Have to notify: %#v", tripReqs)

	for _, tripReq := range tripReqs {
		lp.sendTripNotification(&tripReq, trip)
		time.Sleep(2 * time.Second)
	}

	return nil
}
