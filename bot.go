// Package tgbotapi has bindings for interacting with the Telegram Bot API.
package tgbotapi

import (
	"bytes"
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
	"time"
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

func (bot *BotAPI) MakeMessageRequest(endpoint string, params url.Values) (Message, error) {
	resp, err := bot.MakeRequest(endpoint, params)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	bot.DebugLog(endpoint, params, message)

	return message, nil
}

// UploadFile makes a request to the API with a file.
//
// Requires the parameter to hold the file not be in the params.
// File should be a string to a file path, a FileBytes struct, or a FileReader struct.
func (bot *BotAPI) UploadFile(endpoint string, params map[string]string, fieldname string, file interface{}) (APIResponse, error) {
	ms := multipartstreamer.New()
	ms.WriteFields(params)

	switch f := file.(type) {
	case string:
		fileHandle, err := os.Open(f)
		if err != nil {
			return APIResponse{}, err
		}
		defer fileHandle.Close()

		fi, err := os.Stat(f)
		if err != nil {
			return APIResponse{}, err
		}

		ms.WriteReader(fieldname, fileHandle.Name(), fi.Size(), fileHandle)
	case FileBytes:
		buf := bytes.NewBuffer(f.Bytes)
		ms.WriteReader(fieldname, f.Name, int64(len(f.Bytes)), buf)
	case FileReader:
		if f.Size == -1 {
			data, err := ioutil.ReadAll(f.Reader)
			if err != nil {
				return APIResponse{}, err
			}
			buf := bytes.NewBuffer(data)

			ms.WriteReader(fieldname, f.Name, int64(len(data)), buf)

			break
		}

		ms.WriteReader(fieldname, f.Name, f.Size, f.Reader)
	default:
		return APIResponse{}, errors.New("bad file type")
	}

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

	if !apiResp.Ok {
		return APIResponse{}, errors.New(apiResp.Description)
	}

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

func (bot *BotAPI) Send(c Chattable) error {
	return nil
}

func (bot *BotAPI) DebugLog(context string, v url.Values, message interface{}) {
	if bot.Debug {
		log.Printf("%s req : %+v\n", context, v)
		log.Printf("%s resp: %+v\n", context, message)
	}
}

// SendMessage sends a Message to a chat.
//
// Requires ChatID and Text.
// DisableWebPagePreview, ReplyToMessageID, and ReplyMarkup are optional.
func (bot *BotAPI) SendMessage(config MessageConfig) (Message, error) {
	v, err := config.Values()

	if err != nil {
		return Message{}, err
	}

	message, err := bot.MakeMessageRequest("SendMessage", v)

	if err != nil {
		return Message{}, err
	}

	return message, nil
}

// ForwardMessage forwards a message from one chat to another.
//
// Requires ChatID (destination), FromChatID (source), and MessageID.
func (bot *BotAPI) ForwardMessage(config ForwardConfig) (Message, error) {
	v, _ := config.Values()

	message, err := bot.MakeMessageRequest("forwardMessage", v)
	if err != nil {
		return Message{}, err
	}

	return message, nil
}

// SendPhoto sends or uploads a photo to a chat.
//
// Requires ChatID and FileID OR File.
// Caption, ReplyToMessageID, and ReplyMarkup are optional.
// File should be either a string, FileBytes, or FileReader.
func (bot *BotAPI) SendPhoto(config PhotoConfig) (Message, error) {
	if config.UseExistingPhoto {
		v, err := config.Values()

		if err != nil {
			return Message{}, err
		}

		message, err := bot.MakeMessageRequest("SendPhoto", v)
		if err != nil {
			return Message{}, err
		}

		return message, nil
	}

	params := make(map[string]string)
	if config.ChannelUsername != "" {
		params["chat_id"] = config.ChannelUsername
	} else {
		params["chat_id"] = strconv.Itoa(config.ChatID)
	}
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

	var file interface{}
	if config.FilePath == "" {
		file = config.File
	} else {
		file = config.FilePath
	}

	resp, err := bot.UploadFile("SendPhoto", params, "photo", file)
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
// If using a file, the file must be in the .mp3 format.
//
// When the fields title and performer are both empty and
// the mime-type of the file to be sent is not audio/mpeg,
// the file must be an .ogg file encoded with OPUS.
// You may use the tgutils.EncodeAudio func to assist you with this, if needed.
//
// Requires ChatID and FileID OR File.
// ReplyToMessageID and ReplyMarkup are optional.
// File should be either a string, FileBytes, or FileReader.
func (bot *BotAPI) SendAudio(config AudioConfig) (Message, error) {
	if config.UseExistingAudio {
		v, err := config.Values()
		if err != nil {
			return Message{}, err
		}

		message, err := bot.MakeMessageRequest("sendAudio", v)
		if err != nil {
			return Message{}, err
		}

		return message, nil
	}

	params := make(map[string]string)

	if config.ChannelUsername != "" {
		params["chat_id"] = config.ChannelUsername
	} else {
		params["chat_id"] = strconv.Itoa(config.ChatID)
	}
	if config.ReplyToMessageID != 0 {
		params["reply_to_message_id"] = strconv.Itoa(config.ReplyToMessageID)
	}
	if config.Duration != 0 {
		params["duration"] = strconv.Itoa(config.Duration)
	}
	if config.ReplyMarkup != nil {
		data, err := json.Marshal(config.ReplyMarkup)
		if err != nil {
			return Message{}, err
		}

		params["reply_markup"] = string(data)
	}
	if config.Performer != "" {
		params["performer"] = config.Performer
	}
	if config.Title != "" {
		params["title"] = config.Title
	}

	var file interface{}
	if config.FilePath == "" {
		file = config.File
	} else {
		file = config.FilePath
	}

	resp, err := bot.UploadFile("sendAudio", params, "audio", file)
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
// Requires ChatID and FileID OR File.
// ReplyToMessageID and ReplyMarkup are optional.
// File should be either a string, FileBytes, or FileReader.
func (bot *BotAPI) SendDocument(config DocumentConfig) (Message, error) {
	if config.UseExistingDocument {
		v, err := config.Values()
		if err != nil {
			return Message{}, err
		}

		message, err := bot.MakeMessageRequest("sendDocument", v)
		if err != nil {
			return Message{}, err
		}

		return message, nil
	}

	params := make(map[string]string)

	if config.ChannelUsername != "" {
		params["chat_id"] = config.ChannelUsername
	} else {
		params["chat_id"] = strconv.Itoa(config.ChatID)
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

	var file interface{}
	if config.FilePath == "" {
		file = config.File
	} else {
		file = config.FilePath
	}

	resp, err := bot.UploadFile("sendDocument", params, "document", file)
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

// SendVoice sends or uploads a playable voice to a chat.
// If using a file, the file must be encoded as an .ogg with OPUS.
// You may use the tgutils.EncodeAudio func to assist you with this, if needed.
//
// Requires ChatID and FileID OR File.
// ReplyToMessageID and ReplyMarkup are optional.
// File should be either a string, FileBytes, or FileReader.
func (bot *BotAPI) SendVoice(config VoiceConfig) (Message, error) {
	if config.UseExistingVoice {
		v, err := config.Values()
		if err != nil {
			return Message{}, err
		}

		message, err := bot.MakeMessageRequest("sendVoice", v)
		if err != nil {
			return Message{}, err
		}

		return message, nil
	}

	params := make(map[string]string)

	if config.ChannelUsername != "" {
		params["chat_id"] = config.ChannelUsername
	} else {
		params["chat_id"] = strconv.Itoa(config.ChatID)
	}
	if config.ReplyToMessageID != 0 {
		params["reply_to_message_id"] = strconv.Itoa(config.ReplyToMessageID)
	}
	if config.Duration != 0 {
		params["duration"] = strconv.Itoa(config.Duration)
	}
	if config.ReplyMarkup != nil {
		data, err := json.Marshal(config.ReplyMarkup)
		if err != nil {
			return Message{}, err
		}

		params["reply_markup"] = string(data)
	}

	var file interface{}
	if config.FilePath == "" {
		file = config.File
	} else {
		file = config.FilePath
	}

	resp, err := bot.UploadFile("SendVoice", params, "voice", file)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.Debug {
		log.Printf("SendVoice resp: %+v\n", message)
	}

	return message, nil
}

// SendSticker sends or uploads a sticker to a chat.
//
// Requires ChatID and FileID OR File.
// ReplyToMessageID and ReplyMarkup are optional.
// File should be either a string, FileBytes, or FileReader.
func (bot *BotAPI) SendSticker(config StickerConfig) (Message, error) {
	if config.UseExistingSticker {
		v, err := config.Values()
		if err != nil {
			return Message{}, err
		}

		message, err := bot.MakeMessageRequest("sendSticker", v)
		if err != nil {
			return Message{}, err
		}

		return message, nil
	}

	params := make(map[string]string)

	if config.ChannelUsername != "" {
		params["chat_id"] = config.ChannelUsername
	} else {
		params["chat_id"] = strconv.Itoa(config.ChatID)
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

	var file interface{}
	if config.FilePath == "" {
		file = config.File
	} else {
		file = config.FilePath
	}

	resp, err := bot.UploadFile("sendSticker", params, "sticker", file)
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
// Requires ChatID and FileID OR File.
// ReplyToMessageID and ReplyMarkup are optional.
// File should be either a string, FileBytes, or FileReader.
func (bot *BotAPI) SendVideo(config VideoConfig) (Message, error) {
	if config.UseExistingVideo {
		v, err := config.Values()
		if err != nil {
			return Message{}, err
		}

		message, err := bot.MakeMessageRequest("sendVideo", v)
		if err != nil {
			return Message{}, err
		}

		return message, nil
	}

	params := make(map[string]string)

	if config.ChannelUsername != "" {
		params["chat_id"] = config.ChannelUsername
	} else {
		params["chat_id"] = strconv.Itoa(config.ChatID)
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

	var file interface{}
	if config.FilePath == "" {
		file = config.File
	} else {
		file = config.FilePath
	}

	resp, err := bot.UploadFile("sendVideo", params, "video", file)
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
	v, err := config.Values()
	if err != nil {
		return Message{}, err
	}

	message, err := bot.MakeMessageRequest("sendLocation", v)
	if err != nil {
		return Message{}, err
	}

	return message, nil
}

// SendChatAction sets a current action in a chat.
//
// Requires ChatID and a valid Action (see Chat constants).
func (bot *BotAPI) SendChatAction(config ChatActionConfig) error {
	v, _ := config.Values()

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

	bot.DebugLog("GetUserProfilePhoto", v, profilePhotos)

	return profilePhotos, nil
}

// GetFile returns a file_id required to download a file.
//
// Requires FileID.
func (bot *BotAPI) GetFile(config FileConfig) (File, error) {
	v := url.Values{}
	v.Add("file_id", config.FileID)

	resp, err := bot.MakeRequest("getFile", v)
	if err != nil {
		return File{}, err
	}

	var file File
	json.Unmarshal(resp.Result, &file)

	bot.DebugLog("GetFile", v, file)

	return file, nil
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
// Requires URL OR to set Clear to true.
func (bot *BotAPI) SetWebhook(config WebhookConfig) (APIResponse, error) {
	if config.Certificate == nil {
		v := url.Values{}
		if !config.Clear {
			v.Add("url", config.URL.String())
		}

		return bot.MakeRequest("setWebhook", v)
	}

	params := make(map[string]string)
	params["url"] = config.URL.String()

	resp, err := bot.UploadFile("setWebhook", params, "certificate", config.Certificate)
	if err != nil {
		return APIResponse{}, err
	}

	var apiResp APIResponse
	json.Unmarshal(resp.Result, &apiResp)

	if bot.Debug {
		log.Printf("setWebhook resp: %+v\n", apiResp)
	}

	return apiResp, nil
}

// UpdatesChan starts a channel for getting updates.
func (bot *BotAPI) UpdatesChan(config UpdateConfig) error {
	bot.Updates = make(chan Update, 100)

	go func() {
		for {
			updates, err := bot.GetUpdates(config)
			if err != nil {
				log.Println(err)
				log.Println("Failed to get updates, retrying in 3 seconds...")
				time.Sleep(time.Second * 3)

				continue
			}

			for _, update := range updates {
				if update.UpdateID >= config.Offset {
					config.Offset = update.UpdateID + 1
					bot.Updates <- update
				}
			}
		}
	}()

	return nil
}

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
