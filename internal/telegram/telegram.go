package telegram

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/sherry-tay/personal-finance/internal/firestore"
	"github.com/sherry-tay/personal-finance/internal/sgx"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/olekukonko/tablewriter"
)

const dateInputFormat = "20060102" // 2006-Jan-02

var telegramBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
var webhookURL = os.Getenv("TELEGRAM_WEBHOOK_URL")
var authorizedUser = os.Getenv("AUTHORIZED_USER")
var port = os.Getenv("PORT")

// Initialize the Telegram bot
func Initialize() {
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Fatalf("Failed to get Telegram bot: %v", err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	if (webhookURL != "") {
		if _, err = bot.SetWebhook(tgbotapi.NewWebhook(webhookURL + "/" + bot.Token)); err != nil {
			log.Fatalf("Failed to set webhook for Telegram bot: %v", err)
		}
	}

	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServe("0.0.0.0:" + port, nil)

	fs := newCustomFirestore()

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
				msg.Text = fs.getStatistics(newCustomSgx(), formatTable)
			case "add":
				msg.Text = fs.addHoldings(update.Message.CommandArguments())
			case "detailed":
				msg.ParseMode = "Markdown"
				msg.Text = fs.getStatistics(newCustomSgx(), formatDetailedTable)
			case "withArgument":
				msg.Text = "You supplied the following argument: " + update.Message.CommandArguments()
			case "html":
				msg.ParseMode = "html"
				msg.Text = "This will be interpreted as HTML, click <a href=\"https://www.example.com\">here</a>"
			default:
				msg.Text = `Use one of the commands available. To add, follow the format: "/add <code> <price> <volume> <yyyymmdd> <stored location>" e.g. /add ABC 1.23 100 20210101 mybroker`
			}
			bot.Send(msg)
		}
	}
}

func (fs *customFirestore) getStatistics(sgx *customSgx, formatter func(stockInfos []stockInfo) string) string {
	var data []stockInfo
	totalPortfolio, totalInvested := 0.0, 0.0

	averagePrice, err := fs.read()
	if err != nil {
		return "Something went wrong while retrieving your holdings"
	}

	sortedCodes := make([]string, 0, len(averagePrice))
	for key := range averagePrice {
		sortedCodes = append(sortedCodes, key)
	}
	sort.Strings(sortedCodes)

	for _, key := range sortedCodes {
		value := averagePrice[key]
		current, err := sgx.get(key)
		if err != nil {
			break
		}
		diff, percentage, profit := getProfit(current, value.Price, value.Volume)

		s := stockInfo{
			Code:           key,
			CurrentPrice:   current,
			Average:        value.Price,
			Diff:           diff,
			DiffPercentage: percentage,
			Volume:         value.Volume,
			Profit:         profit,
		}
		data = append(data, s)

		totalPortfolio += current * float64(value.Volume)
		totalInvested += value.Price * float64(value.Volume)
	}

	totalDiff, totalPercentage, _ := getProfit(totalPortfolio, totalInvested, 0)

	return fmt.Sprintf(
		"```\n%v```\nCurrent portfolio: %.3f\nTotal invested: %.3f\nCapital gains: %.3f (%.3f%%)",
		formatter(data), totalPortfolio, totalInvested, totalDiff, totalPercentage,
	)
}

func getProfit(current, average float64, volume int) (float64, float64, float64) {
	diff := (current - average)
	return round(diff), round(diff / average * 100), round(diff * float64(volume))
}

func round(val float64) float64 {
	precisionFactor := 1000.0
	return math.Round(val*precisionFactor) / precisionFactor
}

func formatTable(stockInfos []stockInfo) string {
	var formatted [][]string
	for _, i := range stockInfos {
		strings := []string{
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
	table.SetHeader([]string{"Code", "Now", "Avg", "Diff", "%", "Vol", "Total"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(formatted)
	table.Render()
	return tableString.String()
}

func formatDetailedTable(stockInfos []stockInfo) string {
	var formatted [][]string
	for _, i := range stockInfos {
		formatted = append(formatted, []string{i.Code, "Vol", "1", fmt.Sprintf("%v", i.Volume)})
		formatted = append(formatted, []string{i.Code, "Now", fmt.Sprintf("%.3f", i.CurrentPrice), fmt.Sprintf("%.3f", i.CurrentPrice*float64(i.Volume))})
		formatted = append(formatted, []string{i.Code, "Avg", fmt.Sprintf("%.3f", i.Average), fmt.Sprintf("%.3f", i.Average*float64(i.Volume))})
		formatted = append(formatted, []string{i.Code, "Dif", fmt.Sprintf("%.3f", i.Diff), fmt.Sprintf("%.3f", i.Profit)})
		formatted = append(formatted, []string{i.Code, "%", "", fmt.Sprintf("%.3f", i.DiffPercentage)})
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"", "", "Per", "Total"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.AppendBulk(formatted)
	table.Render()
	return tableString.String()
}

func (fs *customFirestore) addHoldings(arg string) string {
	params := strings.SplitN(arg, " ", 5)
	if len(params) != 5 {
		return "Not enough arguments supplied! Aborting..."
	}

	code := params[0]
	price, priceErr := strconv.ParseFloat(params[1], 64)
	volume, volumeErr := strconv.Atoi(params[2])
	date, dateErr := time.Parse(dateInputFormat, params[3])
	in := params[4]

	if priceErr != nil || volumeErr != nil || dateErr != nil {
		return "Something went wrong while parsing the arguments... Please enter using the correct format."
	}

	id := params[3] + "-" + code

	s := firestore.Stock{
		Code:     code,
		Price:    price,
		Volume:   volume,
		Date:     date,
		StoredIn: in,
	}

	if err := fs.add(id, s); err != nil {
		return "Something went wrong while adding your holdings"
	}
	return "Successfully added holdings!"
}

type stockInfo struct {
	Code           string
	CurrentPrice   float64
	Average        float64
	Diff           float64
	DiffPercentage float64
	Volume         int
	Profit         float64
}

type readFunction func() (map[string]firestore.Stock, error)
type addFunction func(id string, s firestore.Stock) error

type customFirestore struct {
	read readFunction
	add  addFunction
}

func newCustomFirestore() *customFirestore {
	return &customFirestore{
		read: firestore.GetHoldings,
		add:  firestore.AddHoldings,
	}
}

type getCurrentPrice func(code string) (float64, error)

type customSgx struct {
	get getCurrentPrice
}

func newCustomSgx() *customSgx {
	return &customSgx{get: sgx.GetCurrentPrice}
}
