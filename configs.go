package tgbotapi

import (
	"encoding/json"
	"io"
	"net/url"
	"strconv"
)

// Telegram constants
const (
	// APIEndpoint is the endpoint for all API methods, with formatting for Sprintf
	APIEndpoint = "https://api.telegram.org/bot%s/%s"
	// FileEndpoint is the endpoint for downloading a file from Telegram
	FileEndpoint = "https://api.telegram.org/file/bot%s/%s"
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

// Constant values for ParseMode in MessageConfig
const (
	ModeMarkdown = "Markdown"
)

type Chattable interface {
	Values() (url.Values, error)
	Method() string
}

type Fileable interface {
	Chattable
	Params() (map[string]string, error)
	Name() string
	GetFile() interface{}
	UseExistingFile() bool
}

// Base struct for all chat event(Message, Photo and so on)
type BaseChat struct {
	ChatID          int
	ChannelUsername string
}

func (chat *BaseChat) Values() (url.Values, error) {
	v := url.Values{}
	if chat.ChannelUsername != "" {
		v.Add("chat_id", chat.ChannelUsername)
	} else {
		v.Add("chat_id", strconv.Itoa(chat.ChatID))
	}
	return v, nil
}

type BaseFile struct {
	BaseChat
	FilePath    string
	File        interface{}
	FileID      string
	UseExisting bool
}

func (file BaseFile) Params() (map[string]string, error) {
	params := make(map[string]string)

	if file.ChannelUsername != "" {
		params["chat_id"] = file.ChannelUsername
	} else {
		params["chat_id"] = strconv.Itoa(file.ChatID)
	}

	return params, nil
}

func (file BaseFile) GetFile() interface{} {
	var result interface{}
	if file.FilePath == "" {
		result = file.File
	} else {
		result = file.FilePath
	}

	return result
}

func (file BaseFile) UseExistingFile() bool {
	return file.UseExisting
}

// MessageConfig contains information about a SendMessage request.
type MessageConfig struct {
	BaseChat
	Text                  string
	ParseMode             string
	DisableWebPagePreview bool
	ReplyToMessageID      int
	ReplyMarkup           interface{}
}

func (config MessageConfig) Values() (url.Values, error) {
	v, _ := config.BaseChat.Values()
	v.Add("text", config.Text)
	v.Add("disable_web_page_preview", strconv.FormatBool(config.DisableWebPagePreview))
	if config.ParseMode != "" {
		v.Add("parse_mode", config.ParseMode)
	}
	if config.ReplyToMessageID != 0 {
		v.Add("reply_to_message_id", strconv.Itoa(config.ReplyToMessageID))
	}
	if config.ReplyMarkup != nil {
		data, err := json.Marshal(config.ReplyMarkup)
		if err != nil {
			return v, err
		}

		v.Add("reply_markup", string(data))
	}

	return v, nil
}

func (config MessageConfig) Method() string {
	return "SendMessage"
}

// ForwardConfig contains information about a ForwardMessage request.
type ForwardConfig struct {
	BaseChat
	FromChatID          int
	FromChannelUsername string
	MessageID           int
}

func (config ForwardConfig) Values() (url.Values, error) {
	v, _ := config.BaseChat.Values()
	v.Add("message_id", strconv.Itoa(config.MessageID))
	return v, nil
}

func (config ForwardConfig) Method() string {
	return "forwardMessage"
}

// PhotoConfig contains information about a SendPhoto request.
type PhotoConfig struct {
	BaseFile
	Caption          string
	ReplyToMessageID int
	ReplyMarkup      interface{}
}

func (config PhotoConfig) Params() (map[string]string, error) {
	params, _ := config.BaseFile.Params()

	if config.Caption != "" {
		params["caption"] = config.Caption
	}
	if config.ReplyToMessageID != 0 {
		params["reply_to_message_id"] = strconv.Itoa(config.ReplyToMessageID)
	}
	if config.ReplyMarkup != nil {
		data, err := json.Marshal(config.ReplyMarkup)
		if err != nil {
			return params, err
		}

		params["reply_markup"] = string(data)
	}

	return params, nil
}

func (config PhotoConfig) Values() (url.Values, error) {
	v, _ := config.BaseChat.Values()

	v.Add("photo", config.FileID)
	if config.Caption != "" {
		v.Add("caption", config.Caption)
	}
	if config.ReplyToMessageID != 0 {
		v.Add("reply_to_message_id", strconv.Itoa(config.ReplyToMessageID))
	}
	if config.ReplyMarkup != nil {
		data, err := json.Marshal(config.ReplyMarkup)
		if err != nil {
			return v, err
		}

		v.Add("reply_markup", string(data))
	}

	return v, nil
}

func (config PhotoConfig) Name() string {
	return "photo"
}

func (config PhotoConfig) Method() string {
	return "SendPhoto"
}

// AudioConfig contains information about a SendAudio request.
type AudioConfig struct {
	BaseFile
	Duration         int
	Performer        string
	Title            string
	ReplyToMessageID int
	ReplyMarkup      interface{}
}

func (config AudioConfig) Values() (url.Values, error) {
	v, _ := config.BaseChat.Values()

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
			return v, err
		}

		v.Add("reply_markup", string(data))
	}
	if config.Performer != "" {
		v.Add("performer", config.Performer)
	}
	if config.Title != "" {
		v.Add("title", config.Title)
	}

	return v, nil
}

func (config AudioConfig) Params() (map[string]string, error) {
	params, _ := config.BaseFile.Params()

	if config.ReplyToMessageID != 0 {
		params["reply_to_message_id"] = strconv.Itoa(config.ReplyToMessageID)
	}
	if config.Duration != 0 {
		params["duration"] = strconv.Itoa(config.Duration)
	}
	if config.ReplyMarkup != nil {
		data, err := json.Marshal(config.ReplyMarkup)
		if err != nil {
			return params, err
		}

		params["reply_markup"] = string(data)
	}
	if config.Performer != "" {
		params["performer"] = config.Performer
	}
	if config.Title != "" {
		params["title"] = config.Title
	}

	return params, nil
}

func (config AudioConfig) Name() string {
	return "audio"
}

func (config AudioConfig) Method() string {
	return "SendAudio"
}

// DocumentConfig contains information about a SendDocument request.
type DocumentConfig struct {
	BaseFile
	ReplyToMessageID int
	ReplyMarkup      interface{}
}

func (config DocumentConfig) Values() (url.Values, error) {
	v, _ := config.BaseChat.Values()

	v.Add("document", config.FileID)
	if config.ReplyToMessageID != 0 {
		v.Add("reply_to_message_id", strconv.Itoa(config.ReplyToMessageID))
	}
	if config.ReplyMarkup != nil {
		data, err := json.Marshal(config.ReplyMarkup)
		if err != nil {
			return v, err
		}

		v.Add("reply_markup", string(data))
	}

	return v, nil
}

func (config DocumentConfig) Params() (map[string]string, error) {
	params, _ := config.BaseFile.Params()

	if config.ReplyToMessageID != 0 {
		params["reply_to_message_id"] = strconv.Itoa(config.ReplyToMessageID)
	}
	if config.ReplyMarkup != nil {
		data, err := json.Marshal(config.ReplyMarkup)
		if err != nil {
			return params, err
		}

		params["reply_markup"] = string(data)
	}

	return params, nil
}

func (config DocumentConfig) Name() string {
	return "document"
}

func (config DocumentConfig) Method() string {
	return "sendDocument"
}

// StickerConfig contains information about a SendSticker request.
type StickerConfig struct {
	BaseFile
	ReplyToMessageID int
	ReplyMarkup      interface{}
}

func (config StickerConfig) Values() (url.Values, error) {
	v, _ := config.BaseChat.Values()

	v.Add("sticker", config.FileID)
	if config.ReplyToMessageID != 0 {
		v.Add("reply_to_message_id", strconv.Itoa(config.ReplyToMessageID))
	}
	if config.ReplyMarkup != nil {
		data, err := json.Marshal(config.ReplyMarkup)
		if err != nil {
			return v, err
		}

		v.Add("reply_markup", string(data))
	}

	return v, nil
}

func (config StickerConfig) Params() (map[string]string, error) {
	params, _ := config.BaseFile.Params()

	if config.ReplyToMessageID != 0 {
		params["reply_to_message_id"] = strconv.Itoa(config.ReplyToMessageID)
	}
	if config.ReplyMarkup != nil {
		data, err := json.Marshal(config.ReplyMarkup)
		if err != nil {
			return params, err
		}

		params["reply_markup"] = string(data)
	}

	return params, nil
}

func (config StickerConfig) Name() string {
	return "sticker"
}

func (config StickerConfig) Method() string {
	return "sendSticker"
}

// VideoConfig contains information about a SendVideo request.
type VideoConfig struct {
	BaseFile
	Duration         int
	Caption          string
	ReplyToMessageID int
	ReplyMarkup      interface{}
}

func (config VideoConfig) Values() (url.Values, error) {
	v, _ := config.BaseChat.Values()

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
			return v, err
		}

		v.Add("reply_markup", string(data))
	}

	return v, nil
}

func (config VideoConfig) Params() (map[string]string, error) {
	params, _ := config.BaseFile.Params()

	if config.ReplyToMessageID != 0 {
		params["reply_to_message_id"] = strconv.Itoa(config.ReplyToMessageID)
	}
	if config.ReplyMarkup != nil {
		data, err := json.Marshal(config.ReplyMarkup)
		if err != nil {
			return params, err
		}

		params["reply_markup"] = string(data)
	}

	return params, nil
}

func (config VideoConfig) Name() string {
	return "viceo"
}

func (config VideoConfig) Method() string {
	return "sendVideo"
}

// VoiceConfig contains information about a SendVoice request.
type VoiceConfig struct {
	BaseFile
	Duration         int
	ReplyToMessageID int
	ReplyMarkup      interface{}
}

func (config VoiceConfig) Values() (url.Values, error) {
	v, _ := config.BaseChat.Values()

	v.Add("voice", config.FileID)
	if config.ReplyToMessageID != 0 {
		v.Add("reply_to_message_id", strconv.Itoa(config.ReplyToMessageID))
	}
	if config.Duration != 0 {
		v.Add("duration", strconv.Itoa(config.Duration))
	}
	if config.ReplyMarkup != nil {
		data, err := json.Marshal(config.ReplyMarkup)
		if err != nil {
			return v, err
		}

		v.Add("reply_markup", string(data))
	}

	return v, nil
}

func (config VoiceConfig) Params() (map[string]string, error) {
	params, _ := config.BaseFile.Params()

	if config.ReplyToMessageID != 0 {
		params["reply_to_message_id"] = strconv.Itoa(config.ReplyToMessageID)
	}
	if config.Duration != 0 {
		params["duration"] = strconv.Itoa(config.Duration)
	}
	if config.ReplyMarkup != nil {
		data, err := json.Marshal(config.ReplyMarkup)
		if err != nil {
			return params, err
		}

		params["reply_markup"] = string(data)
	}

	return params, nil
}

func (config VoiceConfig) Name() string {
	return "voice"
}

func (config VoiceConfig) Method() string {
	return "sendVoice"
}

// LocationConfig contains information about a SendLocation request.
type LocationConfig struct {
	BaseChat
	Latitude         float64
	Longitude        float64
	ReplyToMessageID int
	ReplyMarkup      interface{}
}

func (config LocationConfig) Values() (url.Values, error) {
	v, _ := config.BaseChat.Values()

	v.Add("latitude", strconv.FormatFloat(config.Latitude, 'f', 6, 64))
	v.Add("longitude", strconv.FormatFloat(config.Longitude, 'f', 6, 64))

	if config.ReplyToMessageID != 0 {
		v.Add("reply_to_message_id", strconv.Itoa(config.ReplyToMessageID))
	}
	if config.ReplyMarkup != nil {
		data, err := json.Marshal(config.ReplyMarkup)
		if err != nil {
			return v, err
		}

		v.Add("reply_markup", string(data))
	}

	return v, nil
}

func (config LocationConfig) Method() string {
	return "sendLocation"
}

// ChatActionConfig contains information about a SendChatAction request.
type ChatActionConfig struct {
	BaseChat
	Action string
}

func (config ChatActionConfig) Values() (url.Values, error) {
	v, _ := config.BaseChat.Values()
	v.Add("action", config.Action)
	return v, nil
}

// UserProfilePhotosConfig contains information about a GetUserProfilePhotos request.
type UserProfilePhotosConfig struct {
	UserID int
	Offset int
	Limit  int
}

// FileConfig has information about a file hosted on Telegram
type FileConfig struct {
	FileID string
}

// UpdateConfig contains information about a GetUpdates request.
type UpdateConfig struct {
	Offset  int
	Limit   int
	Timeout int
}

// WebhookConfig contains information about a SetWebhook request.
type WebhookConfig struct {
	Clear       bool
	URL         *url.URL
	Certificate interface{}
}

// FileBytes contains information about a set of bytes to upload as a File.
type FileBytes struct {
	Name  string
	Bytes []byte
}

// FileReader contains information about a reader to upload as a File.
// If Size is -1, it will read the entire Reader into memory to calculate a Size.
type FileReader struct {
	Name   string
	Reader io.Reader
	Size   int64
}
