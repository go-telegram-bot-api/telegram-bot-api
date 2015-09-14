package tgbotapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// ListenForWebhook registers a http handler for a webhook.
func (bot *BotAPI) ListenForWebhook(pattern string) {
	bot.Updates = make(chan Update, 100)

	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := ioutil.ReadAll(r.Body)

		var update Update
		json.Unmarshal(bytes, &update)

		bot.Updates <- update
	})
}
