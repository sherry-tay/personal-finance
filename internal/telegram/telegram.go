package telegram

import (
	"fmt"
	"net/http"
	"log"
	"math"
	"os"
	"strings"

	"personal.finance/internal/firebase"
	"personal.finance/internal/sgx"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/olekukonko/tablewriter"
)

var telegramBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
var webhookURL = os.Getenv("TELEGRAM_WEBHOOK_URL")
var authorizedUser = "AUTHORISED_USER"

// Initialize the Telegram bot
func Initialize() {
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Fatalf("Failed to get Telegram bot: %v", err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	if _, err = bot.SetWebhook(tgbotapi.NewWebhook(webhookURL + "/" + bot.Token)); err != nil {
		log.Fatalf("Failed to set webhook for Telegram bot: %v", err)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServe("0.0.0.0:80", nil)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.Chat.UserName != authorizedUser || update.Message.Chat.Type != "private" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, this is a private personal bot!")
			bot.Send(msg)
			continue
		}

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "portfolio":
				msg.ParseMode = "Markdown"
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
	var data []stockInfo
	totalPortfolio, totalInvested := 0.0, 0.0

	averagePrice := firebase.GetHoldings()

	for key, value := range averagePrice {
		current := sgx.GetCurrentPrice(key)
		diff, percentage, profit := getProfit(current, value.Price, value.Volume)
		
		s := stockInfo {
			Code: key,
			CurrentPrice: current,
			Average: value.Price,
			Diff: diff,
			DiffPercentage: percentage,
			Volume: value.Volume,
			Profit: profit,
		}
		data = append(data, s)

		totalPortfolio += current * float64(value.Volume)
		totalInvested += value.Price * float64(value.Volume)
	}

	totalDiff, totalPercentage, _ := getProfit(totalPortfolio, totalInvested, 0)
	
	return fmt.Sprintf(
		"```\n%v```\nCurrent portfolio: %.3f\nTotal invested: %.3f\nCapital gains: %.3f (%.3f%%)", 
		formatTable(data), totalPortfolio, totalInvested, totalDiff, totalPercentage,
	)
}

func getProfit(current, average float64, volume int) (float64, float64, float64) {
	diff := (current - average)
	return round(diff), round(diff/average*100), round(diff * float64(volume))
}

func round(val float64) float64 {
	precisionFactor := 1000.0
	return math.Round(val * precisionFactor)/precisionFactor
}

func formatTable(stockInfos []stockInfo) string {
	var formatted [][]string
	 for _, i := range stockInfos {
		strings := []string { 
			i.Code, 
			fmt.Sprintf("%.3f", i.CurrentPrice), 
			fmt.Sprintf("%.3f", i.Average), 
			fmt.Sprintf("%.3f", i.Diff), 
			fmt.Sprintf("%.3f", i.DiffPercentage), 
			fmt.Sprintf("%v", i.Volume), 
			fmt.Sprintf("%.3f", i.Profit),
		}
		formatted = append(formatted, strings)
	 }

	 tableString := &strings.Builder{}
	 table := tablewriter.NewWriter(tableString)
	 table.SetHeader([]string { "Code", "Now", "Avg", "Diff", "%", "Vol", "Total" })
	 table.SetBorders(tablewriter.Border {Left: true, Top: false, Right: true, Bottom: false})
	 table.SetCenterSeparator("|")
	 table.AppendBulk(formatted)
	 table.Render()
	 return tableString.String()
}

type stockInfo struct {
	Code string
	CurrentPrice float64
	Average float64
	Diff float64
	DiffPercentage float64
	Volume int
	Profit float64
}