package chooseCurrency

import (
	"database/sql"
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rob-bender/invest-market-frontend/pkg/telegram/helperFunc"
	requestProject "github.com/rob-bender/invest-market-frontend/pkg/telegram/request"
	"github.com/rob-bender/invest-market-frontend/pkg/telegram/sqlite"
)

func ChooseCurrency(bot *tgbotapi.BotAPI, sqliteDb *sql.DB, msg tgbotapi.MessageConfig, teleId int64, userName string, languageChoose string) error {
	if len(languageChoose) > 0 {
		resUpdateLanguage, err := requestProject.UpdateLanguage(teleId, languageChoose)
		if err != nil {
			return err
		}
		if resUpdateLanguage {
			_, err = sqliteDb.Exec("INSERT INTO bot_params(tele_id, lang) VALUES($1, $2)", teleId, languageChoose)
			if err != nil {
				return err
			}
			err := sqlite.TurnOnListenerWatchingWriteCaptcha(sqliteDb, teleId)
			if err != nil {
				return err
			}
			var isTesting string = os.Getenv("IS_TESTING")
			var needPath string = ""
			if isTesting == "true" {
				needPath = "/home/dale/job/work/my-project/invest-market/frontend/img/captcha/"
			} else {
				needPath = "/home/nft-market-bot/frontend/invest-market-frontend/img/captcha/"
			}
			var resRandomRangeInt int = helperFunc.RandomRangeInt(1, 9)
			photo := tgbotapi.NewPhoto(teleId, tgbotapi.FilePath(fmt.Sprintf("%s%d.jpg", needPath, resRandomRangeInt)))
			photo.ParseMode = "Markdown"
			if languageChoose == "ru" {
				photo.ParseMode = "HTML"
				var text string = "Пожалуйста, введите цифры с картинки.\n\nНеоднократное введение неправильных значений может привести к блокировке.\n\nВводя капчу, вы принимаете условия <a href='https://telegra.ph/Polzovatelskoe-soglashenie-Open-Sea-06-18'>пользовательского соглашения</a>."
				photo.Caption = text
				_, err := bot.Send(photo)
				if err != nil {
					return err
				}
			}
			if languageChoose == "en" {
				photo.ParseMode = "HTML"
				var text string = "Please, solve captcha.\n\nConsequent wrong answers may lead to a ban."
				photo.Caption = text
				_, err := bot.Send(photo)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
