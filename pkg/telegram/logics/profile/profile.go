package profile

import (
	"database/sql"
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rob-bender/invest-market-frontend/pkg/telegram/keyboard"
	requestProject "github.com/rob-bender/invest-market-frontend/pkg/telegram/request"
	"github.com/rob-bender/invest-market-frontend/pkg/telegram/sqlite"
)

func Profile(bot *tgbotapi.BotAPI, sqliteDb *sql.DB, msg tgbotapi.MessageConfig, teleId int64, userName string, languageUser string) error {
	resGetUserProfile, err := requestProject.GetUserProfile(teleId)
	if err != nil {
		return err
	}
	if len(resGetUserProfile) > 0 {
		err = sqlite.TurnOffListeners(sqliteDb, teleId)
		if err != nil {
			return err
		}
		resGetUserTotalEarn, err := requestProject.GetUserTotalEarn(teleId)
		if err != nil {
			return err
		}
		if len(resGetUserTotalEarn) > 0 {
			var isTesting string = os.Getenv("IS_TESTING")
			var needPath string = ""
			if isTesting == "true" {
				needPath = "/home/dale/job/work/my-project/invest-market/frontend/img"
			} else {
				needPath = "/home/nft-market-bot/frontend/invest-market-frontend/img"
			}
			photo := tgbotapi.NewPhoto(teleId, tgbotapi.FilePath(fmt.Sprintf("%s%s", needPath, "/img-need/1.jpg")))
			photo.ParseMode = "Markdown"
			if err != nil {
				return err
			}
			resCheckIsAdmin, err := requestProject.CheckIsAdmin(teleId)
			if err != nil {
				return err
			}
			if resCheckIsAdmin {
				resCheckIsVisibleName, err := requestProject.CheckIsVisibleName(teleId)
				if err != nil {
					return err
				}
				if languageUser == "ru" {
					var isVerification string
					var isPremium string
					var isNickName string
					if resGetUserProfile[0].Verification {
						isVerification = "✅ *Верифицирован*"
					} else {
						isVerification = "❌ *Не верифицирован*"
					}
					if resGetUserProfile[0].IsPremium {
						isPremium = "💎 *Премиум*"
					} else {
						isPremium = "⚠️ *Обычный*"
					}
					if resCheckIsVisibleName {
						isNickName = "🪫 Скрыть никнейм в выплатах"
					} else {
						isNickName = "🔋 Показать никнейм в выплатах"
					}
					photo.Caption = fmt.Sprintf("*Личный кабинет*\n\n💰 Баланс: *%.2f $*\n📤 На выводе: *%.2f $*\n💵 Заработано: *%.2f $*\n\n👤 Верификация:  %s\n🌐 Статус аккаунта:  %s\n🆔 Ваш ID: [%d](tg://user?id=%d)",
						resGetUserProfile[0].Balance,
						resGetUserProfile[0].Conclusion,
						resGetUserTotalEarn[0].TotalEarn,
						isVerification,
						isPremium,
						teleId,
						teleId,
					)
					photo.ReplyMarkup = keyboard.GenKeyboardInlineForProfileMenuAdmin("📥 Пополнить", "📤 Вывести", "📝 Верификация", "🇺🇸 English language", "en", isNickName)
				}

				if languageUser == "en" {
					var isVerification string
					var isPremium string
					var isNickName string
					if resGetUserProfile[0].Verification {
						isVerification = "✅ *Verified*"
					} else {
						isVerification = "❌ *Not verified*"
					}
					if resGetUserProfile[0].IsPremium {
						isPremium = "💎 *Premium*"
					} else {
						isPremium = "⚠️ *Ordinary*"
					}
					if resCheckIsVisibleName {
						isNickName = "🪫 Hide nickname in payouts"
					} else {
						isNickName = "🔋 Show nickname in payouts"
					}
					photo.Caption = fmt.Sprintf("*Personal account*\n\n💰 Balance: *%.2f $*\n📤 Withdrawal: *%.2f $*\n💵 Earned: *%.2f $*\n\n👤 Verification:  %s\n🌐 Status Account:  %s\n🆔 Your ID: [%d](tg://user?id=%d)",
						resGetUserProfile[0].Balance,
						resGetUserProfile[0].Conclusion,
						resGetUserTotalEarn[0].TotalEarn,
						isVerification,
						isPremium,
						teleId,
						teleId,
					)
					photo.ReplyMarkup = keyboard.GenKeyboardInlineForProfileMenuAdmin("📥 Deposit", "📤 Withdraw", "📝 Verification", "🇷🇺 Русский язык", "ru", isNickName)
				}

				_, err = bot.Send(photo)
				if err != nil {
					return err
				}
			} else {
				if languageUser == "ru" {
					var isVerification string
					var isPremium string
					if resGetUserProfile[0].Verification {
						isVerification = "✅ *Верифицирован*"
					} else {
						isVerification = "⚠️ *Не верифицирован*"
					}
					if resGetUserProfile[0].IsPremium {
						isPremium = "✅ *Премиум*"
					} else {
						isPremium = "❌ *Не премиум*"
					}
					photo.Caption = fmt.Sprintf("*Личный кабинет*\n\n💰 Баланс: *%.2f $*\n📤 На выводе: *%.2f $*\n💵 Заработано: *%.2f $*\n\n👤 Верификация: %s\n🌐 Статус аккаунта: %s\n🆔 Ваш ID: [%d](tg://user?id=%d)",
						resGetUserProfile[0].Balance,
						resGetUserProfile[0].Conclusion,
						resGetUserTotalEarn[0].TotalEarn,
						isVerification,
						isPremium,
						teleId,
						teleId,
					)
					photo.ReplyMarkup = keyboard.GenKeyboardInlineForProfileMenu("📥 Пополнить", "📤 Вывести", "📝 Верификация", "🇺🇸 English language", "en")
				}

				if languageUser == "en" {
					var isVerification string
					var isPremium string
					if resGetUserProfile[0].Verification {
						isVerification = "✅ *Verified*"
					} else {
						isVerification = "⚠️ *Not verified*"
					}
					if resGetUserProfile[0].IsPremium {
						isPremium = "✅ *Premium*"
					} else {
						isPremium = "❌ *Not premium*"
					}
					photo.Caption = fmt.Sprintf("*Personal account*\n\n💰 Balance: *%.2f $*\n📤 Withdrawal: *%.2f $*\n💵 Earned: *%.2f $*\n\n👤 Verification: %s\n🌐 Status Account: %s\n🆔 Your ID: [%d](tg://user?id=%d)",
						resGetUserProfile[0].Balance,
						resGetUserProfile[0].Conclusion,
						resGetUserTotalEarn[0].TotalEarn,
						isVerification,
						isPremium,
						teleId,
						teleId,
					)
					photo.ReplyMarkup = keyboard.GenKeyboardInlineForProfileMenu("📥 Deposit", "📤 Withdraw", "📝 Verification", "🇷🇺 Русский язык", "ru")
				}

				_, err = bot.Send(photo)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
