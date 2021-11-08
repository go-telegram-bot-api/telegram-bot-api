# Golang bindings for the Telegram Bot API

[![Go Reference](https://pkg.go.dev/badge/github.com/go-telegram-bot-api/telegram-bot-api/v5.svg)](https://pkg.go.dev/github.com/go-telegram-bot-api/telegram-bot-api/v5)
[![Test](https://github.com/go-telegram-bot-api/telegram-bot-api/actions/workflows/test.yml/badge.svg)](https://github.com/go-telegram-bot-api/telegram-bot-api/actions/workflows/test.yml)

All methods are fairly self explanatory, and reading the [godoc](http://godoc.org/github.com/go-telegram-bot-api/telegram-bot-api) page should
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
`go get -u github.com/go-telegram-bot-api/telegram-bot-api/v5`.

This is a very simple bot that just displays any gotten updates,
then replies it to that chat.

```go
package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

There are more examples on the [site](https://go-telegram-bot-api.dev/)
with detailed information on how to do many different kinds of things.
It's a great place to get started on using keyboards, commands, or other
kinds of reply markup.

If you need to use webhooks (if you wish to run on Google App Engine),
you may use a slightly different method.

```go
package main

import (
	"log"
	"net/http"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("MyAwesomeBotToken")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhookWithCert("https://www.google.com:8443/"+bot.Token, "cert.pem"))
	if err != nil {
		log.Fatal(err)
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServeTLS("0.0.0.0:8443", "cert.pem", "key.pem", nil)

	for update := range updates {
		log.Printf("%+v\n", update)
	}
}
```

If you need to publish your bot on AWS Lambda(or something like it) and AWS API Gateway,
you can use such example:

In this code used AWS Lambda Go net/http server adapter [algnhsa](https://github.com/akrylysov/algnhsa)

```go
package main

import (
	"github.com/akrylysov/algnhsa"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
)

func answer(w http.ResponseWriter, r *http.Request) {
	bot, err := tgbotapi.NewBotAPI("MyAwesomeBotToken")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	updates := bot.ListenForWebhookRespReqFormat(w, r)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID
		_, err := bot.Send(msg)
		if err != nil {
			log.Printf("Error send message: %s | Error: %s", msg.Text, err.Error())
		}
	}
}

func setWebhook(_ http.ResponseWriter, _ *http.Request) {
	bot, err := tgbotapi.NewBotAPI("MyAwesomeBotToken")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook("https://your_api_gateway_address.com/"+bot.Token))
	if err != nil {
		log.Fatal(err)
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}
}

func main() {
	http.HandleFunc("/set_webhook", setWebhook)
	http.HandleFunc("/MyAwesomeBotToken", answer)
	algnhsa.ListenAndServe(http.DefaultServeMux, nil)
}
```

If you need, you may generate a self signed certficate, as this requires
HTTPS / TLS. The above example tells Telegram that this is your
certificate and that it should be trusted, even though it is not
properly signed.

    openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 3560 -subj "//O=Org\CN=Test" -nodes

Now that [Let's Encrypt](https://letsencrypt.org) is available,
you may wish to generate your free TLS certificate there.
