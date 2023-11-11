# Suggestions with custom commands

Suggestions with personalized commands, so that when the user types "/", they will receive suggestions for possible commands. Also, add a menu to view all the available commands.

```go
package main

import (
	"bufio"
	"context"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	bot *tgbotapi.BotAPI
)

func main() {
	var err error
	bot, err = tgbotapi.NewBotAPI("<TELEGRAM_APITOKEN>")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	updates := bot.GetUpdatesChan(u)

	commands := []tgbotapi.BotCommand{
		{
			Command:     "repo",
			Description: "repository",
		},
		{
			Command:     "issues",
			Description: "contribute reviews to issues",
		},
		{
			Command:     "pr",
			Description: "pull requests",
		},
		{
			Command:     "doc",
			Description: "documentation",
		},
	}

	// Defines the ReplyKeyboardMarkup object with custom commands as suggestions
	replyKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/repo"),
			tgbotapi.NewKeyboardButton("/issues"),
			tgbotapi.NewKeyboardButton("/pr"),
			tgbotapi.NewKeyboardButton("/doc"),
		),
	)
	replyKeyboard.ResizeKeyboard = true

	config := tgbotapi.NewSetMyCommands(commands...)
	bot.Send(config)

	go receiveUpdates(ctx, updates)

	log.Println("INFO: Start listening for updates. Press enter to stop")

	_, err = bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
}

func receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	<-updates

	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			handleUpdate(update)
		}
	}
}

func Commands(message *tgbotapi.Message, command string) error {
	chatID := message.Chat.ID

	switch message.Command() {
	case "repo":
		return SendMessage(chatID, "Repository: https://github.com/go-telegram-bot-api/telegram-bot-api")

	case "issues":
		return SendMessage(chatID, "Issues: https://github.com/go-telegram-bot-api/telegram-bot-api/issues")

	case "pr":
		return SendMessage(chatID, "Pull requests: https://github.com/go-telegram-bot-api/telegram-bot-api/pulls")

	case "doc":
		return SendMessage(chatID, "Documentation: https://go-telegram-bot-api.dev/")
	default:
		return SendMessage(chatID, "Not Found Command")
	}
}

func handleUpdate(update tgbotapi.Update) {
	switch {
	case update.Message != nil:
		handleMessage(update.Message)
	}
}

func handleMessage(message *tgbotapi.Message) {
	user := message.From
	text := message.Text

	if user == nil {
		return
	}

	if message.IsCommand() {
		err := Commands(message, text)
		if err != nil {
			log.Printf("ERROR: %s", err.Error())
			return
		}
	}
}

func SendMessage(chatID int64, msg string) error {
	_, err := bot.Send(tgbotapi.NewMessage(chatID, msg))
	return err
}
```
