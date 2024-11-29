package telegram

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/google/uuid"
	"github.com/jellydator/ttlcache/v3"
	"github.com/oybek/choguuket/database"
	"github.com/oybek/choguuket/model"
	"github.com/samber/lo"
	"github.com/sashabaranov/go-openai"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func (lp *LongPoll) handleVoice(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveMessage.Chat
	voice := ctx.EffectiveMessage.Voice

	if voice.Duration > 20 {
		return lp.sendText(chat.Id, TextTooLongVoice)
	}

	lp.sendText(chat.Id, "Ищу подходящие аптеки")

	text, err := lp.transcribeVoice(voice)
	if err != nil {
		return err
	}

	medicineNames := lo.Map(strings.Split(text, ","), func(s string, _ int) string {
		return strings.TrimSpace(s)
	})

	tuples, err := lp.searchApteka(medicineNames)
	if err != nil {
		return err
	}

	if len(tuples) == 0 {
		return lp.sendText(
			chat.Id,
			fmt.Sprintf("Не нашел данные лекарства ни в одной из аптек: %s", text),
		)
	}

	requestId := uuid.New()
	lp.requestCache.Set(requestId, tuples, ttlcache.DefaultTTL)
	url := "https://wolfrepos.github.io/apteka/search/index.html?request_id=" + requestId.String()
	_, err = lp.bot.SendMessage(
		chat.Id, fmt.Sprintf("По запросу: %s\nНайдено %d аптек", text, len(tuples)),
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

func (lp *LongPoll) transcribeVoice(voice *gotgbot.Voice) (string, error) {
	file, err := lp.bot.GetFile(voice.FileId, &gotgbot.GetFileOpts{})
	if err != nil {
		return "", err
	}

	resp, err := http.Get(file.URL(lp.bot, &gotgbot.RequestOpts{}))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		Reader:   resp.Body,
		FilePath: file.FilePath,
		Prompt:   "Парацетамол, ТайлолХот, Тримол",
		Language: "ru",
	}

	context := context.Background()
	openaiResp, err := lp.openaiClient.CreateTranscription(context, req)
	if err != nil {
		fmt.Printf("Transcription error: %v\n", err)
		return "", err
	}

	return openaiResp.Text, nil
}

func toMessage(t lo.Tuple2[model.Apteka, []string]) (string, *gotgbot.SendMessageOpts) {
	a, ms := t.A, t.B
	baseInfo := fmt.Sprintf(
		"%s %s\n%s %s\n%s %s",
		EmojiHospital, cases.Title(language.Und).String(a.Name),
		EmojiPin, cases.Title(language.Und).String(a.Address),
		EmojiPhone, a.Phone,
	)
	presenceInfo := fmt.Sprintf("В наличии: %s", strings.Join(ms, ", "))

	return baseInfo + "\n\n" + presenceInfo, &gotgbot.SendMessageOpts{}
}
