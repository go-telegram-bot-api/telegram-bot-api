package tgbotapi

import (
	"io"
	"net/url"
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

// MessageConfig contains information about a SendMessage request.
type MessageConfig struct {
	Chattable
	Text                  string
	ParseMode             string
	DisableWebPagePreview bool
	ReplyToMessageID      int
	ReplyMarkup           interface{}
}

// ForwardConfig contains information about a ForwardMessage request.
type ForwardConfig struct {
	Chattable
	FromChatID          int
	FromChannelUsername string
	MessageID           int
}

// PhotoConfig contains information about a SendPhoto request.
type PhotoConfig struct {
	Chattable
	Caption          string
	ReplyToMessageID int
	ReplyMarkup      interface{}
	UseExistingPhoto bool
	FilePath         string
	File             interface{}
	FileID           string
}

// AudioConfig contains information about a SendAudio request.
type AudioConfig struct {
	Chattable
	Duration         int
	Performer        string
	Title            string
	ReplyToMessageID int
	ReplyMarkup      interface{}
	UseExistingAudio bool
	FilePath         string
	File             interface{}
	FileID           string
}

// DocumentConfig contains information about a SendDocument request.
type DocumentConfig struct {
	Chattable
	ReplyToMessageID    int
	ReplyMarkup         interface{}
	UseExistingDocument bool
	FilePath            string
	File                interface{}
	FileID              string
}

// StickerConfig contains information about a SendSticker request.
type StickerConfig struct {
	Chattable
	ReplyToMessageID   int
	ReplyMarkup        interface{}
	UseExistingSticker bool
	FilePath           string
	File               interface{}
	FileID             string
}

// VideoConfig contains information about a SendVideo request.
type VideoConfig struct {
	Chattable
	Duration         int
	Caption          string
	ReplyToMessageID int
	ReplyMarkup      interface{}
	UseExistingVideo bool
	FilePath         string
	File             interface{}
	FileID           string
}

// VoiceConfig contains information about a SendVoice request.
type VoiceConfig struct {
	Chattable
	Duration         int
	ReplyToMessageID int
	ReplyMarkup      interface{}
	UseExistingVoice bool
	FilePath         string
	File             interface{}
	FileID           string
}

// LocationConfig contains information about a SendLocation request.
type LocationConfig struct {
	Chattable
	Latitude         float64
	Longitude        float64
	ReplyToMessageID int
	ReplyMarkup      interface{}
}

// ChatActionConfig contains information about a SendChatAction request.
type ChatActionConfig struct {
	Chattable
	Action string
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
