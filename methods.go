package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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

type ApiResponse struct {
	Ok     bool            `json:"ok"`
	Result json.RawMessage `json:"result"`
}

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type User struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
}

type GroupChat struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type UserOrGroupChat struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
	Title     string `json:"title"`
}

type Message struct {
	MessageId           int             `json:"message_id"`
	From                User            `json:"from"`
	Date                int             `json:"date"`
	Chat                UserOrGroupChat `json:"chat"`
	ForwardFrom         User            `json:"forward_from"`
	ForwardDate         int             `json:"forward_date"`
	ReplyToMessage      *Message        `json:"reply_to_message"`
	Text                string          `json:"text"`
	Audio               Audio           `json:"audio"`
	Document            Document        `json:"document"`
	Photo               []PhotoSize     `json:"photo"`
	Sticker             Sticker         `json:"sticker"`
	Video               Video           `json:"video"`
	Contact             Contact         `json:"contact"`
	Location            Location        `json:"location"`
	NewChatParticipant  User            `json:"new_chat_participant"`
	LeftChatParticipant User            `json:"left_chat_participant"`
	NewChatTitle        string          `json:"new_chat_title"`
	NewChatPhoto        string          `json:"new_chat_photo"`
	DeleteChatPhoto     bool            `json:"delete_chat_photo"`
	GroupChatCreated    bool            `json:"group_chat_created"`
}

type PhotoSize struct {
	FileId   string `json:"file_id"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	FileSize int    `json:"file_size"`
}

type Audio struct {
	FileId   string `json:"file_id"`
	Duration int    `json:"duration"`
	MimeType string `json:"mime_type"`
	FileSize int    `json:"file_size"`
}

type Document struct {
	FileId   string    `json:"file_id"`
	Thumb    PhotoSize `json:"thumb"`
	FileName string    `json:"file_name"`
	MimeType string    `json:"mime_type"`
	FileSize int       `json:"file_size"`
}

type Sticker struct {
	FileId   string    `json:"file_id"`
	Width    int       `json:"width"`
	Height   int       `json:"height"`
	Thumb    PhotoSize `json:"thumb"`
	FileSize int       `json:"file_size"`
}

type Video struct {
	FileId   string    `json:"file_id"`
	Width    int       `json:"width"`
	Height   int       `json:"height"`
	Duration int       `json:"duration"`
	Thumb    PhotoSize `json:"thumb"`
	MimeType string    `json:"mime_type"`
	FileSize int       `json:"file_size"`
	Caption  string    `json:"caption"`
}

type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	UserId      string `json:"user_id"`
}

type Location struct {
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
}

type UserProfilePhotos struct {
	TotalCount int         `json:"total_count"`
	Photos     []PhotoSize `json:"photos"`
}

type ReplyKeyboardMarkup struct {
	Keyboard        map[string]map[string]string `json:"keyboard"`
	ResizeKeyboard  bool                         `json:"resize_keyboard"`
	OneTimeKeyboard bool                         `json:"one_time_keyboard"`
	Selective       bool                         `json:"selective"`
}

type ReplyKeyboardHide struct {
	HideKeyboard bool `json:"hide_keyboard"`
	Selective    bool `json:"selective"`
}

type ForceReply struct {
	ForceReply bool `json:"force_reply"`
	Selective  bool `json:"force_reply"`
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
		return ApiResponse{}, nil
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
