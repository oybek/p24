package telegram

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/jellydator/ttlcache/v3"
	"github.com/oybek/choguuket/model"
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

func (lp *LongPoll) sendTrip(chat *gotgbot.Chat, tripReqId int64, trips []model.Trip) error {
	sendMessageOpts := gotgbot.SendMessageOpts{
		ParseMode: "markdown",
	}

	if len(trips) > 1 {
		sendMessageOpts.ReplyMarkup = gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
				gotgbot.InlineKeyboardButton{
					Text:         fmt.Sprintf("Еще %d", len(trips)-1),
					CallbackData: fmt.Sprintf("next;%d", tripReqId),
				},
			}},
		}
	}

	_, err := lp.bot.SendMessage(chat.Id, trips[0].String(), &sendMessageOpts)
	return err
}
