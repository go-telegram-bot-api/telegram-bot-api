package tgbotapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

// MakeRequest makes a request to a specific endpoint with our token.
// All requests are POSTs because Telegram doesn't care, and it's easier.
func (bot *BotAPI) MakeRequest(endpoint string, params url.Values) (APIResponse, error) {
	resp, err := http.PostForm("https://api.telegram.org/bot"+bot.Token+"/"+endpoint, params)
	defer resp.Body.Close()
	if err != nil {
		return APIResponse{}, err
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return APIResponse{}, err
	}

	if bot.Debug {
		log.Println(endpoint, string(bytes))
	}

	var apiResp APIResponse
	json.Unmarshal(bytes, &apiResp)

	if !apiResp.Ok {
		return APIResponse{}, errors.New(apiResp.Description)
	}

	return apiResp, nil
}

// UploadFile makes a request to the API with a file.
//
// Requires the parameter to hold the file not be in the params.
func (bot *BotAPI) UploadFile(endpoint string, params map[string]string, fieldname string, filename string) (APIResponse, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	f, err := os.Open(filename)
	if err != nil {
		return APIResponse{}, err
	}

	fw, err := w.CreateFormFile(fieldname, filename)
	if err != nil {
		return APIResponse{}, err
	}

	if _, err = io.Copy(fw, f); err != nil {
		return APIResponse{}, err
	}

	for key, val := range params {
		if fw, err = w.CreateFormField(key); err != nil {
			return APIResponse{}, err
		}

		if _, err = fw.Write([]byte(val)); err != nil {
			return APIResponse{}, err
		}
	}

	w.Close()

	req, err := http.NewRequest("POST", "https://api.telegram.org/bot"+bot.Token+"/"+endpoint, &b)
	if err != nil {
		return APIResponse{}, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return APIResponse{}, err
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return APIResponse{}, err
	}

	if bot.Debug {
		log.Println(string(bytes[:]))
	}

	var apiResp APIResponse
	json.Unmarshal(bytes, &apiResp)

	return apiResp, nil
}

// GetMe fetches the currently authenticated bot.
//
// There are no parameters for this method.
func (bot *BotAPI) GetMe() (User, error) {
	resp, err := bot.MakeRequest("getMe", nil)
	if err != nil {
		return User{}, err
	}

	var user User
	json.Unmarshal(resp.Result, &user)

	if bot.Debug {
		log.Printf("getMe: %+v\n", user)
	}

	return user, nil
}

// SendMessage sends a Message to a chat.
//
// Requires ChatID and Text.
// DisableWebPagePreview, ReplyToMessageID, and ReplyMarkup are optional.
func (bot *BotAPI) SendMessage(config MessageConfig) (Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(config.ChatID))
	v.Add("text", config.Text)
	v.Add("disable_web_page_preview", strconv.FormatBool(config.DisableWebPagePreview))
	if config.ReplyToMessageID != 0 {
		v.Add("reply_to_message_id", strconv.Itoa(config.ReplyToMessageID))
	}
	if config.ReplyMarkup != nil {
		data, err := json.Marshal(config.ReplyMarkup)
		if err != nil {
			return Message{}, err
		}

		v.Add("reply_markup", string(data))
	}

	resp, err := bot.MakeRequest("SendMessage", v)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.Debug {
		log.Printf("SendMessage req : %+v\n", v)
		log.Printf("SendMessage resp: %+v\n", message)
	}

	return message, nil
}

// ForwardMessage forwards a message from one chat to another.
//
// Requires ChatID (destionation), FromChatID (source), and MessageID.
func (bot *BotAPI) ForwardMessage(config ForwardConfig) (Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(config.ChatID))
	v.Add("from_chat_id", strconv.Itoa(config.FromChatID))
	v.Add("message_id", strconv.Itoa(config.MessageID))

	resp, err := bot.MakeRequest("forwardMessage", v)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.Debug {
		log.Printf("forwardMessage req : %+v\n", v)
		log.Printf("forwardMessage resp: %+v\n", message)
	}

	return message, nil
}

// SendPhoto sends or uploads a photo to a chat.
//
// Requires ChatID and FileID OR FilePath.
// Caption, ReplyToMessageID, and ReplyMarkup are optional.
func (bot *BotAPI) SendFile(config FileConfig) (Message, error) {
	if !config.FileType.ValidFileType() {
		return Message{}, errors.New("Invalid file type.")
	}
	if config.UseExistingPhoto {
		v := url.Values{}
		v.Add("chat_id", strconv.Itoa(config.ChatID))
		v.Add(config.FileType, config.FileID)
		if config.FileType == "photo" && config.Caption != "" {
			v.Add("caption", config.Caption)
		}
		if config.ReplyToMessageID != 0 {
			v.Add("reply_to_message_id", strconv.Itoa(config.ChatID))
		}
		if config.ReplyMarkup != nil {
			data, err := json.Marshal(config.ReplyMarkup)
			if err != nil {
				return Message{}, err
			}

			v.Add("reply_markup", string(data))
		}

		resp, err := bot.MakeRequest("SendPhoto", v)
		if err != nil {
			return Message{}, err
		}

		var message Message
		json.Unmarshal(resp.Result, &message)

		if bot.Debug {
			log.Printf("SendPhoto req : %+v\n", v)
			log.Printf("SendPhoto resp: %+v\n", message)
		}

		return message, nil
	}

	params := make(map[string]string)
	params["chat_id"] = strconv.Itoa(config.ChatID)
	if config.Caption != "" {
		params["caption"] = config.Caption
	}
	if config.ReplyToMessageID != 0 {
		params["reply_to_message_id"] = strconv.Itoa(config.ReplyToMessageID)
	}
	if config.ReplyMarkup != nil {
		data, err := json.Marshal(config.ReplyMarkup)
		if err != nil {
			return Message{}, err
		}

		params["reply_markup"] = string(data)
	}

	resp, err := bot.UploadFile("SendPhoto", params, "photo", config.FilePath)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.Debug {
		log.Printf("SendPhoto resp: %+v\n", message)
	}

	return message, nil
}

// SendLocation sends a location to a chat.
//
// Requires ChatID, Latitude, and Longitude.
// ReplyToMessageID and ReplyMarkup are optional.
func (bot *BotAPI) SendLocation(config LocationConfig) (Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(config.ChatID))
	v.Add("latitude", strconv.FormatFloat(config.Latitude, 'f', 6, 64))
	v.Add("longitude", strconv.FormatFloat(config.Longitude, 'f', 6, 64))
	if config.ReplyToMessageID != 0 {
		v.Add("reply_to_message_id", strconv.Itoa(config.ReplyToMessageID))
	}
	if config.ReplyMarkup != nil {
		data, err := json.Marshal(config.ReplyMarkup)
		if err != nil {
			return Message{}, err
		}

		v.Add("reply_markup", string(data))
	}

	resp, err := bot.MakeRequest("sendLocation", v)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.Debug {
		log.Printf("sendLocation req : %+v\n", v)
		log.Printf("sendLocation resp: %+v\n", message)
	}

	return message, nil
}

// SendChatAction sets a current action in a chat.
//
// Requires ChatID and a valid Action (see Chat constants).
func (bot *BotAPI) SendChatAction(config ChatActionConfig) error {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(config.ChatID))
	v.Add("action", config.Action)

	_, err := bot.MakeRequest("sendChatAction", v)
	if err != nil {
		return err
	}

	return nil
}

// GetUserProfilePhotos gets a user's profile photos.
//
// Requires UserID.
// Offset and Limit are optional.
func (bot *BotAPI) GetUserProfilePhotos(config UserProfilePhotosConfig) (UserProfilePhotos, error) {
	v := url.Values{}
	v.Add("user_id", strconv.Itoa(config.UserID))
	if config.Offset != 0 {
		v.Add("offset", strconv.Itoa(config.Offset))
	}
	if config.Limit != 0 {
		v.Add("limit", strconv.Itoa(config.Limit))
	}

	resp, err := bot.MakeRequest("getUserProfilePhotos", v)
	if err != nil {
		return UserProfilePhotos{}, err
	}

	var profilePhotos UserProfilePhotos
	json.Unmarshal(resp.Result, &profilePhotos)

	if bot.Debug {
		log.Printf("getUserProfilePhotos req : %+v\n", v)
		log.Printf("getUserProfilePhotos resp: %+v\n", profilePhotos)
	}

	return profilePhotos, nil
}

// GetUpdates fetches updates.
// If a WebHook is set, this will not return any data!
//
// Offset, Limit, and Timeout are optional.
// To not get old items, set Offset to one higher than the previous item.
// Set Timeout to a large number to reduce requests and get responses instantly.
func (bot *BotAPI) GetUpdates(config UpdateConfig) ([]Update, error) {
	v := url.Values{}
	if config.Offset > 0 {
		v.Add("offset", strconv.Itoa(config.Offset))
	}
	if config.Limit > 0 {
		v.Add("limit", strconv.Itoa(config.Limit))
	}
	if config.Timeout > 0 {
		v.Add("timeout", strconv.Itoa(config.Timeout))
	}

	resp, err := bot.MakeRequest("getUpdates", v)
	if err != nil {
		return []Update{}, err
	}

	var updates []Update
	json.Unmarshal(resp.Result, &updates)

	if bot.Debug {
		log.Printf("getUpdates: %+v\n", updates)
	}

	return updates, nil
}

// SetWebhook sets a webhook.
// If this is set, GetUpdates will not get any data!
//
// Requires Url OR to set Clear to true.
func (bot *BotAPI) SetWebhook(config WebhookConfig) error {
	v := url.Values{}
	if !config.Clear {
		v.Add("url", config.URL.String())
	}

	_, err := bot.MakeRequest("setWebhook", v)

	return err
}
