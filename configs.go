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

// Base struct for all chat event(Message, Photo and so on)
type Chattable struct {
	ChatID          int
	ChannelUsername string
}

type Fileable struct {
	FilePath string
	File     interface{}
	FileID   string
}

func (chattable *Chattable) Values() (url.Values, error) {
	v := url.Values{}
	if chattable.ChannelUsername != "" {
		v.Add("chat_id", chattable.ChannelUsername)
	} else {
		v.Add("chat_id", strconv.Itoa(chattable.ChatID))
	}
	return v, nil
}

// MessageConfig contains information about a SendMessage request.
type MessageConfig struct {
	Chattable
	Text                  string
	ParseMode             string
	DisableWebPagePreview bool
	ReplyToMessageID      int
	ReplyMarkup           interface{}
}

func (config *MessageConfig) Values() (url.Values, error) {
	v, _ := config.Chattable.Values()
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

// ForwardConfig contains information about a ForwardMessage request.
type ForwardConfig struct {
	Chattable
	FromChatID          int
	FromChannelUsername string
	MessageID           int
}

func (config *ForwardConfig) Values() (url.Values, error) {
	v, _ := config.Chattable.Values()

	if config.FromChannelUsername != "" {
		v.Add("chat_id", config.FromChannelUsername)
	} else {
		v.Add("chat_id", strconv.Itoa(config.FromChatID))
	}
	v.Add("message_id", strconv.Itoa(config.MessageID))

	return v, nil
}

// PhotoConfig contains information about a SendPhoto request.
type PhotoConfig struct {
	Chattable
	Fileable
	Caption          string
	ReplyToMessageID int
	ReplyMarkup      interface{}
	UseExistingPhoto bool
}

func (config *PhotoConfig) Values() (url.Values, error) {
	v, _ := config.Chattable.Values()

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

// AudioConfig contains information about a SendAudio request.
type AudioConfig struct {
	Chattable
	Fileable
	Duration         int
	Performer        string
	Title            string
	ReplyToMessageID int
	ReplyMarkup      interface{}
	UseExistingAudio bool
}

func (config *AudioConfig) Values() (url.Values, error) {
	v, _ := config.Chattable.Values()

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

// DocumentConfig contains information about a SendDocument request.
type DocumentConfig struct {
	Chattable
	Fileable
	ReplyToMessageID    int
	ReplyMarkup         interface{}
	UseExistingDocument bool
}

func (config *DocumentConfig) Values() (url.Values, error) {
	v, _ := config.Chattable.Values()

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

// StickerConfig contains information about a SendSticker request.
type StickerConfig struct {
	Chattable
	Fileable
	ReplyToMessageID   int
	ReplyMarkup        interface{}
	UseExistingSticker bool
}

func (config *StickerConfig) Values() (url.Values, error) {
	v, _ := config.Chattable.Values()

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

// VideoConfig contains information about a SendVideo request.
type VideoConfig struct {
	Chattable
	Fileable
	Duration         int
	Caption          string
	ReplyToMessageID int
	ReplyMarkup      interface{}
	UseExistingVideo bool
}

func (config *VideoConfig) Values() (url.Values, error) {
	v, _ := config.Chattable.Values()

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

// VoiceConfig contains information about a SendVoice request.
type VoiceConfig struct {
	Chattable
	Fileable
	Duration         int
	ReplyToMessageID int
	ReplyMarkup      interface{}
	UseExistingVoice bool
}

func (config *VoiceConfig) Values() (url.Values, error) {
	v, _ := config.Chattable.Values()

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

// LocationConfig contains information about a SendLocation request.
type LocationConfig struct {
	Chattable
	Latitude         float64
	Longitude        float64
	ReplyToMessageID int
	ReplyMarkup      interface{}
}

func (config *LocationConfig) Values() (url.Values, error) {
	v, _ := config.Chattable.Values()

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

// ChatActionConfig contains information about a SendChatAction request.
type ChatActionConfig struct {
	Chattable
	Action string
}

func (config *ChatActionConfig) Values() (url.Values, error) {
	v, _ := config.Chattable.Values()
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
