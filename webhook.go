package tgbotapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// ListenForWebhook registers a http handler for a webhook.
// Useful for Google App Engine or other places where you cannot
// use a normal update chan.
func (bot *BotAPI) ListenForWebhook(config WebhookConfig) {
	bot.Updates = make(chan Update, 100)

	http.HandleFunc("/"+config.Url.Path, func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := ioutil.ReadAll(r.Body)

		var update Update
		json.Unmarshal(bytes, &update)

		bot.Updates <- update
	})
}
