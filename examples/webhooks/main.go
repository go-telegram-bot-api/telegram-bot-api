package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"net/http"
)

var (
	bot *tgbotapi.BotAPI
)

const BotToken string = "YOUR-BOT-TOKEN"

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprint(w, "could not read data")
			return
		}
		var update *tgbotapi.Update
		err = json.Unmarshal(body, &update) // Not update contain your update
		
		if update.Message != nil {
			switch update.Message.Text {
			case "/start":
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "It works!"))
			}
		}
		// Enjoy
		// ...
	default:
		fmt.Fprint(w, "Method not allowed")
	}
}

func main() {
	var err error
	bot, err = tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		panic(err)
	}
	
	http.HandleFunc("/"+BotToken, handler)
	
	// If you need, you may generate a self signed certficate, as this requires HTTPS / TLS.
	// The above example tells Telegram that this is your certificate and that it should be trusted,
	// even though it is not properly signed.
	// Command:
	// $ openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 3560 -subj "//O=Org\CN=Test" -nodes
	if err := http.ListenAndServeTLS(":80", "cert.pem", "key.pem", nil); err != nil {
		panic(err)
	}
}
