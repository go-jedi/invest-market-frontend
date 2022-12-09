package tradingAccountWaitTime

import (
	"database/sql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rob-bender/invest-market-frontend/pkg/telegram/keyboard"
	requestProject "github.com/rob-bender/invest-market-frontend/pkg/telegram/request"
)

func TradingAccountWaitTime(bot *tgbotapi.BotAPI, sqliteDb *sql.DB, msg tgbotapi.MessageConfig, teleId int64, userName string, languageUser string, chooseMovePrice string) error {
	var userChooseMovePrice string = chooseMovePrice
	if chooseMovePrice == "dn" {
		userChooseMovePrice = "down"
	}
	if chooseMovePrice == "dc" {
		userChooseMovePrice = "not change"
	}
	if len(languageUser) > 0 {
		resUpdateTradingMovePrice, err := requestProject.UpdateTradingMovePrice(teleId, userChooseMovePrice)
		if err != nil {
			return err
		}
		if resUpdateTradingMovePrice {
			if languageUser == "ru" {
				msg.Text = "⏰ Выберите время ожидания."
				msg.ReplyMarkup = keyboard.GenKeyboardInlineForTradingAcWaitTime("10 секунд [x1.3]", "30 секунд [x1.5]", "60 секунд [x2.0]", "❌ Отмена")
			}
			if languageUser == "en" {
				msg.Text = "⏰ Выберите время ожидания."
				msg.ReplyMarkup = keyboard.GenKeyboardInlineForTradingAcWaitTime("10 seconds [x1.3]", "30 seconds [x1.5]", "60 seconds [x2.0]", "❌ Cancel")
			}
			_, err := bot.Send(msg)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
