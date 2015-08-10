package tgbotapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/technoweenie/multipartstreamer"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

// Telegram constants
const (
	// APIEndpoint is the endpoint for all API methods, with formatting for Sprintf
	APIEndpoint = "https://api.telegram.org/bot%s/%s"
)

// Constant values for ChatActions
const (
	ChatTyping         = "typing"
	ChatUploadPhoto    = "upload_photo"
	ChatRecordVideo    = "record_video"
	ChatUploadVideo    = "upload_video"
	ChatRecordAudio    = "record_audio"
	ChatUploadAudio    = "upload_audio"
	ChatUploadDocument = "upload_document"
	ChatFindLocation   = "find_location"
)

// API errors
const (
	// APIForbidden happens when a token is bad
	APIForbidden = "forbidden"
)

// MessageConfig contains information about a SendMessage request.
type MessageConfig struct {
	ChatID                int
	Text                  string
	DisableWebPagePreview bool
	ReplyToMessageID      int
	ReplyMarkup           interface{}
}

// ForwardConfig contains information about a ForwardMessage request.
type ForwardConfig struct {
	ChatID     int
	FromChatID int
	MessageID  int
}

// PhotoConfig contains information about a SendPhoto request.
type PhotoConfig struct {
	ChatID           int
	Caption          string
	ReplyToMessageID int
	ReplyMarkup      interface{}
	UseExistingPhoto bool
	FilePath         string
	FileID           string
}

// AudioConfig contains information about a SendAudio request.
type AudioConfig struct {
	ChatID           int
	Duration         int
	ReplyToMessageID int
	ReplyMarkup      interface{}
	UseExistingAudio bool
	FilePath         string
	FileID           string
}

// DocumentConfig contains information about a SendDocument request.
type DocumentConfig struct {
	ChatID              int
	ReplyToMessageID    int
	ReplyMarkup         interface{}
	UseExistingDocument bool
	FilePath            string
	FileID              string
}

// StickerConfig contains information about a SendSticker request.
type StickerConfig struct {
	ChatID             int
	ReplyToMessageID   int
	ReplyMarkup        interface{}
	UseExistingSticker bool
	FilePath           string
	FileID             string
}

// VideoConfig contains information about a SendVideo request.
type VideoConfig struct {
	ChatID           int
	Duration         int
	Caption          string
	ReplyToMessageID int
	ReplyMarkup      interface{}
	UseExistingVideo bool
	FilePath         string
	FileID           string
}

// LocationConfig contains information about a SendLocation request.
type LocationConfig struct {
	ChatID           int
	Latitude         float64
	Longitude        float64
	ReplyToMessageID int
	ReplyMarkup      interface{}
}

// ChatActionConfig contains information about a SendChatAction request.
type ChatActionConfig struct {
	ChatID int
	Action string
}

// UserProfilePhotosConfig contains information about a GetUserProfilePhotos request.
type UserProfilePhotosConfig struct {
	UserID int
	Offset int
	Limit  int
}

// UpdateConfig contains information about a GetUpdates request.
type UpdateConfig struct {
	Offset  int
	Limit   int
	Timeout int
}

// WebhookConfig contains information about a SetWebhook request.
type WebhookConfig struct {
	Clear bool
	URL   *url.URL
}

// MakeRequest makes a request to a specific endpoint with our token.
// All requests are POSTs because Telegram doesn't care, and it's easier.
func (bot *BotAPI) MakeRequest(endpoint string, params url.Values) (APIResponse, error) {
	resp, err := bot.Client.PostForm(fmt.Sprintf(APIEndpoint, bot.Token, endpoint), params)
	if err != nil {
		return APIResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return APIResponse{}, errors.New(APIForbidden)
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
	f, err := os.Open(filename)
	if err != nil {
		return APIResponse{}, err
	}
	defer f.Close()

	fi, err := os.Stat(filename)
	if err != nil {
		return APIResponse{}, err
	}

	ms := multipartstreamer.New()
	ms.WriteFields(params)
	ms.WriteReader(fieldname, f.Name(), fi.Size(), f)

	req, err := http.NewRequest("POST", fmt.Sprintf(APIEndpoint, bot.Token, endpoint), nil)
	ms.SetupRequest(req)
	if err != nil {
		return APIResponse{}, err
	}

	res, err := bot.Client.Do(req)
	if err != nil {
		return APIResponse{}, err
	}
	defer res.Body.Close()

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
// Requires ChatID (destination), FromChatID (source), and MessageID.
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
func (bot *BotAPI) SendPhoto(config PhotoConfig) (Message, error) {
	if config.UseExistingPhoto {
		v := url.Values{}
		v.Add("chat_id", strconv.Itoa(config.ChatID))
		v.Add("photo", config.FileID)
		if config.Caption != "" {
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

// SendAudio sends or uploads an audio clip to a chat.
// If using a file, the file must be encoded as an .ogg with OPUS.
//
// Requires ChatID and FileID OR FilePath.
// ReplyToMessageID and ReplyMarkup are optional.
func (bot *BotAPI) SendAudio(config AudioConfig) (Message, error) {
	if config.UseExistingAudio {
		v := url.Values{}
		v.Add("chat_id", strconv.Itoa(config.ChatID))
		v.Add("audio", config.FileID)
		if config.ReplyToMessageID != 0 {
			v.Add("reply_to_message_id", strconv.Itoa(config.ReplyToMessageID))
		}
		if config.Duration != 0 {
			v.Add("duration", strconv.Itoa(config.Duration))
		}
		if config.ReplyMarkup != nil {
			data, err := json.Marshal(config.ReplyMarkup)
			if err != nil {
				return Message{}, err
			}

			v.Add("reply_markup", string(data))
		}

		resp, err := bot.MakeRequest("sendAudio", v)
		if err != nil {
			return Message{}, err
		}

		var message Message
		json.Unmarshal(resp.Result, &message)

		if bot.Debug {
			log.Printf("sendAudio req : %+v\n", v)
			log.Printf("sendAudio resp: %+v\n", message)
		}

		return message, nil
	}

	params := make(map[string]string)

	params["chat_id"] = strconv.Itoa(config.ChatID)
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

	resp, err := bot.UploadFile("sendAudio", params, "audio", config.FilePath)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.Debug {
		log.Printf("sendAudio resp: %+v\n", message)
	}

	return message, nil
}

// SendDocument sends or uploads a document to a chat.
//
// Requires ChatID and FileID OR FilePath.
// ReplyToMessageID and ReplyMarkup are optional.
func (bot *BotAPI) SendDocument(config DocumentConfig) (Message, error) {
	if config.UseExistingDocument {
		v := url.Values{}
		v.Add("chat_id", strconv.Itoa(config.ChatID))
		v.Add("document", config.FileID)
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

		resp, err := bot.MakeRequest("sendDocument", v)
		if err != nil {
			return Message{}, err
		}

		var message Message
		json.Unmarshal(resp.Result, &message)

		if bot.Debug {
			log.Printf("sendDocument req : %+v\n", v)
			log.Printf("sendDocument resp: %+v\n", message)
		}

		return message, nil
	}

	params := make(map[string]string)

	params["chat_id"] = strconv.Itoa(config.ChatID)
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

	resp, err := bot.UploadFile("sendDocument", params, "document", config.FilePath)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.Debug {
		log.Printf("sendDocument resp: %+v\n", message)
	}

	return message, nil
}

// SendSticker sends or uploads a sticker to a chat.
//
// Requires ChatID and FileID OR FilePath.
// ReplyToMessageID and ReplyMarkup are optional.
func (bot *BotAPI) SendSticker(config StickerConfig) (Message, error) {
	if config.UseExistingSticker {
		v := url.Values{}
		v.Add("chat_id", strconv.Itoa(config.ChatID))
		v.Add("sticker", config.FileID)
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

		resp, err := bot.MakeRequest("sendSticker", v)
		if err != nil {
			return Message{}, err
		}

		var message Message
		json.Unmarshal(resp.Result, &message)

		if bot.Debug {
			log.Printf("sendSticker req : %+v\n", v)
			log.Printf("sendSticker resp: %+v\n", message)
		}

		return message, nil
	}

	params := make(map[string]string)

	params["chat_id"] = strconv.Itoa(config.ChatID)
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

	resp, err := bot.UploadFile("sendSticker", params, "sticker", config.FilePath)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.Debug {
		log.Printf("sendSticker resp: %+v\n", message)
	}

	return message, nil
}

// SendVideo sends or uploads a video to a chat.
//
// Requires ChatID and FileID OR FilePath.
// ReplyToMessageID and ReplyMarkup are optional.
func (bot *BotAPI) SendVideo(config VideoConfig) (Message, error) {
	if config.UseExistingVideo {
		v := url.Values{}
		v.Add("chat_id", strconv.Itoa(config.ChatID))
		v.Add("video", config.FileID)
		if config.ReplyToMessageID != 0 {
			v.Add("reply_to_message_id", strconv.Itoa(config.ReplyToMessageID))
		}
		if config.Duration != 0 {
			v.Add("duration", strconv.Itoa(config.Duration))
		}
		if config.Caption != "" {
			v.Add("caption", config.Caption)
		}
		if config.ReplyMarkup != nil {
			data, err := json.Marshal(config.ReplyMarkup)
			if err != nil {
				return Message{}, err
			}

			v.Add("reply_markup", string(data))
		}

		resp, err := bot.MakeRequest("sendVideo", v)
		if err != nil {
			return Message{}, err
		}

		var message Message
		json.Unmarshal(resp.Result, &message)

		if bot.Debug {
			log.Printf("sendVideo req : %+v\n", v)
			log.Printf("sendVideo resp: %+v\n", message)
		}

		return message, nil
	}

	params := make(map[string]string)

	params["chat_id"] = strconv.Itoa(config.ChatID)
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

	resp, err := bot.UploadFile("sendVideo", params, "video", config.FilePath)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.Debug {
		log.Printf("sendVideo resp: %+v\n", message)
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
