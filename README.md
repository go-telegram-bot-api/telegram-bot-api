# Golang Telegram bindings for the Bot API

[![GoDoc](https://godoc.org/src.foxpaw.in/Syfaro/telegram-bot-api?status.svg)](https://godoc.org/src.foxpaw.in/Syfaro/telegram-bot-api)

All methods have been added, and all features should be available.
If you want a feature that hasn't been added yet or something is broken, open an issue and I'll see what I can do.

## Example

This is a very simple bot that just displays any gotten updates, then replies it to that chat.

```go
package main

import (
	"log"
	"src.foxpaw.in/Syfaro/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("MyAwesomeBotToken")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.UpdatesChan(u)

	for update := range updates {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.SendMessage(msg)
	}
}
```
