package homeAfterReg

import (
	"database/sql"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rob-bender/invest-market-frontend/pkg/telegram/keyboard"
	requestProject "github.com/rob-bender/invest-market-frontend/pkg/telegram/request"
)

func HomeAfterReg(bot *tgbotapi.BotAPI, sqliteDb *sql.DB, msg tgbotapi.MessageConfig, teleId int64, userName string, isAgreeTerms bool, languageUser string) error {
	if isAgreeTerms {
		resAgreeTerms, err := requestProject.AgreeTerms(teleId)
		if err != nil {
			return err
		}
		if resAgreeTerms {
			var textTwo string = ""
			if languageUser == "ru" {
				textTwo = "✅ Капча введена верно. Можете спокойно пользоваться ботом.\n\nЕсли у вас не появилось меню, то напишите /start"
			}
			if languageUser == "en" {
				textTwo = "✅ Captcha entered correctly. Feel free to use the bot.\n\nIf you don't see a menu, type /start"
			}
			msg.Text = textTwo
			_, err = bot.Send(msg)
			if err != nil {
				return err
			}
			resCheckIsAdmin, err := requestProject.CheckIsAdmin(teleId)
			if err != nil {
				return err
			}
			if resCheckIsAdmin {
				if languageUser == "ru" {
					msg.ReplyMarkup = keyboard.GenKeyboardHomeAdmin("Мой торговый счёт", "Личный кабинет", "Информация", "Поддержка")
					msg.Text = "Главное меню"
				}
				if languageUser == "en" {
					msg.ReplyMarkup = keyboard.GenKeyboardHomeAdmin("My trading account", "Profile", "About", "Support")
					msg.Text = "Main menu"
				}
				_, err = bot.Send(msg)
				if err != nil {
					return err
				}
			} else {
				resGetAdminByUser, err := requestProject.GetAdminByUser(teleId)
				if err != nil {
					return err
				}
				if len(resGetAdminByUser) > 0 {
					msg.ChatID = resGetAdminByUser[0].TeleId
					msg.ParseMode = "HTML"
					msg.Text = fmt.Sprintf("✅ Мамонт %s @%s /u%d) <b>зарегистрировался по твоей ссылке</b>", userName, userName, teleId)
					_, err := bot.Send(msg)
					if err != nil {
						return err
					}
				}
				if languageUser == "ru" {
					msg.ChatID = teleId
					msg.ReplyMarkup = keyboard.GenKeyboardHome("Мой торговый счёт", "Личный кабинет", "Информация", "Поддержка")
					msg.Text = "Главное меню"
				}
				if languageUser == "en" {
					msg.ChatID = teleId
					msg.ReplyMarkup = keyboard.GenKeyboardHome("My trading account", "Profile", "About", "Support")
					msg.Text = "Main menu"
				}
				_, err = bot.Send(msg)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
