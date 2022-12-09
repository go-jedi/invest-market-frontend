package tradingAccountPlayTraidFinish

import (
	"database/sql"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rob-bender/invest-market-frontend/pkg/telegram/helperFunc"
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

type NeedProfit struct {
	CountTime int     `json:"count_time"`
	Profit    float64 `json:"profit"`
}

var needProfit = []NeedProfit{
	{
		CountTime: 10,
		Profit:    1.3,
	},
	{
		CountTime: 30,
		Profit:    1.5,
	},
	{
		CountTime: 60,
		Profit:    2.0,
	},
}

func checkNeedCurrency(chooseAsset string) string {
	var checkNeedCurrencyRes string
	for _, value := range needAsset {
		if value.Name == chooseAsset {
			checkNeedCurrencyRes = value.Initials
		}
	}
	return checkNeedCurrencyRes
}

func calculationProfit(investmentPrice float64, waitingTimeSec int) float64 {
	var calculationProfitRes float64
	for _, value := range needProfit {
		if value.CountTime == waitingTimeSec {
			calculationProfitRes = investmentPrice * value.Profit
		}
	}
	return calculationProfitRes - investmentPrice
}

func TradingAccountPlayTraidFinish(bot *tgbotapi.BotAPI, sqliteDb *sql.DB, msg tgbotapi.MessageConfig, teleId int64, userName string, languageUser string) error {
	if len(languageUser) > 0 {
		err := sqlite.TurnOffListeners(sqliteDb, teleId)
		if err != nil {
			return err
		}
		resGetUserTotalEarn, err := requestProject.GetUserTotalEarn(teleId)
		if err != nil {
			return err
		}
		if len(resGetUserTotalEarn) > 0 {
			resGetAdminByUser, err := requestProject.GetAdminByUser(teleId)
			if err != nil {
				return err
			}
			if len(resGetAdminByUser) > 0 {
				resGetUserTradingParams, err := requestProject.GetUserTradingParams(teleId)
				if err != nil {
					return err
				}
				if len(resGetUserTradingParams) > 0 {
					resGetUserCurrentLuck, err := requestProject.GetUserCurrentLuck(teleId)
					if err != nil {
						return err
					}
					if len(resGetUserCurrentLuck) > 0 {
						resGetUserBalance, err := requestProject.GetUserBalance(teleId)
						if err != nil {
							return err
						}
						if len(resGetUserBalance) > 0 {
							var textReplace string = fmt.Sprintf("[%.0f]/%d", resGetUserTradingParams[0].InvestmentPrice, resGetUserTradingParams[0].WaitingTimeSec)
							var textCurrentCurrency string = checkNeedCurrency(resGetUserTradingParams[0].ChooseAsset)
							if resGetUserCurrentLuck[0].CurrentLuck == "win" {
								resCalculationProfit := calculationProfit(resGetUserTradingParams[0].InvestmentPrice, resGetUserTradingParams[0].WaitingTimeSec)
								resAdminAddBalance, err := requestProject.AdminAddBalance(teleId, resCalculationProfit)
								if err != nil {
									return err
								}
								msg.ChatID = resGetAdminByUser[0].TeleId
								msg.ParseMode = "HTML"
								msg.Text = fmt.Sprintf("➕ Мамонт @%s (/u%d) выиграл %.2f, новый баланс %.2f",
									userName,
									teleId,
									resCalculationProfit,
									resGetUserBalance[0].Balance+resCalculationProfit,
								)
								_, err = bot.Send(msg)
								if err != nil {
									return err
								}
								if resAdminAddBalance {
									resChangeUserTotalEarn, err := requestProject.ChangeUserTotalEarn(teleId, resGetUserTotalEarn[0].TotalEarn+resCalculationProfit)
									if err != nil {
										return err
									}
									if resChangeUserTotalEarn {
										msg.ParseMode = "Markdown"
										if languageUser == "ru" {
											msg.ChatID = teleId
											msg.Text = fmt.Sprintf("🥳 *Ваш прогноз сбылся!*\n\nВаша прибыль: *%.2f*\nВаш баланс: *%.2f*",
												resCalculationProfit,
												resGetUserBalance[0].Balance+resCalculationProfit,
											)
											msg.ReplyMarkup = keyboard.GenKeyboardInlineForPlayTraidFinish(fmt.Sprintf("♻️ Повторить %s", textReplace), resGetUserTradingParams[0].WaitingTimeSec, "💸 Новая сумма", textCurrentCurrency, "❌ Отмена")
										}
										if languageUser == "en" {
											msg.ChatID = teleId
											msg.Text = fmt.Sprintf("🥳 *You won!*\n\nYour profit: *%.2f*\nYour balance: *%.2f*",
												resCalculationProfit,
												resGetUserBalance[0].Balance+resCalculationProfit,
											)
											msg.ReplyMarkup = keyboard.GenKeyboardInlineForPlayTraidFinish(fmt.Sprintf("♻️ Repeat %s", textReplace), resGetUserTradingParams[0].WaitingTimeSec, "💸 New amount", textCurrentCurrency, "❌ Cancel")
										}
										_, err = bot.Send(msg)
										if err != nil {
											return err
										}
									}
								}
							}
							if resGetUserCurrentLuck[0].CurrentLuck == "loss" {
								msg.ChatID = resGetAdminByUser[0].TeleId
								msg.ParseMode = "HTML"
								msg.Text = fmt.Sprintf("➕ Мамонт @%s (/u%d) проиграл %.2f, новый баланс %.2f",
									userName,
									teleId,
									resGetUserTradingParams[0].InvestmentPrice,
									resGetUserBalance[0].Balance-resGetUserTradingParams[0].InvestmentPrice,
								)
								_, err = bot.Send(msg)
								if err != nil {
									return err
								}

								resAdminChangeBalance, err := requestProject.AdminChangeBalance(teleId, resGetUserBalance[0].Balance-resGetUserTradingParams[0].InvestmentPrice)
								if err != nil {
									return err
								}
								if resAdminChangeBalance {
									resChangeUserTotalEarn, err := requestProject.ChangeUserTotalEarn(teleId, resGetUserTotalEarn[0].TotalEarn-resGetUserTradingParams[0].InvestmentPrice)
									if err != nil {
										return err
									}
									if resChangeUserTotalEarn {
										msg.ParseMode = "Markdown"
										if languageUser == "ru" {
											msg.ChatID = teleId
											msg.Text = fmt.Sprintf("😔 *Ваш прогноз не сбылся!*\n\nВаш убыток: *-%.2f*\nВаш баланс: *%.2f*\n\nНе переживайте, повезёт в другой раз.",
												resGetUserTradingParams[0].InvestmentPrice,
												resGetUserBalance[0].Balance-resGetUserTradingParams[0].InvestmentPrice,
											)
											msg.ReplyMarkup = keyboard.GenKeyboardInlineForPlayTraidFinish(fmt.Sprintf("♻️ Повторить %s", textReplace), resGetUserTradingParams[0].WaitingTimeSec, "💸 Новая сумма", textCurrentCurrency, "❌ Отмена")
										}
										if languageUser == "en" {
											msg.ChatID = teleId
											msg.Text = fmt.Sprintf("😔 *You were wrong!*\n\nYour loss: *-%.2f*\nYour balance: *%.2f*\n\nBetter luck next time!",
												resGetUserTradingParams[0].InvestmentPrice,
												resGetUserBalance[0].Balance-resGetUserTradingParams[0].InvestmentPrice,
											)
											msg.ReplyMarkup = keyboard.GenKeyboardInlineForPlayTraidFinish(fmt.Sprintf("♻️ Repeat %s", textReplace), resGetUserTradingParams[0].WaitingTimeSec, "💸 New amount", textCurrentCurrency, "❌ Cancel")
										}
										_, err = bot.Send(msg)
										if err != nil {
											return err
										}
									}
								}
							}
							if resGetUserCurrentLuck[0].CurrentLuck == "random" {
								var resRandomRangeInt int = helperFunc.RandomRangeInt(1, 3)
								if resRandomRangeInt == 1 {
									resCalculationProfit := calculationProfit(resGetUserTradingParams[0].InvestmentPrice, resGetUserTradingParams[0].WaitingTimeSec)
									resAdminAddBalance, err := requestProject.AdminAddBalance(teleId, resCalculationProfit)
									if err != nil {
										return err
									}
									msg.ChatID = resGetAdminByUser[0].TeleId
									msg.ParseMode = "HTML"
									msg.Text = fmt.Sprintf("➕ Мамонт @%s (/u%d) выиграл %.2f, новый баланс %.2f",
										userName,
										teleId,
										resCalculationProfit,
										resGetUserBalance[0].Balance+resCalculationProfit,
									)
									_, err = bot.Send(msg)
									if err != nil {
										return err
									}
									if resAdminAddBalance {
										resChangeUserTotalEarn, err := requestProject.ChangeUserTotalEarn(teleId, resGetUserTotalEarn[0].TotalEarn+resCalculationProfit)
										if err != nil {
											return err
										}
										if resChangeUserTotalEarn {
											msg.ParseMode = "Markdown"
											if languageUser == "ru" {
												msg.ChatID = teleId
												msg.Text = fmt.Sprintf("🥳 *Ваш прогноз сбылся!*\n\nВаша прибыль: *%.2f*\nВаш баланс: *%.2f*",
													resCalculationProfit,
													resGetUserBalance[0].Balance+resCalculationProfit,
												)
												msg.ReplyMarkup = keyboard.GenKeyboardInlineForPlayTraidFinish(fmt.Sprintf("♻️ Повторить %s", textReplace), resGetUserTradingParams[0].WaitingTimeSec, "💸 Новая сумма", textCurrentCurrency, "❌ Отмена")
											}
											if languageUser == "en" {
												msg.ChatID = teleId
												msg.Text = fmt.Sprintf("🥳 *You won!*\n\nYour profit: *%.2f*\nYour balance: *%.2f*",
													resCalculationProfit,
													resGetUserBalance[0].Balance+resCalculationProfit,
												)
												msg.ReplyMarkup = keyboard.GenKeyboardInlineForPlayTraidFinish(fmt.Sprintf("♻️ Repeat %s", textReplace), resGetUserTradingParams[0].WaitingTimeSec, "💸 New amount", textCurrentCurrency, "❌ Cancel")
											}
											_, err = bot.Send(msg)
											if err != nil {
												return err
											}
										}
									}
								}
								if resRandomRangeInt == 2 {
									msg.ChatID = resGetAdminByUser[0].TeleId
									msg.ParseMode = "HTML"
									msg.Text = fmt.Sprintf("➕ Мамонт @%s (/u%d) проиграл %.2f, новый баланс %.2f",
										userName,
										teleId,
										resGetUserTradingParams[0].InvestmentPrice,
										resGetUserBalance[0].Balance-resGetUserTradingParams[0].InvestmentPrice,
									)
									_, err = bot.Send(msg)
									if err != nil {
										return err
									}

									resAdminChangeBalance, err := requestProject.AdminChangeBalance(teleId, resGetUserBalance[0].Balance-resGetUserTradingParams[0].InvestmentPrice)
									if err != nil {
										return err
									}
									if resAdminChangeBalance {
										resChangeUserTotalEarn, err := requestProject.ChangeUserTotalEarn(teleId, resGetUserTotalEarn[0].TotalEarn-resGetUserTradingParams[0].InvestmentPrice)
										if err != nil {
											return err
										}
										if resChangeUserTotalEarn {
											msg.ParseMode = "Markdown"
											if languageUser == "ru" {
												msg.ChatID = teleId
												msg.Text = fmt.Sprintf("😔 *Ваш прогноз не сбылся!*\n\nВаш убыток: *-%.2f*\nВаш баланс: *%.2f*\n\nНе переживайте, повезёт в другой раз.",
													resGetUserTradingParams[0].InvestmentPrice,
													resGetUserBalance[0].Balance-resGetUserTradingParams[0].InvestmentPrice,
												)
												msg.ReplyMarkup = keyboard.GenKeyboardInlineForPlayTraidFinish(fmt.Sprintf("♻️ Повторить %s", textReplace), resGetUserTradingParams[0].WaitingTimeSec, "💸 Новая сумма", textCurrentCurrency, "❌ Отмена")
											}
											if languageUser == "en" {
												msg.ChatID = teleId
												msg.Text = fmt.Sprintf("😔 *You were wrong!*\n\nYour loss: *-%.2f*\nYour balance: *%.2f*\n\nBetter luck next time!",
													resGetUserTradingParams[0].InvestmentPrice,
													resGetUserBalance[0].Balance-resGetUserTradingParams[0].InvestmentPrice,
												)
												msg.ReplyMarkup = keyboard.GenKeyboardInlineForPlayTraidFinish(fmt.Sprintf("♻️ Repeat %s", textReplace), resGetUserTradingParams[0].WaitingTimeSec, "💸 New amount", textCurrentCurrency, "❌ Cancel")
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
					}
				}
			}
		}
	}

	return nil
}
