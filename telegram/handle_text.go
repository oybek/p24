package telegram

import (
	"fmt"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/google/uuid"
	"github.com/jellydator/ttlcache/v3"
	"github.com/oybek/choguuket/database"
	"github.com/oybek/choguuket/model"
	"github.com/samber/lo"
)

func (lp *LongPoll) handleText(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveMessage.Chat
	text := ctx.EffectiveMessage.Text
	return lp.searchByText(chat.Id, text)
}

func (lp *LongPoll) searchByText(chatId int64, text string) error {
	medicineNames := lo.Map(strings.Split(text, ","), func(s string, _ int) string {
		return strings.TrimSpace(s)
	})

	tuples, err := lp.searchApteka(medicineNames)
	if err != nil {
		return err
	}

	if len(tuples) == 0 {
		return lp.sendText(
			chatId,
			fmt.Sprintf("Не нашел данные лекарства ни в одной из аптек: %s\n\n"+
				"Если нужна помощь наберите команду /help", text),
		)
	}

	requestId := uuid.New()
	lp.requestCache.Set(requestId, tuples, ttlcache.DefaultTTL)
	url := "https://wolfrepos.github.io/apteka/search/index.html?request_id=" + requestId.String()
	_, err = lp.bot.SendMessage(
		chatId, fmt.Sprintf("По запросу: %s\n\nНайдено %d аптек", text, len(tuples)),
		&gotgbot.SendMessageOpts{
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{{Text: "Открыть список", Url: url}},
				},
			},
		},
	)

	return err
}

func (lp *LongPoll) searchApteka(medicineNames []string) ([]lo.Tuple2[model.Apteka, []string], error) {
	return database.Transact(
		lp.db,
		func(tx database.TransactionOps) ([]lo.Tuple2[model.Apteka, []string], error) {
			mIds, err := database.MedicineSearch(tx, medicineNames)
			if err != nil {
				return nil, err
			}

			return database.AptekaSearch(tx, mIds)
		},
	)
}
