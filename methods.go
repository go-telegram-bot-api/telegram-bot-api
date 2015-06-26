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

const (
	CHAT_TYPING          = "typing"
	CHAT_UPLOAD_PHOTO    = "upload_photo"
	CHAT_RECORD_VIDEO    = "record_video"
	CHAT_UPLOAD_VIDEO    = "upload_video"
	CHAT_RECORD_AUDIO    = "record_audio"
	CHAT_UPLOAD_AUDIO    = "upload_audio"
	CHAT_UPLOAD_DOCUMENT = "upload_document"
	CHAT_FIND_LOCATION   = "find_location"
)

type MessageConfig struct {
	ChatId                int
	Text                  string
	DisableWebPagePreview bool
	ReplyToMessageId      int
	ReplyMarkup           interface{}
}

type ForwardConfig struct {
	ChatId     int
	FromChatId int
	MessageId  int
}

type PhotoConfig struct {
	ChatId           int
	Caption          string
	ReplyToMessageId int
	ReplyMarkup      interface{}
	UseExistingPhoto bool
	FilePath         string
	FileId           string
}

type AudioConfig struct {
	ChatId           int
	ReplyToMessageId int
	ReplyMarkup      interface{}
	UseExistingAudio bool
	FilePath         string
	FileId           string
}

type DocumentConfig struct {
	ChatId              int
	ReplyToMessageId    int
	ReplyMarkup         interface{}
	UseExistingDocument bool
	FilePath            string
	FileId              string
}

type StickerConfig struct {
	ChatId             int
	ReplyToMessageId   int
	ReplyMarkup        interface{}
	UseExistingSticker bool
	FilePath           string
	FileId             string
}

type VideoConfig struct {
	ChatId           int
	ReplyToMessageId int
	ReplyMarkup      interface{}
	UseExistingVideo bool
	FilePath         string
	FileId           string
}

type LocationConfig struct {
	ChatId           int
	Latitude         float64
	Longitude        float64
	ReplyToMessageId int
	ReplyMarkup      interface{}
}

type ChatActionConfig struct {
	ChatId int
	Action string
}

type UserProfilePhotosConfig struct {
	UserId int
	Offset int
	Limit  int
}

type UpdateConfig struct {
	Offset  int
	Limit   int
	Timeout int
}

type WebhookConfig struct {
	Clear bool
	Url   *url.URL
}

func (bot *BotApi) MakeRequest(endpoint string, params url.Values) (ApiResponse, error) {
	resp, err := http.PostForm("https://api.telegram.org/bot"+bot.Token+"/"+endpoint, params)
	defer resp.Body.Close()
	if err != nil {
		return ApiResponse{}, err
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ApiResponse{}, err
	}

	if bot.Debug {
		log.Println(endpoint, string(bytes))
	}

	var apiResp ApiResponse
	json.Unmarshal(bytes, &apiResp)

	if !apiResp.Ok {
		return ApiResponse{}, errors.New(apiResp.Description)
	}

	return apiResp, nil
}

func (bot *BotApi) UploadFile(endpoint string, params map[string]string, fieldname string, filename string) (ApiResponse, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	f, err := os.Open(filename)
	if err != nil {
		return ApiResponse{}, err
	}

	fw, err := w.CreateFormFile(fieldname, filename)
	if err != nil {
		return ApiResponse{}, err
	}

	if _, err = io.Copy(fw, f); err != nil {
		return ApiResponse{}, err
	}

	for key, val := range params {
		if fw, err = w.CreateFormField(key); err != nil {
			return ApiResponse{}, err
		}

		if _, err = fw.Write([]byte(val)); err != nil {
			return ApiResponse{}, err
		}
	}

	w.Close()

	req, err := http.NewRequest("POST", "https://api.telegram.org/bot"+bot.Token+"/"+endpoint, &b)
	if err != nil {
		return ApiResponse{}, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return ApiResponse{}, err
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ApiResponse{}, err
	}

	if bot.Debug {
		log.Println(string(bytes[:]))
	}

	var apiResp ApiResponse
	json.Unmarshal(bytes, &apiResp)

	return apiResp, nil
}

func (bot *BotApi) GetMe() (User, error) {
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

func (bot *BotApi) SendMessage(config MessageConfig) (Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(config.ChatId))
	v.Add("text", config.Text)
	v.Add("disable_web_page_preview", strconv.FormatBool(config.DisableWebPagePreview))
	if config.ReplyToMessageId != 0 {
		v.Add("reply_to_message_id", strconv.Itoa(config.ReplyToMessageId))
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

func (bot *BotApi) ForwardMessage(config ForwardConfig) (Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(config.ChatId))
	v.Add("from_chat_id", strconv.Itoa(config.FromChatId))
	v.Add("message_id", strconv.Itoa(config.MessageId))

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

func (bot *BotApi) SendPhoto(config PhotoConfig) (Message, error) {
	if config.UseExistingPhoto {
		v := url.Values{}
		v.Add("chat_id", strconv.Itoa(config.ChatId))
		v.Add("photo", config.FileId)
		if config.Caption != "" {
			v.Add("caption", config.Caption)
		}
		if config.ReplyToMessageId != 0 {
			v.Add("reply_to_message_id", strconv.Itoa(config.ChatId))
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
	params["chat_id"] = strconv.Itoa(config.ChatId)
	if config.Caption != "" {
		params["caption"] = config.Caption
	}
	if config.ReplyToMessageId != 0 {
		params["reply_to_message_id"] = strconv.Itoa(config.ReplyToMessageId)
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

func (bot *BotApi) SendAudio(config AudioConfig) (Message, error) {
	if config.UseExistingAudio {
		v := url.Values{}
		v.Add("chat_id", strconv.Itoa(config.ChatId))
		v.Add("audio", config.FileId)
		if config.ReplyToMessageId != 0 {
			v.Add("reply_to_message_id", strconv.Itoa(config.ReplyToMessageId))
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

	params["chat_id"] = strconv.Itoa(config.ChatId)
	if config.ReplyToMessageId != 0 {
		params["reply_to_message_id"] = strconv.Itoa(config.ReplyToMessageId)
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

func (bot *BotApi) SendDocument(config DocumentConfig) (Message, error) {
	if config.UseExistingDocument {
		v := url.Values{}
		v.Add("chat_id", strconv.Itoa(config.ChatId))
		v.Add("document", config.FileId)
		if config.ReplyToMessageId != 0 {
			v.Add("reply_to_message_id", strconv.Itoa(config.ReplyToMessageId))
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

	params["chat_id"] = strconv.Itoa(config.ChatId)
	if config.ReplyToMessageId != 0 {
		params["reply_to_message_id"] = strconv.Itoa(config.ReplyToMessageId)
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

func (bot *BotApi) SendSticker(config StickerConfig) (Message, error) {
	if config.UseExistingSticker {
		v := url.Values{}
		v.Add("chat_id", strconv.Itoa(config.ChatId))
		v.Add("sticker", config.FileId)
		if config.ReplyToMessageId != 0 {
			v.Add("reply_to_message_id", strconv.Itoa(config.ReplyToMessageId))
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

	params["chat_id"] = strconv.Itoa(config.ChatId)
	if config.ReplyToMessageId != 0 {
		params["reply_to_message_id"] = strconv.Itoa(config.ReplyToMessageId)
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

func (bot *BotApi) SendVideo(config VideoConfig) (Message, error) {
	if config.UseExistingVideo {
		v := url.Values{}
		v.Add("chat_id", strconv.Itoa(config.ChatId))
		v.Add("video", config.FileId)
		if config.ReplyToMessageId != 0 {
			v.Add("reply_to_message_id", strconv.Itoa(config.ReplyToMessageId))
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

	params["chat_id"] = strconv.Itoa(config.ChatId)
	if config.ReplyToMessageId != 0 {
		params["reply_to_message_id"] = strconv.Itoa(config.ReplyToMessageId)
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

func (bot *BotApi) SendLocation(config LocationConfig) (Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(config.ChatId))
	v.Add("latitude", strconv.FormatFloat(config.Latitude, 'f', 6, 64))
	v.Add("longitude", strconv.FormatFloat(config.Longitude, 'f', 6, 64))
	if config.ReplyToMessageId != 0 {
		v.Add("reply_to_message_id", strconv.Itoa(config.ReplyToMessageId))
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

func (bot *BotApi) SendChatAction(config ChatActionConfig) error {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(config.ChatId))
	v.Add("action", config.Action)

	_, err := bot.MakeRequest("sendChatAction", v)
	if err != nil {
		return err
	}

	return nil
}

func (bot *BotApi) GetUserProfilePhotos(config UserProfilePhotosConfig) (UserProfilePhotos, error) {
	v := url.Values{}
	v.Add("user_id", strconv.Itoa(config.UserId))
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

func (bot *BotApi) GetUpdates(config UpdateConfig) ([]Update, error) {
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

func (bot *BotApi) SetWebhook(config WebhookConfig) error {
	v := url.Values{}
	if !config.Clear {
		v.Add("url", config.Url.String())
	}

	_, err := bot.MakeRequest("setWebhook", v)

	return err
}
