package tgbotapi

import (
	"encoding/json"
	"strconv"
)

// GetManagedBotToken retrieves the current token of a bot managed by this bot
// (Bot API 9.6+, Managed Bots). userID is the managed bot's own id (bots are
// Users in the Bot API's type system).
func (bot *BotAPI) GetManagedBotToken(userID int64) (string, error) {
	params := Params{"user_id": strconv.FormatInt(userID, 10)}

	resp, err := bot.MakeRequest("getManagedBotToken", params)
	if err != nil {
		return "", err
	}

	var token string
	err = json.Unmarshal(resp.Result, &token)
	return token, err
}

// ReplaceManagedBotToken rotates the token of a bot managed by this bot,
// invalidating the previous one (Bot API 9.6+, Managed Bots).
func (bot *BotAPI) ReplaceManagedBotToken(userID int64) (string, error) {
	params := Params{"user_id": strconv.FormatInt(userID, 10)}

	resp, err := bot.MakeRequest("replaceManagedBotToken", params)
	if err != nil {
		return "", err
	}

	var token string
	err = json.Unmarshal(resp.Result, &token)
	return token, err
}
