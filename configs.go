package tgbotapi

import (
	"encoding/json"
	"io"
	"log"
	"net/url"
	"strconv"
)

// Telegram constants
const (
	// APIEndpoint is the endpoint for all API methods,
	// with formatting for Sprintf.
	APIEndpoint = "https://api.telegram.org/bot%s/%s"
	// FileEndpoint is the endpoint for downloading a file from Telegram.
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

// Chattable is any config type that can be sent.
type Chattable interface {
	values() (url.Values, error)
	method() string
}

// Fileable is any config type that can be sent that includes a file.
type Fileable interface {
	Chattable
	params() (map[string]string, error)
	name() string
	getFile() interface{}
	useExistingFile() bool
}

// BaseChat is base type for all chat config types.
type BaseChat struct {
	ChatID           int // required
	ChannelUsername  string
	ReplyToMessageID int
	ReplyMarkup      interface{}
}

// values returns url.Values representation of BaseChat
func (chat *BaseChat) values() (url.Values, error) {
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

// BaseFile is a base type for all file config types.
type BaseFile struct {
	BaseChat
	FilePath    string
	File        interface{}
	FileID      string
	UseExisting bool
	MimeType    string
	FileSize    int
}

// params returns a map[string]string representation of BaseFile.
func (file BaseFile) params() (map[string]string, error) {
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

	if file.MimeType != "" {
		params["mime_type"] = file.MimeType
	}

	if file.FileSize > 0 {
		params["file_size"] = strconv.Itoa(file.FileSize)
	}

	return params, nil
}

// getFile returns the file.
func (file BaseFile) getFile() interface{} {
	var result interface{}
	if file.FilePath == "" {
		result = file.File
	} else {
		log.Println("FilePath is deprecated.")
		log.Println("Please use BaseFile.File instead.")
		result = file.FilePath
	}

	return result
}

// useExistingFile returns if the BaseFile has already been uploaded.
func (file BaseFile) useExistingFile() bool {
	return file.UseExisting
}

// MessageConfig contains information about a SendMessage request.
type MessageConfig struct {
	BaseChat
	Text                  string
	ParseMode             string
	DisableWebPagePreview bool
}

// values returns a url.Values representation of MessageConfig.
func (config MessageConfig) values() (url.Values, error) {
	v, _ := config.BaseChat.values()
	v.Add("text", config.Text)
	v.Add("disable_web_page_preview", strconv.FormatBool(config.DisableWebPagePreview))
	if config.ParseMode != "" {
		v.Add("parse_mode", config.ParseMode)
	}

	return v, nil
}

// method returns Telegram API method name for sending Message.
func (config MessageConfig) method() string {
	return "sendMessage"
}

// ForwardConfig contains information about a ForwardMessage request.
type ForwardConfig struct {
	BaseChat
	FromChatID          int // required
	FromChannelUsername string
	MessageID           int // required
}

// values returns a url.Values representation of ForwardConfig.
func (config ForwardConfig) values() (url.Values, error) {
	v, _ := config.BaseChat.values()
	v.Add("from_chat_id", strconv.Itoa(config.FromChatID))
	v.Add("message_id", strconv.Itoa(config.MessageID))
	return v, nil
}

// method returns Telegram API method name for sending Forward.
func (config ForwardConfig) method() string {
	return "forwardMessage"
}

// PhotoConfig contains information about a SendPhoto request.
type PhotoConfig struct {
	BaseFile
	Caption string
}

// Params returns a map[string]string representation of PhotoConfig.
func (config PhotoConfig) params() (map[string]string, error) {
	params, _ := config.BaseFile.params()

	if config.Caption != "" {
		params["caption"] = config.Caption
	}

	return params, nil
}

// Values returns a url.Values representation of PhotoConfig.
func (config PhotoConfig) values() (url.Values, error) {
	v, _ := config.BaseChat.values()

	v.Add(config.name(), config.FileID)
	if config.Caption != "" {
		v.Add("caption", config.Caption)
	}
	return v, nil
}

// name returns the field name for the Photo.
func (config PhotoConfig) name() string {
	return "photo"
}

// method returns Telegram API method name for sending Photo.
func (config PhotoConfig) method() string {
	return "sendPhoto"
}

// AudioConfig contains information about a SendAudio request.
type AudioConfig struct {
	BaseFile
	Duration  int
	Performer string
	Title     string
}

// values returns a url.Values representation of AudioConfig.
func (config AudioConfig) values() (url.Values, error) {
	v, _ := config.BaseChat.values()

	v.Add(config.name(), config.FileID)
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

// params returns a map[string]string representation of AudioConfig.
func (config AudioConfig) params() (map[string]string, error) {
	params, _ := config.BaseFile.params()

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

// name returns the field name for the Audio.
func (config AudioConfig) name() string {
	return "audio"
}

// method returns Telegram API method name for sending Audio.
func (config AudioConfig) method() string {
	return "sendAudio"
}

// DocumentConfig contains information about a SendDocument request.
type DocumentConfig struct {
	BaseFile
}

// values returns a url.Values representation of DocumentConfig.
func (config DocumentConfig) values() (url.Values, error) {
	v, _ := config.BaseChat.values()

	v.Add(config.name(), config.FileID)

	return v, nil
}

// params returns a map[string]string representation of DocumentConfig.
func (config DocumentConfig) params() (map[string]string, error) {
	params, _ := config.BaseFile.params()

	return params, nil
}

// name returns the field name for the Document.
func (config DocumentConfig) name() string {
	return "document"
}

// method returns Telegram API method name for sending Document.
func (config DocumentConfig) method() string {
	return "sendDocument"
}

// StickerConfig contains information about a SendSticker request.
type StickerConfig struct {
	BaseFile
}

// values returns a url.Values representation of StickerConfig.
func (config StickerConfig) values() (url.Values, error) {
	v, _ := config.BaseChat.values()

	v.Add(config.name(), config.FileID)

	return v, nil
}

// params returns a map[string]string representation of StickerConfig.
func (config StickerConfig) params() (map[string]string, error) {
	params, _ := config.BaseFile.params()

	return params, nil
}

// name returns the field name for the Sticker.
func (config StickerConfig) name() string {
	return "sticker"
}

// method returns Telegram API method name for sending Sticker.
func (config StickerConfig) method() string {
	return "sendSticker"
}

// VideoConfig contains information about a SendVideo request.
type VideoConfig struct {
	BaseFile
	Duration int
	Caption  string
}

// values returns a url.Values representation of VideoConfig.
func (config VideoConfig) values() (url.Values, error) {
	v, _ := config.BaseChat.values()

	v.Add(config.name(), config.FileID)
	if config.Duration != 0 {
		v.Add("duration", strconv.Itoa(config.Duration))
	}
	if config.Caption != "" {
		v.Add("caption", config.Caption)
	}

	return v, nil
}

// params returns a map[string]string representation of VideoConfig.
func (config VideoConfig) params() (map[string]string, error) {
	params, _ := config.BaseFile.params()

	return params, nil
}

// name returns the field name for the Video.
func (config VideoConfig) name() string {
	return "video"
}

// method returns Telegram API method name for sending Video.
func (config VideoConfig) method() string {
	return "sendVideo"
}

// VoiceConfig contains information about a SendVoice request.
type VoiceConfig struct {
	BaseFile
	Duration int
}

// values returns a url.Values representation of VoiceConfig.
func (config VoiceConfig) values() (url.Values, error) {
	v, _ := config.BaseChat.values()

	v.Add(config.name(), config.FileID)
	if config.Duration != 0 {
		v.Add("duration", strconv.Itoa(config.Duration))
	}

	return v, nil
}

// params returns a map[string]string representation of VoiceConfig.
func (config VoiceConfig) params() (map[string]string, error) {
	params, _ := config.BaseFile.params()

	if config.Duration != 0 {
		params["duration"] = strconv.Itoa(config.Duration)
	}

	return params, nil
}

// name returns the field name for the Voice.
func (config VoiceConfig) name() string {
	return "voice"
}

// method returns Telegram API method name for sending Voice.
func (config VoiceConfig) method() string {
	return "sendVoice"
}

// LocationConfig contains information about a SendLocation request.
type LocationConfig struct {
	BaseChat
	Latitude  float64 // required
	Longitude float64 // required
}

// values returns a url.Values representation of LocationConfig.
func (config LocationConfig) values() (url.Values, error) {
	v, _ := config.BaseChat.values()

	v.Add("latitude", strconv.FormatFloat(config.Latitude, 'f', 6, 64))
	v.Add("longitude", strconv.FormatFloat(config.Longitude, 'f', 6, 64))

	return v, nil
}

// method returns Telegram API method name for sending Location.
func (config LocationConfig) method() string {
	return "sendLocation"
}

// ChatActionConfig contains information about a SendChatAction request.
type ChatActionConfig struct {
	BaseChat
	Action string // required
}

// values returns a url.Values representation of ChatActionConfig.
func (config ChatActionConfig) values() (url.Values, error) {
	v, _ := config.BaseChat.values()
	v.Add("action", config.Action)
	return v, nil
}

// method returns Telegram API method name for sending ChatAction.
func (config ChatActionConfig) method() string {
	return "sendChatAction"
}

// UserProfilePhotosConfig contains information about a
// GetUserProfilePhotos request.
type UserProfilePhotosConfig struct {
	UserID int
	Offset int
	Limit  int
}

// FileConfig has information about a file hosted on Telegram.
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

// FileBytes contains information about a set of bytes to upload
// as a File.
type FileBytes struct {
	Name  string
	Bytes []byte
}

// FileReader contains information about a reader to upload as a File.
// If Size is -1, it will read the entire Reader into memory to
// calculate a Size.
type FileReader struct {
	Name   string
	Reader io.Reader
	Size   int64
}

// InlineConfig contains information on making an InlineQuery response.
type InlineConfig struct {
	InlineQueryID string              `json:"inline_query_id"`
	Results       []InlineQueryResult `json:"results"`
	CacheTime     int                 `json:"cache_time"`
	IsPersonal    bool                `json:"is_personal"`
	NextOffset    string              `json:"next_offset"`
}
