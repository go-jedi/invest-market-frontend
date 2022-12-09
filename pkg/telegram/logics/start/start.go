package start

import (
	"database/sql"
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rob-bender/invest-market-frontend/pkg/telegram/helperFunc"
	"github.com/rob-bender/invest-market-frontend/pkg/telegram/keyboard"
	requestProject "github.com/rob-bender/invest-market-frontend/pkg/telegram/request"
	"github.com/rob-bender/invest-market-frontend/pkg/telegram/sqlite"
)

func GetStart(bot *tgbotapi.BotAPI, sqliteDb *sql.DB, msg tgbotapi.MessageConfig, teleId int64, userName string, teleIdAdmin int64) error {
	resCheckAuth, err := requestProject.CheckAuth(teleId)
	if err != nil {
		return err
	}
	if resCheckAuth {
		resGetUserLanguage, err := requestProject.GetUserLanguage(teleId)
		if err != nil {
			return err
		}
		if len(resGetUserLanguage[0].Lang) > 0 {
			resGetUserCurrency, err := requestProject.GetUserCurrency(teleId)
			if err != nil {
				return err
			}
			if len(resGetUserCurrency[0].Currency) > 0 {
				resCheckIsTerms, err := requestProject.CheckIsTerms(teleId)
				if err != nil {
					return err
				}
				if resCheckIsTerms {
					var textTwo string = ""
					if resGetUserLanguage[0].Lang == "ru" {
						textTwo = "Если у вас не появилось меню, то напишите /start"
					}
					if resGetUserLanguage[0].Lang == "en" {
						textTwo = "In case the menu did not appear, send /start or press it"
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
						if resGetUserLanguage[0].Lang == "ru" {
							msg.ReplyMarkup = keyboard.GenKeyboardHomeAdmin("Мой торговый счёт", "Личный кабинет", "Информация", "Поддержка")
							msg.Text = "Главное меню"
						}
						if resGetUserLanguage[0].Lang == "en" {
							msg.ReplyMarkup = keyboard.GenKeyboardHomeAdmin("My trading account", "Profile", "About", "Support")
							msg.Text = "Main menu"
						}
					} else {
						if resGetUserLanguage[0].Lang == "ru" {
							msg.ReplyMarkup = keyboard.GenKeyboardHome("Мой торговый счёт", "Личный кабинет", "Информация", "Поддержка")
							msg.Text = "Главное меню"
						}
						if resGetUserLanguage[0].Lang == "en" {
							msg.ReplyMarkup = keyboard.GenKeyboardHome("My trading account", "Profile", "About", "Support")
							msg.Text = "Main menu"
						}
					}
					_, err = bot.Send(msg)
					if err != nil {
						return err
					}
				} else {
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
					if resGetUserLanguage[0].Lang == "ru" {
						photo.ParseMode = "HTML"
						var text string = "Пожалуйста, введите цифры с картинки.\n\nНеоднократное введение неправильных значений может привести к блокировке.\n\nВводя капчу, вы принимаете условия <a href='https://telegra.ph/Polzovatelskoe-soglashenie-Open-Sea-06-18'>пользовательского соглашения</a>."
						photo.Caption = text
						_, err := bot.Send(photo)
						if err != nil {
							return err
						}
					}
					if resGetUserLanguage[0].Lang == "en" {
						photo.ParseMode = "HTML"
						var text string = "Please, solve captcha.\n\nConsequent wrong answers may lead to a ban."
						photo.Caption = text
						_, err := bot.Send(photo)
						if err != nil {
							return err
						}
					}
				}
			} else {
				if resGetUserLanguage[0].Lang == "ru" {
					msg.Text = "Выберите свою валюту"
				}
				if resGetUserLanguage[0].Lang == "en" {
					msg.Text = "Choose your currency"
				}
				msg.ReplyMarkup = keyboard.DgCurrencyKeyboardInline
				_, err := bot.Send(msg)
				if err != nil {
					return err
				}
			}
		} else {
			var text string = "🏳️?"
			msg.ReplyMarkup = keyboard.DgLangKeyboardInline
			msg.Text = text
			_, err := bot.Send(msg)
			if err != nil {
				return err
			}
		}
	} else {
		resRegisterUser, err := requestProject.RegisterUser(teleId, userName, teleIdAdmin)
		if err != nil {
			return err
		}
		if resRegisterUser {
			var text string = "🏳️?"
			msg.ReplyMarkup = keyboard.DgLangKeyboardInline
			msg.Text = text
			_, err := bot.Send(msg)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
