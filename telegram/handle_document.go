package telegram

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/oybek/choguuket/database"
	"github.com/oybek/choguuket/model"
)

func (lp *LongPoll) handleDocument(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveMessage.Chat
	doc := ctx.EffectiveMessage.Document

	docReader, err := lp.downloadDocument(doc)
	if err != nil {
		return err
	}
	defer docReader.Close()

	user, err := lp.readUser(chat.Id)
	if err != nil {
		return err
	}

	excelReader, ok := lp.readers[user.Reader]
	if !ok {
		return fmt.Errorf("reader '%s' not found for chatId=%d", user.Reader, chat.Id)
	}

	medicines, err := excelReader.Read(docReader)
	if err != nil {
		return err
	}

	updn, err := lp.writeMedicines(user.AptekaId, medicines)
	if err != nil {
		return err
	}

	return lp.sendText(chat.Id, fmt.Sprintf("Обновлено %d записей", updn))
}

func (lp *LongPoll) writeMedicines(aptekaId int64, medicines []model.Medicine) (int, error) {
	return database.Transact(lp.db, func(tx database.TransactionOps) (counter int, err error) {
		now := time.Now()
		counter = 0
		for _, medicine := range medicines {
			medicineId, err := database.MedicineInsert(tx, &medicine)
			if err != nil {
				continue
			}

			err = database.AptekaMedicineInsert(tx, aptekaId, medicineId, medicine.Amount, now)
			if err != nil {
				continue
			}

			counter++
		}

		return counter, err
	})
}

func (lp *LongPoll) readUser(chatId int64) (model.User, error) {
	return database.Transact(lp.db, func(tx database.TransactionOps) (model.User, error) {
		user, err := database.UserSelect(tx, chatId)
		return user, err
	})
}

func (lp *LongPoll) downloadDocument(document *gotgbot.Document) (io.ReadCloser, error) {
	file, err := lp.bot.GetFile(document.FileId, &gotgbot.GetFileOpts{})
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(file.URL(lp.bot, &gotgbot.RequestOpts{}))
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
