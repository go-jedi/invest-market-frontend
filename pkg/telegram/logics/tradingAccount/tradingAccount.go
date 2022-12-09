package tradingAccount

import (
	"database/sql"
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rob-bender/invest-market-frontend/pkg/telegram/keyboard"
	"github.com/rob-bender/invest-market-frontend/pkg/telegram/sqlite"
)

func TradingAccount(bot *tgbotapi.BotAPI, sqliteDb *sql.DB, msg tgbotapi.MessageConfig, teleId int64, userName string, languageUser string) error {
	err := sqlite.TurnOffListeners(sqliteDb, teleId)
	if err != nil {
		return err
	}
	if len(languageUser) > 0 {
		var isTesting string = os.Getenv("IS_TESTING")
		var needPath string = ""
		if isTesting == "true" {
			needPath = "/home/dale/job/work/my-project/invest-market/frontend/img"
		} else {
			needPath = "/home/nft-market-bot/frontend/invest-market-frontend/img"
		}
		photo := tgbotapi.NewPhoto(teleId, tgbotapi.FilePath(fmt.Sprintf("%s%s", needPath, "/img-need/1.jpg")))
		photo.ParseMode = "Markdown"
		if languageUser == "ru" {
			photo.Caption = "*📈 Ваш личный торговый-счет*\n\n💠 Выберите актив для торговли:"
		}
		if languageUser == "en" {
			photo.Caption = "*📈 Your personal trading-account*\n\n💠 Choose an asset to trade:"
		}
		photo.ReplyMarkup = keyboard.GenKeyboardInlineForTradingAccount()
		_, err := bot.Send(photo)
		if err != nil {
			return err
		}
	}

	return nil
}
