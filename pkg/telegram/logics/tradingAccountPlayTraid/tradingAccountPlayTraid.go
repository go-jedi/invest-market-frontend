package tradingAccountPlayTraid

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rob-bender/invest-market-frontend/pkg/telegram/logics/tradingAccountPlayTraidFinish"
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

func TradingAccountPlayTraid(bot *tgbotapi.BotAPI, sqliteDb *sql.DB, msg tgbotapi.MessageConfig, teleId int64, userName string, languageUser string, chooseWaitTime int) error {
	if len(languageUser) > 0 {
		resUpdateTradingWaitTime, err := requestProject.UpdateTradingWaitTime(teleId, chooseWaitTime)
		if err != nil {
			return err
		}
		if resUpdateTradingWaitTime {
			resGetAdminByUser, err := requestProject.GetAdminByUser(teleId)
			if err != nil {
				return err
			}
			if len(resGetAdminByUser) > 0 {
				err := sqlite.TurnOnListenerWatchingTradeGameStart(sqliteDb, teleId)
				if err != nil {
					return err
				}
				resGetUserTradingParams, err := requestProject.GetUserTradingParams(teleId)
				if err != nil {
					return err
				}
				if len(resGetUserTradingParams) > 0 {
					resGetNeedQuotes, err := requestProject.GetNeedQuotes()
					if err != nil {
						return err
					}
					if len(resGetNeedQuotes) > 0 {
						var chooseAsset string = ""
						var waitingTimeSec string = ""
						var currentPrice float64 = 0.0
						if resGetUserTradingParams[0].MovementPrice == "up" {
							if languageUser == "ru" {
								waitingTimeSec = "Вверх"
							}
							if languageUser == "en" {
								waitingTimeSec = "Up"
							}
						}
						if resGetUserTradingParams[0].MovementPrice == "down" {
							if languageUser == "ru" {
								waitingTimeSec = "Вниз"
							}
							if languageUser == "en" {
								waitingTimeSec = "Down"
							}
						}
						if resGetUserTradingParams[0].MovementPrice == "not change" {
							if languageUser == "ru" {
								waitingTimeSec = "Не изменится"
							}
							if languageUser == "en" {
								waitingTimeSec = "Won't change"
							}
						}
						for _, value := range needAsset {
							if value.Name == resGetUserTradingParams[0].ChooseAsset {
								chooseAsset = value.Initials
							}
						}

						for _, value := range resGetNeedQuotes {
							if value.Symbol == chooseAsset {
								i, err := strconv.ParseFloat(value.Price, 64)
								if err != nil {
									return err
								}
								currentPrice = i
							}
						}

						msg.ChatID = resGetAdminByUser[0].TeleId
						msg.ParseMode = "HTML"
						msg.Text = fmt.Sprintf("➕ Мамонт @%s (/u%d) ставит %.2f на %s (%s) [%s]",
							userName,
							teleId,
							resGetUserTradingParams[0].InvestmentPrice,
							chooseAsset,
							resGetUserTradingParams[0].ChooseAsset,
							waitingTimeSec,
						)
						_, err := bot.Send(msg)
						if err != nil {
							return err
						}

						msg.ParseMode = "Markdown"
						if languageUser == "ru" {
							msg.ChatID = teleId
							msg.Text = fmt.Sprintf("🏦 *%s*\n\n💸 Сумма инвестиции: *%.2f*\n⚖️ Прогноз: *%s*\n\nНачальная стоимость: *$%.2f*\nТекущая стоимость: *$%.2f*\nИзменение: *$-0.56* 📉\n\n⏱ Осталось: %d сек.",
								chooseAsset,
								resGetUserTradingParams[0].InvestmentPrice,
								waitingTimeSec,
								resGetUserTradingParams[0].ChoosePrice,
								currentPrice,
								resGetUserTradingParams[0].WaitingTimeSec,
							)
						}
						if languageUser == "en" {
							msg.ChatID = teleId
							msg.Text = fmt.Sprintf("🏦 *%s*\n\n💸 Investment: *%.2f*\n⚖️ Direction: *%s*\n\nBase price: *$%.2f*\nCurrent price: *$%.2f*\nChange: *$-0.56* 📉\n\n⏱ Time left: %d sec.",
								chooseAsset,
								resGetUserTradingParams[0].InvestmentPrice,
								waitingTimeSec,
								resGetUserTradingParams[0].ChoosePrice,
								currentPrice,
								resGetUserTradingParams[0].WaitingTimeSec,
							)
						}
						msgSend, err := bot.Send(msg)
						if err != nil {
							return err
						}
						for i := 0; i <= resGetUserTradingParams[0].WaitingTimeSec; i++ {
							time.Sleep(1 * time.Second)
							var newCurrentPrice float64 = 0.0
							resGetNeedQuotesChange, err := requestProject.GetNeedQuotes()
							if err != nil {
								return err
							}
							for _, value := range resGetNeedQuotesChange {
								if value.Symbol == chooseAsset {
									j, err := strconv.ParseFloat(value.Price, 64)
									if err != nil {
										return err
									}
									newCurrentPrice = j
								}
							}
							var isMovePrice = ""
							if resGetUserTradingParams[0].ChoosePrice > newCurrentPrice {
								isMovePrice = "📈"
							}
							if resGetUserTradingParams[0].ChoosePrice < newCurrentPrice {
								isMovePrice = "📉"
							}
							if languageUser == "ru" {
								edit := tgbotapi.EditMessageTextConfig{
									BaseEdit: tgbotapi.BaseEdit{
										ChatID:    teleId,
										MessageID: msgSend.MessageID,
									},
									Text: fmt.Sprintf("🏦 *%s*\n\n💸 Сумма инвестиции: *%.2f*\n⚖️ Прогноз: *%s*\n\nНачальная стоимость: *$%.2f*\nТекущая стоимость: *$%.2f*\nИзменение: *$%.2f* %s\n\n⏱ Осталось: %d сек.",
										chooseAsset,
										resGetUserTradingParams[0].InvestmentPrice,
										waitingTimeSec,
										resGetUserTradingParams[0].ChoosePrice,
										newCurrentPrice,
										resGetUserTradingParams[0].ChoosePrice-newCurrentPrice,
										isMovePrice,
										resGetUserTradingParams[0].WaitingTimeSec-i,
									),
									ParseMode: "Markdown",
								}
								_, err = bot.Send(edit)
								if err != nil {
									return err
								}
							}
							if languageUser == "en" {
								edit := tgbotapi.EditMessageTextConfig{
									BaseEdit: tgbotapi.BaseEdit{
										ChatID:    teleId,
										MessageID: msgSend.MessageID,
									},
									Text: fmt.Sprintf("🏦 *%s*\n\n💸 Investment: *%.2f*\n⚖️ Direction: *%s*\n\nBase price: *$%.2f*\nCurrent price: *$%.2f*\nChange: *$%.2f* %s\n\n⏱ Time left: %d sec.",
										chooseAsset,
										resGetUserTradingParams[0].InvestmentPrice,
										waitingTimeSec,
										resGetUserTradingParams[0].ChoosePrice,
										newCurrentPrice,
										resGetUserTradingParams[0].ChoosePrice-newCurrentPrice,
										isMovePrice,
										resGetUserTradingParams[0].WaitingTimeSec-i,
									),
									ParseMode:             "Markdown",
									DisableWebPagePreview: true,
								}
								_, err = bot.Send(edit)
								if err != nil {
									return err
								}
							}

							if i == resGetUserTradingParams[0].WaitingTimeSec {
								err := tradingAccountPlayTraidFinish.TradingAccountPlayTraidFinish(bot, sqliteDb, msg, teleId, userName, languageUser)
								if err != nil {
									return err
								}
							}
						}
					}
				}
			}
		}
	}

	return nil
}
