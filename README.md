# Golang bindings for the Telegram Bot API

[![GoDoc](https://godoc.org/github.com/go-telegram-bot-api/telegram-bot-api?status.svg)](http://godoc.org/github.com/go-telegram-bot-api/telegram-bot-api)
[![Travis](https://travis-ci.org/go-telegram-bot-api/telegram-bot-api.svg)](https://travis-ci.org/go-telegram-bot-api/telegram-bot-api)

All methods have been added, and all features should be available.
If you want a feature that hasn't been added yet or something is broken,
open an issue and I'll see what I can do.

All methods are fairly self explanatory, and reading the godoc page should
explain everything. If something isn't clear, open an issue or submit
a pull request.

The scope of this project is just to provide a wrapper around the API
without any additional features. There are other projects for creating
something with plugins and command handlers without having to design
all that yourself.

Use `github.com/go-telegram-bot-api/telegram-bot-api` for the latest
version, or use `gopkg.in/telegram-bot-api.v4` for the stable build.

Join [the development group](https://telegram.me/go_telegram_bot_api) if
you want to ask questions or discuss development.

## Example

This is a very simple bot that just displays any gotten updates,
then replies it to that chat.

```go
package main

import (
	"log"
	"gopkg.in/telegram-bot-api.v4"
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
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
```

If you need to use webhooks (if you wish to run on Google App Engine),
you may use a slightly different method.

```go
package main

import (
	"flag"
	"log"
	"net/http"

	"gopkg.in/telegram-bot-api.v4"
)

var (
	reverse_proxy = flag.Bool("reverse_proxy", false, "Used reverse proxy (e.g., nginx)")
)

func main() {
	flag.Parse()

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

	updates := bot.ListenForWebhook("/" + bot.Token)
	go func() {
		var err error
		if *reverse_proxy {
			err = http.ListenAndServe("127.0.0.1:8444", nil)
		} else {
			err = http.ListenAndServeTLS("0.0.0.0:8443", "cert.pem", "key.pem", nil)
		}
		log.Fatal(err)
	}()

	for update := range updates {
		log.Printf("%+v\n", update)
	}
}
```

Example nginx config for multiple bots

```nginx
server {
	listen 8443 ssl;
	# listen [::]:80 default_server;
	server_name _;

	# ssl on;
	ssl_certificate		 /etc/nginx/cert/cert.pem;
	ssl_certificate_key  /etc/nginx/cert/key.pem;

	ssl_protocols TLSv1 TLSv1.1 TLSv1.2;

	root /usr/share/nginx/html;
	index index.html;

	# MyAwesomeBotToken
	location ^~ /MyAwesomeBotToken {
		proxy_pass http://127.0.0.1:8444;
		proxy_redirect	   off;
		proxy_set_header   Host $host;
		proxy_set_header   X-Real-IP $remote_addr;
	}

	# MyAwesomeBotToken2
	location ^~ /MyAwesomeBotToken2 {
		proxy_pass http://127.0.0.1:8445;
		proxy_redirect	   off;
		proxy_set_header   Host $host;
		proxy_set_header   X-Real-IP $remote_addr;
	}

	location / {
			# First attempt to serve request as file, then
			# as directory, then fall back to displaying a 404.
			try_files $uri $uri/ =404;
	}
}
```

If you need, you may generate a self signed certficate, as this requires
HTTPS / TLS. The above example tells Telegram that this is your
certificate and that it should be trusted, even though it is not
properly signed.

    openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 3560 -subj "//O=Org\CN=Test" -nodes

Now that [Let's Encrypt](https://letsencrypt.org) has entered public beta,
you may wish to generate your free TLS certificate there.
