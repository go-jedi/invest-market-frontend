package changeMamLuck

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rob-bender/invest-market-frontend/pkg/telegram/keyboard"
	requestProject "github.com/rob-bender/invest-market-frontend/pkg/telegram/request"
)

func ChangeMamLuck(bot *tgbotapi.BotAPI, msg tgbotapi.MessageConfig, teleId int64, userName string, languageUser string, userChooseTeleId int64) error {
	if len(languageUser) > 0 {
		fmt.Println("222222222222222222222")
		resGetUserCurrentLuck, err := requestProject.GetUserCurrentLuck(userChooseTeleId)
		if err != nil {
			return err
		}
		fmt.Println("resGetUserCurrentLuck -->", resGetUserCurrentLuck)
		if len(resGetUserCurrentLuck) > 0 {
			var currentLuck string = ""
			if resGetUserCurrentLuck[0].CurrentLuck == "win" {
				currentLuck = "Выигрыш"
			}
			if resGetUserCurrentLuck[0].CurrentLuck == "loss" {
				currentLuck = "Проигрыш"
			}
			if resGetUserCurrentLuck[0].CurrentLuck == "random" {
				currentLuck = "Рандом"
			}
			msg.Text = fmt.Sprintf("Текущая удача мамонта: *%s*\n\nВыбери новую:",
				currentLuck,
			)
			msg.ReplyMarkup = keyboard.GenKeyboardInlineForChangeMamLuck(userChooseTeleId)
			_, err := bot.Send(msg)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
