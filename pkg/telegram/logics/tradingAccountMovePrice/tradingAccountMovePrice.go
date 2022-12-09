package tradingAccountMovePrice

import (
	"database/sql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rob-bender/invest-market-frontend/pkg/telegram/keyboard"
)

func TradingAccountMovePrice(bot *tgbotapi.BotAPI, sqliteDb *sql.DB, msg tgbotapi.MessageConfig, teleId int64, userName string, languageUser string) error {
	if len(languageUser) > 0 {
		if languageUser == "ru" {
			msg.Text = "💭 *Куда пойдет цена актива?*"
			msg.ReplyMarkup = keyboard.GenKeyboardInlineForTradingAccountMovePrice("☝️ Вверх", "👇 Вниз", "⚖️  Не изменится", "❌ Отмена")
		}
		if languageUser == "en" {
			msg.Text = "💭 *What direction will the asset price go?*"
			msg.ReplyMarkup = keyboard.GenKeyboardInlineForTradingAccountMovePrice("☝️ Up", "👇 Down", "⚖️  Won't change", "❌ Cancel")
		}
		_, err := bot.Send(msg)
		if err != nil {
			return err
		}
	}

	return nil
}
