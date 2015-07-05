package tgbotapi

import (
	"encoding/json"
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

// APIResponse is a response from the Telegram API with the result stored raw.
type APIResponse struct {
	Ok          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	ErrorCode   int             `json:"error_code"`
	Description string          `json:"description"`
}

// Update is an update response, from GetUpdates.
type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

// User is a user, contained in Message and returned by GetSelf.
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
}

// GroupChat is a group chat, and not currently in use.
type GroupChat struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// UserOrGroupChat is returned in Message, because it's not clear which it is.
type UserOrGroupChat struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
	Title     string `json:"title"`
}

// Message is returned by almost every request, and contains data about almost anything.
type Message struct {
	MessageID           int             `json:"message_id"`
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

// PhotoSize contains information about photos, including ID and Width and Height.
type PhotoSize struct {
	FileID   string `json:"file_id"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	FileSize int    `json:"file_size"`
}

// Audio contains information about audio, including ID and Duration.
type Audio struct {
	FileID   string `json:"file_id"`
	Duration int    `json:"duration"`
	MimeType string `json:"mime_type"`
	FileSize int    `json:"file_size"`
}

// Document contains information about a document, including ID and a Thumbnail.
type Document struct {
	FileID    string    `json:"file_id"`
	Thumbnail PhotoSize `json:"thumb"`
	FileName  string    `json:"file_name"`
	MimeType  string    `json:"mime_type"`
	FileSize  int       `json:"file_size"`
}

// Sticker contains information about a sticker, including ID and Thumbnail.
type Sticker struct {
	FileID    string    `json:"file_id"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	Thumbnail PhotoSize `json:"thumb"`
	FileSize  int       `json:"file_size"`
}

// Video contains information about a video, including ID and duration and Thumbnail.
type Video struct {
	FileID    string    `json:"file_id"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	Duration  int       `json:"duration"`
	Thumbnail PhotoSize `json:"thumb"`
	MimeType  string    `json:"mime_type"`
	FileSize  int       `json:"file_size"`
	Caption   string    `json:"caption"`
}

// Contact contains information about a contact, such as PhoneNumber and UserId.
type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	UserID      string `json:"user_id"`
}

// Location contains information about a place, such as Longitude and Latitude.
type Location struct {
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
}

// UserProfilePhotos contains information a set of user profile photos.
type UserProfilePhotos struct {
	TotalCount int         `json:"total_count"`
	Photos     []PhotoSize `json:"photos"`
}

// ReplyKeyboardMarkup allows the Bot to set a custom keyboard.
type ReplyKeyboardMarkup struct {
	Keyboard        [][]string `json:"keyboard"`
	ResizeKeyboard  bool       `json:"resize_keyboard"`
	OneTimeKeyboard bool       `json:"one_time_keyboard"`
	Selective       bool       `json:"selective"`
}

// ReplyKeyboardHide allows the Bot to hide a custom keyboard.
type ReplyKeyboardHide struct {
	HideKeyboard bool `json:"hide_keyboard"`
	Selective    bool `json:"selective"`
}

// ForceReply allows the Bot to have users directly reply to it without additional interaction.
type ForceReply struct {
	ForceReply bool `json:"force_reply"`
	Selective  bool `json:"force_reply"`
}

// MessageConfig contains information about a SendMessage request.
type MessageConfig struct {
	ChatID                int
	Text                  string
	DisableWebPagePreview bool
	ReplyToMessageID      int
	ReplyMarkup           interface{}
}

// ForwardConfig contains infomation about a ForwardMessage request.
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
