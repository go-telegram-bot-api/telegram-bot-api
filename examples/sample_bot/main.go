package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const BotToken string = "YOUR-BOT-TOKEN"

func main() {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		panic(err)
	}
	updates := bot.GetUpdatesChan(tgbotapi.UpdateConfig{})
	for update := range updates {
		if update.Message != nil {
			switch update.Message.Text {
			case "/start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				if update.Message.From.FirstName != "" {
					msg.Text = fmt.Sprintf("Hello, *%s*!", update.Message.From.FirstName)
				} else {
					msg.Text = "Hello!"
				}
				msg.ParseMode = tgbotapi.ModeMarkdown
				bot.Send(msg)
			}
		}
	}
}
