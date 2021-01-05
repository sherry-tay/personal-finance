package telegram

import (
	"fmt"
	"log"
	"math"

	"personal.finance/internal/firebase"
	"personal.finance/internal/sgx"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func Initialize() {
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Fatalf("Failed to get Telegram bot: %v", err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.Chat.UserName != authorizedUser || update.Message.Chat.Type != "private" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, this is a private personal bot!")
			bot.Send(msg)
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "portfolio":
				msg.Text = getStatistics()
			case "sayhi":
				msg.Text = "Hi :)"
			case "status":
				msg.Text = "I'm ok."
			case "withArgument":
				msg.Text = "You supplied the following argument: " + update.Message.CommandArguments()
			case "html":
				msg.ParseMode = "html"
				msg.Text = "This will be interpreted as HTML, click <a href=\"https://www.example.com\">here</a>"
			default:
				msg.Text = "I don't know that command"
			}
			bot.Send(msg)
		}
	}
}

func getStatistics() string {
	message := ""
	totalPortfolio, totalInvested := 0.0, 0.0
	averagePrice := firebase.GetHoldings()
	for key, value := range averagePrice {
		current := sgx.GetCurrentPrice(key)
		diff, percentage, profit := getProfit(current, value.Price, value.Volume)
		message += fmt.Sprintf("%v: Current %v, Average %v, Difference %v (%v%%), Volume %v, Profit %v\n", key, current, value.Price, diff, percentage, value.Volume, profit)
		totalPortfolio += current * float64(value.Volume)
		totalInvested += value.Price * float64(value.Volume)
	}
	totalDiff, totalPercentage, _ := getProfit(totalPortfolio, totalInvested, 0)
	message += fmt.Sprintf("Current portfolio: %v\n", round(totalPortfolio))
	message += fmt.Sprintf("Total invested: %v\n", round(totalInvested))
	message += fmt.Sprintf("Capital gains: %v (%v%%)\n", totalDiff, totalPercentage)
	return message
}

func getProfit(current, average float64, volume int) (float64, float64, float64) {
	diff := (current - average)
	return round(diff), round(diff/average*100), round(diff * float64(volume))
}

func round(val float64) float64 {
	precisionFactor := 1000.0
	return math.Round(val * precisionFactor)/precisionFactor
}