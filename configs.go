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

//Chattable represents any event in chat(MessageConfig, PhotoConfig, ChatActionConfig and others)
type Chattable interface {
	Values() (url.Values, error)
	Method() string
}

//Fileable represents any file event(PhotoConfig, DocumentConfig, AudioConfig, VoiceConfig, VideoConfig, StickerConfig)
type Fileable interface {
	Chattable
	Params() (map[string]string, error)
	Name() string
	GetFile() interface{}
	UseExistingFile() bool
}

// BaseChat is base struct for all chat event(Message, Photo and so on)
type BaseChat struct {
	ChatID           int
	ChannelUsername  string
	ReplyToMessageID int
	ReplyMarkup      interface{}
}

// Values returns url.Values representation of BaseChat
func (chat *BaseChat) Values() (url.Values, error) {
	v := url.Values{}
	if chat.ChannelUsername != "" {
		v.Add("chat_id", chat.ChannelUsername)
	} else {
		v.Add("chat_id", strconv.Itoa(chat.ChatID))
	}

	if chat.ReplyToMessageID != 0 {
		v.Add("reply_to_message_id", strconv.Itoa(chat.ReplyToMessageID))
	}

	if chat.ReplyMarkup != nil {
		data, err := json.Marshal(chat.ReplyMarkup)
		if err != nil {
			return v, err
		}

		v.Add("reply_markup", string(data))
	}

	return v, nil
}

// BaseFile is base struct for all file events(PhotoConfig, DocumentConfig, AudioConfig, VoiceConfig, VideoConfig, StickerConfig)
type BaseFile struct {
	BaseChat
	FilePath    string
	File        interface{}
	FileID      string
	UseExisting bool
	MimeType    string
	FileSize    int
}

// Params returns map[string]string representation of BaseFile
func (file BaseFile) Params() (map[string]string, error) {
	params := make(map[string]string)

	if file.ChannelUsername != "" {
		params["chat_id"] = file.ChannelUsername
	} else {
		params["chat_id"] = strconv.Itoa(file.ChatID)
	}

	if file.ReplyToMessageID != 0 {
		params["reply_to_message_id"] = strconv.Itoa(file.ReplyToMessageID)
	}

	if file.ReplyMarkup != nil {
		data, err := json.Marshal(file.ReplyMarkup)
		if err != nil {
			return params, err
		}

		params["reply_markup"] = string(data)
	}

	if len(file.MimeType) > 0 {
		params["mime_type"] = file.MimeType
	}

	if file.FileSize > 0 {
		params["file_size"] = strconv.Itoa(file.FileSize)
	}

	return params, nil
}

// GetFile returns abstract representation of File inside BaseFile
func (file BaseFile) GetFile() interface{} {
	var result interface{}
	if file.FilePath == "" {
		result = file.File
	} else {
		result = file.FilePath
	}

	return result
}

// UseExistingFile returns true if BaseFile contains already uploaded file by FileID
func (file BaseFile) UseExistingFile() bool {
	return file.UseExisting
}

// MessageConfig contains information about a SendMessage request.
type MessageConfig struct {
	BaseChat
	Text                  string
	ParseMode             string
	DisableWebPagePreview bool
}

// Values returns url.Values representation of MessageConfig
func (config MessageConfig) Values() (url.Values, error) {
	v, _ := config.BaseChat.Values()
	v.Add("text", config.Text)
	v.Add("disable_web_page_preview", strconv.FormatBool(config.DisableWebPagePreview))
	if config.ParseMode != "" {
		v.Add("parse_mode", config.ParseMode)
	}

	return v, nil
}

// Method returns Telegram API method name for sending Message
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

// Values returns url.Values representation of ForwardConfig
func (config ForwardConfig) Values() (url.Values, error) {
	v, _ := config.BaseChat.Values()
	v.Add("from_chat_id", strconv.Itoa(config.FromChatID))
	v.Add("message_id", strconv.Itoa(config.MessageID))
	return v, nil
}

// Method returns Telegram API method name for sending Forward
func (config ForwardConfig) Method() string {
	return "forwardMessage"
}

// PhotoConfig contains information about a SendPhoto request.
type PhotoConfig struct {
	BaseFile
	Caption string
}

// Params returns map[string]string representation of PhotoConfig
func (config PhotoConfig) Params() (map[string]string, error) {
	params, _ := config.BaseFile.Params()

	if config.Caption != "" {
		params["caption"] = config.Caption
	}

	return params, nil
}

// Values returns url.Values representation of PhotoConfig
func (config PhotoConfig) Values() (url.Values, error) {
	v, _ := config.BaseChat.Values()

	v.Add(config.Name(), config.FileID)
	if config.Caption != "" {
		v.Add("caption", config.Caption)
	}
	return v, nil
}

// Name return field name for uploading file
func (config PhotoConfig) Name() string {
	return "photo"
}

// Method returns Telegram API method name for sending Photo
func (config PhotoConfig) Method() string {
	return "SendPhoto"
}

// AudioConfig contains information about a SendAudio request.
type AudioConfig struct {
	BaseFile
	Duration  int
	Performer string
	Title     string
}

// Values returns url.Values representation of AudioConfig
func (config AudioConfig) Values() (url.Values, error) {
	v, _ := config.BaseChat.Values()

	v.Add(config.Name(), config.FileID)
	if config.Duration != 0 {
		v.Add("duration", strconv.Itoa(config.Duration))
	}

	if config.Performer != "" {
		v.Add("performer", config.Performer)
	}
	if config.Title != "" {
		v.Add("title", config.Title)
	}

	return v, nil
}

// Params returns map[string]string representation of AudioConfig
func (config AudioConfig) Params() (map[string]string, error) {
	params, _ := config.BaseFile.Params()

	if config.Duration != 0 {
		params["duration"] = strconv.Itoa(config.Duration)
	}

	if config.Performer != "" {
		params["performer"] = config.Performer
	}
	if config.Title != "" {
		params["title"] = config.Title
	}

	return params, nil
}

// Name return field name for uploading file
func (config AudioConfig) Name() string {
	return "audio"
}

// Method returns Telegram API method name for sending Audio
func (config AudioConfig) Method() string {
	return "SendAudio"
}

// DocumentConfig contains information about a SendDocument request.
type DocumentConfig struct {
	BaseFile
}

// Values returns url.Values representation of DocumentConfig
func (config DocumentConfig) Values() (url.Values, error) {
	v, _ := config.BaseChat.Values()

	v.Add(config.Name(), config.FileID)

	return v, nil
}

// Params returns map[string]string representation of DocumentConfig
func (config DocumentConfig) Params() (map[string]string, error) {
	params, _ := config.BaseFile.Params()

	return params, nil
}

// Name return field name for uploading file
func (config DocumentConfig) Name() string {
	return "document"
}

// Method returns Telegram API method name for sending Document
func (config DocumentConfig) Method() string {
	return "sendDocument"
}

// StickerConfig contains information about a SendSticker request.
type StickerConfig struct {
	BaseFile
}

// Values returns url.Values representation of StickerConfig
func (config StickerConfig) Values() (url.Values, error) {
	v, _ := config.BaseChat.Values()

	v.Add(config.Name(), config.FileID)

	return v, nil
}

// Params returns map[string]string representation of StickerConfig
func (config StickerConfig) Params() (map[string]string, error) {
	params, _ := config.BaseFile.Params()

	return params, nil
}

// Name return field name for uploading file
func (config StickerConfig) Name() string {
	return "sticker"
}

// Method returns Telegram API method name for sending Sticker
func (config StickerConfig) Method() string {
	return "sendSticker"
}

// VideoConfig contains information about a SendVideo request.
type VideoConfig struct {
	BaseFile
	Duration int
	Caption  string
}

// Values returns url.Values representation of VideoConfig
func (config VideoConfig) Values() (url.Values, error) {
	v, _ := config.BaseChat.Values()

	v.Add(config.Name(), config.FileID)
	if config.Duration != 0 {
		v.Add("duration", strconv.Itoa(config.Duration))
	}
	if config.Caption != "" {
		v.Add("caption", config.Caption)
	}

	return v, nil
}

// Params returns map[string]string representation of VideoConfig
func (config VideoConfig) Params() (map[string]string, error) {
	params, _ := config.BaseFile.Params()

	return params, nil
}

// Name return field name for uploading file
func (config VideoConfig) Name() string {
	return "video"
}

// Method returns Telegram API method name for sending Video
func (config VideoConfig) Method() string {
	return "sendVideo"
}

// VoiceConfig contains information about a SendVoice request.
type VoiceConfig struct {
	BaseFile
	Duration int
}

// Values returns url.Values representation of VoiceConfig
func (config VoiceConfig) Values() (url.Values, error) {
	v, _ := config.BaseChat.Values()

	v.Add(config.Name(), config.FileID)
	if config.Duration != 0 {
		v.Add("duration", strconv.Itoa(config.Duration))
	}

	return v, nil
}

// Params returns map[string]string representation of VoiceConfig
func (config VoiceConfig) Params() (map[string]string, error) {
	params, _ := config.BaseFile.Params()

	if config.Duration != 0 {
		params["duration"] = strconv.Itoa(config.Duration)
	}

	return params, nil
}

// Name return field name for uploading file
func (config VoiceConfig) Name() string {
	return "voice"
}

// Method returns Telegram API method name for sending Voice
func (config VoiceConfig) Method() string {
	return "sendVoice"
}

// LocationConfig contains information about a SendLocation request.
type LocationConfig struct {
	BaseChat
	Latitude  float64
	Longitude float64
}

// Values returns url.Values representation of LocationConfig
func (config LocationConfig) Values() (url.Values, error) {
	v, _ := config.BaseChat.Values()

	v.Add("latitude", strconv.FormatFloat(config.Latitude, 'f', 6, 64))
	v.Add("longitude", strconv.FormatFloat(config.Longitude, 'f', 6, 64))

	return v, nil
}

// Method returns Telegram API method name for sending Location
func (config LocationConfig) Method() string {
	return "sendLocation"
}

// ChatActionConfig contains information about a SendChatAction request.
type ChatActionConfig struct {
	BaseChat
	Action string
}

// Values returns url.Values representation of ChatActionConfig
func (config ChatActionConfig) Values() (url.Values, error) {
	v, _ := config.BaseChat.Values()
	v.Add("action", config.Action)
	return v, nil
}

// Method returns Telegram API method name for sending ChatAction
func (config ChatActionConfig) Method() string {
	return "sendChatAction"
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
