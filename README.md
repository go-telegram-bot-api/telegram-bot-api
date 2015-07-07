# Golang Telegram bindings for the Bot API

[![GoDoc](https://godoc.org/github.com/Syfaro/telegram-bot-api?status.svg)](http://godoc.org/github.com/Syfaro/telegram-bot-api)

This was forked from [Syfaro](https://github.com/Syfaro/telegram-bot-api).  I consolidate a few of the sending methods and reorganized the code a little bit.  Ultimately, the credit should go to Syfaro.

## Example

This is a very simple bot that just displays any gotten updates, then replies it to that chat.

```go
package main

import (
	"log"
	"github.com/Syfaro/telegram-bot-api"
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
