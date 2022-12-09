package tradingAccountWriteSum

import (
	"database/sql"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rob-bender/invest-market-frontend/pkg/telegram/keyboard"
	requestProject "github.com/rob-bender/invest-market-frontend/pkg/telegram/request"
	"github.com/rob-bender/invest-market-frontend/pkg/telegram/sqlite"
)

type NeedAsset struct {
	Name     string `json:"name"`
	Initials string `json:"initials"`
}

var needAsset = []NeedAsset{
	{
		Name:     "Bitcoin",
		Initials: "BTCUSDT",
	},
	{
		Name:     "Ethereum",
		Initials: "ETHUSDT",
	},
	{
		Name:     "BNB",
		Initials: "BNBUSDT",
	},
	{
		Name:     "Cardano",
		Initials: "ADAUSDT",
	},
	{
		Name:     "Solana",
		Initials: "SOLUSDT",
	},
	{
		Name:     "Dogecoin",
		Initials: "DOGEUSDT",
	},
	{
		Name:     "Polkadot",
		Initials: "DOTUSDT",
	},
	{
		Name:     "Polygon",
		Initials: "MATICUSDT",
	},
	{
		Name:     "TRON",
		Initials: "TRXUSDT",
	},
	{
		Name:     "ETH Classic",
		Initials: "ETCUSDT",
	},
	{
		Name:     "Litecoin",
		Initials: "LTCUSDT",
	},
	{
		Name:     "Monero",
		Initials: "XMRUSDT",
	},
	{
		Name:     "Bitcoin Cash",
		Initials: "BCHUSDT",
	},
	{
		Name:     "XRP",
		Initials: "XRPUSDT",
	},
}

func TradingAccountWriteSum(bot *tgbotapi.BotAPI, sqliteDb *sql.DB, msg tgbotapi.MessageConfig, teleId int64, userName string, languageUser string, userChooseAsset string) error {
	if len(languageUser) > 0 {
		err := sqlite.TurnOnListenerWatchingTradAcWriteSum(sqliteDb, teleId)
		if err != nil {
			return err
		}
		resGetNeedQuotes, err := requestProject.GetNeedQuotes()
		if err != nil {
			return err
		}
		if len(resGetNeedQuotes) > 0 {
			resGetUserBalance, err := requestProject.GetUserBalance(teleId)
			if err != nil {
				return err
			}
			if len(resGetUserBalance) > 0 {
				resGetUserMinPrice, err := requestProject.GetUserMinPrice(teleId)
				if err != nil {
					return err
				}
				if len(resGetUserMinPrice) > 0 {
					var chooseAsset string = ""
					var priceAsset float64 = 0.0
					for _, value := range needAsset {
						if value.Initials == userChooseAsset {
							chooseAsset = value.Name
						}
					}
					for _, value := range resGetNeedQuotes {
						if value.Symbol == userChooseAsset {
							i, err := strconv.ParseFloat(value.Price, 64)
							if err != nil {
								return err
							}
							priceAsset = i
						}
					}
					resUpdateTradingAsset, err := requestProject.UpdateTradingAsset(teleId, chooseAsset, priceAsset)
					if err != nil {
						return err
					}
					if resUpdateTradingAsset {
						if languageUser == "ru" {
							msg.Text = fmt.Sprintf("*🌐 Введите сумму, которую Вы хотите инвестировать.*\n\nМинимальная сумма инвестиций: *%.2f $*\nАктив: *%s*\n\nЦена: *$%.2f*\nВаш баланс: *%.2f $*",
								resGetUserMinPrice[0].MinimPrice,
								chooseAsset,
								priceAsset,
								resGetUserBalance[0].Balance,
							)
							msg.ReplyMarkup = keyboard.GenKeyboardInlineForTradingAcWriteSum("🔄 Обновить цену", userChooseAsset, "🔙 Назад")
						}
						if languageUser == "en" {
							msg.Text = fmt.Sprintf("*🌐 Enter the amount of investment.*\n\nMinimal investment: *%.2f $*\nAsset: *%s*\nPrice: *%.2f*\n\nYour balance: *%.2f $*",
								resGetUserMinPrice[0].MinimPrice,
								chooseAsset,
								priceAsset,
								resGetUserBalance[0].Balance,
							)
							msg.ReplyMarkup = keyboard.GenKeyboardInlineForTradingAcWriteSum("🔄 Update price", userChooseAsset, "🔙 Back")
						}
						_, err = bot.Send(msg)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}
