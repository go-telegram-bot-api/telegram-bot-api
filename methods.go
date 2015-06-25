package main

import (
	"bytes"
	"encoding/json"
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

type BotConfig struct {
	token string
	debug bool
}

type BotApi struct {
	config BotConfig
}

func NewBotApi(config BotConfig) *BotApi {
	return &BotApi{
		config: config,
	}
}

func (bot *BotApi) makeRequest(endpoint string, params url.Values) (ApiResponse, error) {
	resp, err := http.PostForm("https://api.telegram.org/bot"+bot.config.token+"/"+endpoint, params)
	defer resp.Body.Close()
	if err != nil {
		return ApiResponse{}, err
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ApiResponse{}, err
	}

	if bot.config.debug {
		log.Println(string(bytes[:]))
	}

	var apiResp ApiResponse
	json.Unmarshal(bytes, &apiResp)

	return apiResp, nil
}

func (bot *BotApi) uploadFile(endpoint string, params map[string]string, fieldname string, filename string) (ApiResponse, error) {
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

	req, err := http.NewRequest("POST", "https://api.telegram.org/bot"+bot.config.token+"/"+endpoint, &b)
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

	if bot.config.debug {
		log.Println(string(bytes[:]))
	}

	var apiResp ApiResponse
	json.Unmarshal(bytes, &apiResp)

	return apiResp, nil
}

func (bot *BotApi) getMe() (User, error) {
	resp, err := bot.makeRequest("getMe", nil)
	if err != nil {
		return User{}, err
	}

	var user User
	json.Unmarshal(resp.Result, &user)

	if bot.config.debug {
		log.Printf("getMe: %+v\n", user)
	}

	return user, nil
}

func (bot *BotApi) sendMessage(config MessageConfig) (Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(config.ChatId))
	v.Add("text", config.Text)
	v.Add("disable_web_page_preview", strconv.FormatBool(config.DisableWebPagePreview))
	if config.ReplyToMessageId != 0 {
		v.Add("reply_to_message_id", strconv.Itoa(config.ReplyToMessageId))
	}

	resp, err := bot.makeRequest("sendMessage", v)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.config.debug {
		log.Printf("sendMessage req : %+v\n", v)
		log.Printf("sendMessage resp: %+v\n", message)
	}

	return message, nil
}

func (bot *BotApi) forwardMessage(config ForwardConfig) (Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(config.ChatId))
	v.Add("from_chat_id", strconv.Itoa(config.FromChatId))
	v.Add("message_id", strconv.Itoa(config.MessageId))

	resp, err := bot.makeRequest("forwardMessage", v)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.config.debug {
		log.Printf("forwardMessage req : %+v\n", v)
		log.Printf("forwardMessage resp: %+v\n", message)
	}

	return message, nil
}

func (bot *BotApi) sendPhoto(config PhotoConfig) (Message, error) {
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

		resp, err := bot.makeRequest("sendPhoto", v)
		if err != nil {
			return Message{}, err
		}

		var message Message
		json.Unmarshal(resp.Result, &message)

		if bot.config.debug {
			log.Printf("sendPhoto req : %+v\n", v)
			log.Printf("sendPhoto resp: %+v\n", message)
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

	resp, err := bot.uploadFile("sendPhoto", params, "photo", config.FilePath)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.config.debug {
		log.Printf("sendPhoto resp: %+v\n", message)
	}

	return message, nil
}

func (bot *BotApi) sendAudio(config AudioConfig) (Message, error) {
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

		resp, err := bot.makeRequest("sendAudio", v)
		if err != nil {
			return Message{}, err
		}

		var message Message
		json.Unmarshal(resp.Result, &message)

		if bot.config.debug {
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

	resp, err := bot.uploadFile("sendAudio", params, "audio", config.FilePath)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.config.debug {
		log.Printf("sendAudio resp: %+v\n", message)
	}

	return message, nil
}

func (bot *BotApi) sendDocument(config DocumentConfig) (Message, error) {
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

		resp, err := bot.makeRequest("sendDocument", v)
		if err != nil {
			return Message{}, err
		}

		var message Message
		json.Unmarshal(resp.Result, &message)

		if bot.config.debug {
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

	resp, err := bot.uploadFile("sendDocument", params, "document", config.FilePath)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.config.debug {
		log.Printf("sendDocument resp: %+v\n", message)
	}

	return message, nil
}

func (bot *BotApi) sendSticker(config StickerConfig) (Message, error) {
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

		resp, err := bot.makeRequest("sendSticker", v)
		if err != nil {
			return Message{}, err
		}

		var message Message
		json.Unmarshal(resp.Result, &message)

		if bot.config.debug {
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

	resp, err := bot.uploadFile("sendSticker", params, "sticker", config.FilePath)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.config.debug {
		log.Printf("sendSticker resp: %+v\n", message)
	}

	return message, nil
}

func (bot *BotApi) sendVideo(config VideoConfig) (Message, error) {
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

		resp, err := bot.makeRequest("sendVideo", v)
		if err != nil {
			return Message{}, err
		}

		var message Message
		json.Unmarshal(resp.Result, &message)

		if bot.config.debug {
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

	resp, err := bot.uploadFile("sendVideo", params, "video", config.FilePath)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.config.debug {
		log.Printf("sendVideo resp: %+v\n", message)
	}

	return message, nil
}

func (bot *BotApi) sendLocation(config LocationConfig) (Message, error) {
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

	resp, err := bot.makeRequest("sendLocation", v)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.config.debug {
		log.Printf("sendLocation req : %+v\n", v)
		log.Printf("sendLocation resp: %+v\n", message)
	}

	return message, nil
}

func (bot *BotApi) sendChatAction(config ChatActionConfig) error {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(config.ChatId))
	v.Add("action", config.Action)

	_, err := bot.makeRequest("sendChatAction", v)
	if err != nil {
		return err
	}

	return nil
}

func (bot *BotApi) getUserProfilePhotos(config UserProfilePhotosConfig) (UserProfilePhotos, error) {
	v := url.Values{}
	v.Add("user_id", strconv.Itoa(config.UserId))
	if config.Offset != 0 {
		v.Add("offset", strconv.Itoa(config.Offset))
	}
	if config.Limit != 0 {
		v.Add("limit", strconv.Itoa(config.Limit))
	}

	resp, err := bot.makeRequest("getUserProfilePhotos", v)
	if err != nil {
		return UserProfilePhotos{}, err
	}

	var profilePhotos UserProfilePhotos
	json.Unmarshal(resp.Result, &profilePhotos)

	if bot.config.debug {
		log.Printf("getUserProfilePhotos req : %+v\n", v)
		log.Printf("getUserProfilePhotos resp: %+v\n", profilePhotos)
	}

	return profilePhotos, nil
}

func (bot *BotApi) getUpdates(config UpdateConfig) ([]Update, error) {
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

	resp, err := bot.makeRequest("getUpdates", v)
	if err != nil {
		return []Update{}, err
	}

	var updates []Update
	json.Unmarshal(resp.Result, &updates)

	if bot.config.debug {
		log.Printf("getUpdates: %+v\n", updates)
	}

	return updates, nil
}

func (bot *BotApi) setWebhook(v url.Values) error {
	_, err := bot.makeRequest("setWebhook", v)

	return err
}

type UpdateConfig struct {
	Offset  int
	Limit   int
	Timeout int
}

type MessageConfig struct {
	ChatId                int
	Text                  string
	DisableWebPagePreview bool
	ReplyToMessageId      int
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

func NewMessage(chatId int, text string) MessageConfig {
	return MessageConfig{
		ChatId: chatId,
		Text:   text,
		DisableWebPagePreview: false,
		ReplyToMessageId:      0,
	}
}

func NewForward(chatId int, fromChatId int, messageId int) ForwardConfig {
	return ForwardConfig{
		ChatId:     chatId,
		FromChatId: fromChatId,
		MessageId:  messageId,
	}
}

func NewPhotoUpload(chatId int, filename string) PhotoConfig {
	return PhotoConfig{
		ChatId:           chatId,
		UseExistingPhoto: false,
		FilePath:         filename,
	}
}

func NewPhotoShare(chatId int, fileId string) PhotoConfig {
	return PhotoConfig{
		ChatId:           chatId,
		UseExistingPhoto: true,
		FileId:           fileId,
	}
}

func NewAudioUpload(chatId int, filename string) AudioConfig {
	return AudioConfig{
		ChatId:           chatId,
		UseExistingAudio: false,
		FilePath:         filename,
	}
}

func NewAudioShare(chatId int, fileId string) AudioConfig {
	return AudioConfig{
		ChatId:           chatId,
		UseExistingAudio: true,
		FileId:           fileId,
	}
}

func NewDocumentUpload(chatId int, filename string) DocumentConfig {
	return DocumentConfig{
		ChatId:              chatId,
		UseExistingDocument: false,
		FilePath:            filename,
	}
}

func NewDocumentShare(chatId int, fileId string) DocumentConfig {
	return DocumentConfig{
		ChatId:              chatId,
		UseExistingDocument: true,
		FileId:              fileId,
	}
}

func NewStickerUpload(chatId int, filename string) StickerConfig {
	return StickerConfig{
		ChatId:             chatId,
		UseExistingSticker: false,
		FilePath:           filename,
	}
}

func NewStickerShare(chatId int, fileId string) StickerConfig {
	return StickerConfig{
		ChatId:             chatId,
		UseExistingSticker: true,
		FileId:             fileId,
	}
}

func NewVideoUpload(chatId int, filename string) VideoConfig {
	return VideoConfig{
		ChatId:           chatId,
		UseExistingVideo: false,
		FilePath:         filename,
	}
}

func NewVideoShare(chatId int, fileId string) VideoConfig {
	return VideoConfig{
		ChatId:           chatId,
		UseExistingVideo: true,
		FileId:           fileId,
	}
}

func NewLocation(chatId int, latitude float64, longitude float64) LocationConfig {
	return LocationConfig{
		ChatId:           chatId,
		Latitude:         latitude,
		Longitude:        longitude,
		ReplyToMessageId: 0,
		ReplyMarkup:      nil,
	}
}

func NewChatAction(chatId int, action string) ChatActionConfig {
	return ChatActionConfig{
		ChatId: chatId,
		Action: action,
	}
}

func NewUserProfilePhotos(userId int) UserProfilePhotosConfig {
	return UserProfilePhotosConfig{
		UserId: userId,
		Offset: 0,
		Limit:  0,
	}
}

func NewUpdate(offset int) UpdateConfig {
	return UpdateConfig{
		Offset:  offset,
		Limit:   0,
		Timeout: 0,
	}
}
