# Golang bindings for the Telegram Bot API

[![GoDoc](https://godoc.org/github.com/go-telegram-bot-api/telegram-bot-api?status.svg)](http://godoc.org/github.com/go-telegram-bot-api/telegram-bot-api)
[![Travis](https://travis-ci.org/go-telegram-bot-api/telegram-bot-api.svg)](https://travis-ci.org/go-telegram-bot-api/telegram-bot-api)

All methods are fairly self explanatory, and reading the godoc page should
explain everything. If something isn't clear, open an issue or submit
a pull request.

The scope of this project is just to provide a wrapper around the API
without any additional features. There are other projects for creating
something with plugins and command handlers without having to design
all that yourself.

Join [the development group](https://telegram.me/go_telegram_bot_api) if
you want to ask questions or discuss development.

## Example

First, ensure the library is installed and up to date by running
`go get -u github.com/go-telegram-bot-api/telegram-bot-api`.

This is a very simple bot that just displays any gotten updates,
then replies it to that chat.

```go
package main

import (
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
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

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
```

There are more examples on the [site](https://go-telegram-bot-api.github.io/)
with detailed information on how to do many different kinds of things.
It's a great place to get started on using keyboards, commands, or other
kinds of reply markup.
