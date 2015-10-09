package tgbotapi

import (
	"encoding/json"
	"fmt"
	"time"
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

// String displays a simple text version of a user.
// It is normally a user's username,
// but falls back to a first/last name as available.
func (u *User) String() string {
	if u.UserName != "" {
		return u.UserName
	}

	name := u.FirstName
	if u.LastName != "" {
		name += " " + u.LastName
	}

	return name
}

// GroupChat is a group chat, and not currently in use.
type GroupChat struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// Chat is returned in Message, it contains information about the Chat a message was sent in.
type Chat struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	UserName  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// IsPrivate returns true if the Chat is a private conversation
func (c *Chat) IsPrivate() bool {
	return c.Type == "private"
}

// IsGroup returns true if the Chat is a group conversation
func (c *Chat) IsGroup() bool {
	return c.Type == "group"
}

// IsChannel returns true if the Chat is a channel
func (c *Chat) IsChannel() bool {
	return c.Type == "channel"
}

// Message is returned by almost every request, and contains data about almost anything.
type Message struct {
	MessageID           int         `json:"message_id"`
	From                User        `json:"from"`
	Date                int         `json:"date"`
	Chat                Chat        `json:"chat"`
	ForwardFrom         User        `json:"forward_from"`
	ForwardDate         int         `json:"forward_date"`
	ReplyToMessage      *Message    `json:"reply_to_message"`
	Text                string      `json:"text"`
	Audio               Audio       `json:"audio"`
	Document            Document    `json:"document"`
	Photo               []PhotoSize `json:"photo"`
	Sticker             Sticker     `json:"sticker"`
	Video               Video       `json:"video"`
	Voice               Voice       `json:"voice"`
	Caption             string      `json:"caption"`
	Contact             Contact     `json:"contact"`
	Location            Location    `json:"location"`
	NewChatParticipant  User        `json:"new_chat_participant"`
	LeftChatParticipant User        `json:"left_chat_participant"`
	NewChatTitle        string      `json:"new_chat_title"`
	NewChatPhoto        []PhotoSize `json:"new_chat_photo"`
	DeleteChatPhoto     bool        `json:"delete_chat_photo"`
	GroupChatCreated    bool        `json:"group_chat_created"`
}

// Time converts the message timestamp into a Time.
func (m *Message) Time() time.Time {
	return time.Unix(int64(m.Date), 0)
}

// IsGroup returns if the message was sent to a group.
func (m *Message) IsGroup() bool {
	return m.From.ID != m.Chat.ID
}

// PhotoSize contains information about photos, including ID and Width and Height.
type PhotoSize struct {
	FileID   string `json:"file_id"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	FileSize int    `json:"file_size"`
}

// Audio contains information about audio,
// including ID, Duration, Performer and Title.
type Audio struct {
	FileID    string `json:"file_id"`
	Duration  int    `json:"duration"`
	Performer string `json:"performer"`
	Title     string `json:"title"`
	MimeType  string `json:"mime_type"`
	FileSize  int    `json:"file_size"`
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
}

// Voice contains information about a voice, including ID and duration.
type Voice struct {
	FileID   string `json:"file_id"`
	Duration int    `json:"duration"`
	MimeType string `json:"mime_type"`
	FileSize int    `json:"file_size"`
}

// Contact contains information about a contact, such as PhoneNumber and UserId.
type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	UserID      int    `json:"user_id"`
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

// File contains information about a file to download from Telegram
type File struct {
	FileID   string `json:"file_id"`
	FileSize int    `json:"file_size"`
	FilePath string `json:"file_path"`
}

// Link returns a full path to the download URL for a File.
//
// It requires the Bot Token to create the link.
func (f *File) Link(token string) string {
	return fmt.Sprintf(FileEndpoint, token, f.FilePath)
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
	Selective  bool `json:"selective"`
}
