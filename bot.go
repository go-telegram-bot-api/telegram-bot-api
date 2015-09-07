// Package tgbotapi has bindings for interacting with the Telegram Bot API.
package tgbotapi

import (
	"net/http"
)

// BotAPI has methods for interacting with all of Telegram's Bot API endpoints.
type BotAPI struct {
	Token   string       `json:"token"`
	Debug   bool         `json:"debug"`
	Self    User         `json:"-"`
	Updates chan Update  `json:"-"`
	Client  *http.Client `json:"-"`
}

// NewBotAPI creates a new BotAPI instance.
// Requires a token, provided by @BotFather on Telegram
func NewBotAPI(token string) (*BotAPI, error) {
	return NewBotAPIWithClient(token, &http.Client{})
}

// NewBotAPIWithClient creates a new BotAPI instance passing an http.Client.
// Requires a token, provided by @BotFather on Telegram
func NewBotAPIWithClient(token string, client *http.Client) (*BotAPI, error) {
	bot := &BotAPI{
		Token:  token,
		Client: client,
	}

	self, err := bot.GetMe()
	if err != nil {
		return &BotAPI{}, err
	}

	bot.Self = self

	return bot, nil
}
