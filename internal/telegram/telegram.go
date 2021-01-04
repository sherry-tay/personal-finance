package telegram

import (
	"log"

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
			case "help":
				msg.Text = "type /sayhi or /status."
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