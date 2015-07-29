# Golang bindings for the Telegram Bot API

[![GoDoc](https://godoc.org/github.com/Syfaro/telegram-bot-api?status.svg)](http://godoc.org/github.com/Syfaro/telegram-bot-api)

All methods have been added, and all features should be available.
If you want a feature that hasn't been added yet or something is broken, open an issue and I'll see what I can do.

All methods are fairly self explanatory, and reading the godoc page should explain everything. If something isn't clear, open an issue or submit a pull request.

The scope of this project is just to provide a wrapper around the API without any additional features. There are other projects for creating something with plugins and command handlers without having to design all that yourself.

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
