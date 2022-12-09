package verification

import (
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rob-bender/invest-market-frontend/pkg/telegram/keyboard"
	requestProject "github.com/rob-bender/invest-market-frontend/pkg/telegram/request"
)

func Verification(bot *tgbotapi.BotAPI, msg tgbotapi.MessageConfig, teleId int64, userName string, languageUser string) error {
	if len(languageUser) > 0 {
		resGetUserProfile, err := requestProject.GetUserProfile(teleId)
		if err != nil {
			return err
		}
		if len(resGetUserProfile) > 0 {
			var isTesting string = os.Getenv("IS_TESTING")
			var needPath string = ""
			var varificationTextRu string = ""
			var varificationTextEn string = ""
			if isTesting == "true" {
				needPath = "/home/dale/job/work/my-project/nft-market/frontend/img"
			} else {
				needPath = "/home/nft-market-bot/frontend/invest-market-frontend/img"
			}
			if resGetUserProfile[0].Verification {
				varificationTextRu = "Ваш аккаунт верифицирован"
				varificationTextEn = "Your account is verified"
			} else {
				varificationTextRu = "Ваш аккаунт не верифицирован"
				varificationTextEn = "Your account is not verified"
			}
			photo := tgbotapi.NewPhoto(teleId, tgbotapi.FilePath(fmt.Sprintf("%s%s", needPath, "/img-need/4.jpg")))
			photo.ParseMode = "Markdown"
			if languageUser == "ru" {
				photo.Caption = fmt.Sprintf("*%s*\n\nДля прохождения верификации, обратитесь с службу поддержки. Оператор службы поддержки даст инструкции, необходимые для верификации счёта.", varificationTextRu)
				photo.ReplyMarkup = keyboard.GenKeyboardInlineForVerification("👨‍💻 Поддержка", "🔙 Вернуться в ЛК")
			}
			if languageUser == "en" {
				photo.Caption = fmt.Sprintf("*%s*\n\nTo pass verification, contact the support service. The support operator will give you the instructions necessary to verify your account.", varificationTextEn)
				photo.ReplyMarkup = keyboard.GenKeyboardInlineForVerification("👨‍💻 Support", "🔙 Back to profile")
			}
			_, err := bot.Send(photo)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
