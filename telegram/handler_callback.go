package telegram

import (
	"errors"
	"strconv"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/jellydator/ttlcache/v3"
	"github.com/oybek/choguuket/database"
)

func (lp *LongPoll) handleNextTrip(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat

	cb := ctx.Update.CallbackQuery
	cbs := strings.Split(cb.Data, ";")
	if len(cbs) != 2 {
		return errors.New("invalid callback data")
	}

	tripReqId, err := strconv.ParseInt(cbs[1], 10, 64)
	if err != nil {
		return err
	}

	cb.Message.EditReplyMarkup(b, &gotgbot.EditMessageReplyMarkupOpts{
		InlineMessageId: cb.InlineMessageId,
		ReplyMarkup:     gotgbot.InlineKeyboardMarkup{},
	})

	kv := lp.searchCache.Get(tripReqId)
	trips := kv.Value()

	if len(trips) == 0 {
		return err
	}

	err = lp.sendTrip(chat, tripReqId, trips)
	if err != nil {
		return err
	}

	if len(trips) == 1 {
		lp.searchCache.Delete(tripReqId)
	} else {
		lp.searchCache.Set(tripReqId, trips[1:], ttlcache.DefaultTTL)
	}

	return nil
}

func (lp *LongPoll) handleRmTripReq(b *gotgbot.Bot, ctx *ext.Context) error {
	cb := ctx.Update.CallbackQuery
	cbs := strings.Split(cb.Data, ";")
	if len(cbs) != 2 {
		return errors.New("invalid callback data")
	}

	tripReqId, err := strconv.Atoi(cbs[1])
	if err != nil {
		return err
	}

	cb.Message.EditReplyMarkup(b, &gotgbot.EditMessageReplyMarkupOpts{
		InlineMessageId: cb.InlineMessageId,
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
				gotgbot.InlineKeyboardButton{
					Text:         "Вы отписались от рассылки",
					CallbackData: "null",
				},
			}},
		},
	})

	_, err = database.Transact(lp.db, func(tx database.TransactionOps) (any, error) {
		return database.DeleteTripReq(tx, tripReqId)
	})
	return err
}
