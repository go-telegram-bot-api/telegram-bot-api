package tgbotapi

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

// APIResponse is a response from the Telegram API with the result
// stored raw.
type APIResponse struct {
	Ok          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	ErrorCode   int             `json:"error_code"`
	Description string          `json:"description"`
}

// Update is an update response, from GetUpdates.
type Update struct {
	UpdateID    int         `json:"update_id"`
	Message     Message     `json:"message"`
	InlineQuery InlineQuery `json:"inline_query"`
}

// User is a user on Telegram.
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"` // optional
	UserName  string `json:"username"`  // optional
}

// String displays a simple text version of a user.
//
// It is normally a user's username, but falls back to a first/last
// name as available.
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

// GroupChat is a group chat.
type GroupChat struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// Chat contains information about the place a message was sent.
type Chat struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`      // optional
	UserName  string `json:"username"`   // optional
	FirstName string `json:"first_name"` // optional
	LastName  string `json:"last_name"`  // optional
}

// IsPrivate returns if the Chat is a private conversation.
func (c *Chat) IsPrivate() bool {
	return c.Type == "private"
}

// IsGroup returns if the Chat is a group.
func (c *Chat) IsGroup() bool {
	return c.Type == "group"
}

// IsSuperGroup returns if the Chat is a supergroup.
func (c *Chat) IsSuperGroup() bool {
	return c.Type == "supergroup"
}

// IsChannel returns if the Chat is a channel.
func (c *Chat) IsChannel() bool {
	return c.Type == "channel"
}

// Message is returned by almost every request, and contains data about
// almost anything.
type Message struct {
	MessageID             int         `json:"message_id"`
	From                  User        `json:"from"` // optional
	Date                  int         `json:"date"`
	Chat                  Chat        `json:"chat"`
	ForwardFrom           User        `json:"forward_from"`            // optional
	ForwardDate           int         `json:"forward_date"`            // optional
	ReplyToMessage        *Message    `json:"reply_to_message"`        // optional
	Text                  string      `json:"text"`                    // optional
	Audio                 Audio       `json:"audio"`                   // optional
	Document              Document    `json:"document"`                // optional
	Photo                 []PhotoSize `json:"photo"`                   // optional
	Sticker               Sticker     `json:"sticker"`                 // optional
	Video                 Video       `json:"video"`                   // optional
	Voice                 Voice       `json:"voice"`                   // optional
	Caption               string      `json:"caption"`                 // optional
	Contact               Contact     `json:"contact"`                 // optional
	Location              Location    `json:"location"`                // optional
	NewChatParticipant    User        `json:"new_chat_participant"`    // optional
	LeftChatParticipant   User        `json:"left_chat_participant"`   // optional
	NewChatTitle          string      `json:"new_chat_title"`          // optional
	NewChatPhoto          []PhotoSize `json:"new_chat_photo"`          // optional
	DeleteChatPhoto       bool        `json:"delete_chat_photo"`       // optional
	GroupChatCreated      bool        `json:"group_chat_created"`      // optional
	SuperGroupChatCreated bool        `json:"supergroup_chat_created"` // optional
	ChannelChatCreated    bool        `json:"channel_chat_created"`    // optional
	MigrateToChatID       int         `json:"migrate_to_chat_id"`      // optional
	MigrateFromChatID     int         `json:"migrate_from_chat_id"`    // optional
}

// Time converts the message timestamp into a Time.
func (m *Message) Time() time.Time {
	return time.Unix(int64(m.Date), 0)
}

// IsGroup returns if the message was sent to a group.
//
// Deprecated in favor of Chat.IsGroup.
func (m *Message) IsGroup() bool {
	log.Println("Message.IsGroup is deprecated.")
	log.Println("Please use Chat.IsGroup instead.")
	return m.Chat.IsGroup()
}

// IsCommand returns true if message starts with '/'.
func (m *Message) IsCommand() bool {
	return m.Text != "" && m.Text[0] == '/'
}

// Command checks if the message was a command and if it was, returns the
// command. If the Message was not a command, it returns an empty string.
func (m *Message) Command() string {
	if !m.IsCommand() {
		return ""
	}

	return strings.SplitN(m.Text, " ", 2)[0]
}

// CommandArguments checks if the message was a command and if it was,
// returns all text after the command name. If the Message was not a
// command, it returns an empty string.
func (m *Message) CommandArguments() string {
	if !m.IsCommand() {
		return ""
	}

	split := strings.SplitN(m.Text, " ", 2)
	if len(split) != 2 {
		return ""
	}

	return strings.SplitN(m.Text, " ", 2)[1]
}

// PhotoSize contains information about photos.
type PhotoSize struct {
	FileID   string `json:"file_id"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	FileSize int    `json:"file_size"` // optional
}

// Audio contains information about audio.
type Audio struct {
	FileID    string `json:"file_id"`
	Duration  int    `json:"duration"`
	Performer string `json:"performer"` // optional
	Title     string `json:"title"`     // optional
	MimeType  string `json:"mime_type"` // optional
	FileSize  int    `json:"file_size"` // optional
}

// Document contains information about a document.
type Document struct {
	FileID    string    `json:"file_id"`
	Thumbnail PhotoSize `json:"thumb"`     // optional
	FileName  string    `json:"file_name"` // optional
	MimeType  string    `json:"mime_type"` // optional
	FileSize  int       `json:"file_size"` // optional
}

// Sticker contains information about a sticker.
type Sticker struct {
	FileID    string    `json:"file_id"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	Thumbnail PhotoSize `json:"thumb"`     // optional
	FileSize  int       `json:"file_size"` // optional
}

// Video contains information about a video.
type Video struct {
	FileID    string    `json:"file_id"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	Duration  int       `json:"duration"`
	Thumbnail PhotoSize `json:"thumb"`     // optional
	MimeType  string    `json:"mime_type"` // optional
	FileSize  int       `json:"file_size"` // optional
}

// Voice contains information about a voice.
type Voice struct {
	FileID   string `json:"file_id"`
	Duration int    `json:"duration"`
	MimeType string `json:"mime_type"` // optional
	FileSize int    `json:"file_size"` // optional
}

// Contact contains information about a contact.
//
// Note that LastName and UserID may be empty.
type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"` // optional
	UserID      int    `json:"user_id"`   // optional
}

// Location contains information about a place.
type Location struct {
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
}

// UserProfilePhotos contains information a set of user profile photos.
type UserProfilePhotos struct {
	TotalCount int         `json:"total_count"`
	Photos     []PhotoSize `json:"photos"`
}

// File contains information about a file to download from Telegram.
type File struct {
	FileID   string `json:"file_id"`
	FileSize int    `json:"file_size"` // optional
	FilePath string `json:"file_path"` // optional
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
	ResizeKeyboard  bool       `json:"resize_keyboard"`   // optional
	OneTimeKeyboard bool       `json:"one_time_keyboard"` // optional
	Selective       bool       `json:"selective"`         // optional
}

// ReplyKeyboardHide allows the Bot to hide a custom keyboard.
type ReplyKeyboardHide struct {
	HideKeyboard bool `json:"hide_keyboard"`
	Selective    bool `json:"selective"` // optional
}

// ForceReply allows the Bot to have users directly reply to it without
// additional interaction.
type ForceReply struct {
	ForceReply bool `json:"force_reply"`
	Selective  bool `json:"selective"` // optional
}

// InlineQuery is a Query from Telegram for an inline request.
type InlineQuery struct {
	ID     string `json:"id"`
	From   User   `json:"user"`
	Query  string `json:"query"`
	Offset string `json:"offset"`
}

// InlineQueryResult is the base type that all InlineQuery Results have.
type InlineQueryResult struct {
	Type string `json:"type"` // required
	ID   string `json:"id"`   // required
}

// InlineQueryResultArticle is an inline query response article.
type InlineQueryResultArticle struct {
	InlineQueryResult
	Title                 string `json:"title"`        // required
	MessageText           string `json:"message_text"` // required
	ParseMode             string `json:"parse_mode"`   // required
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
	URL                   string `json:"url"`
	HideURL               bool   `json:"hide_url"`
	Description           string `json:"description"`
	ThumbURL              string `json:"thumb_url"`
	ThumbWidth            int    `json:"thumb_width"`
	ThumbHeight           int    `json:"thumb_height"`
}

// InlineQueryResultPhoto is an inline query response photo.
type InlineQueryResultPhoto struct {
	InlineQueryResult
	URL                   string `json:"photo_url"` // required
	MimeType              string `json:"mime_type"`
	Width                 int    `json:"photo_width"`
	Height                int    `json:"photo_height"`
	ThumbURL              string `json:"thumb_url"`
	Title                 string `json:"title"`
	Description           string `json:"description"`
	Caption               string `json:"caption"`
	MessageText           string `json:"message_text"`
	ParseMode             string `json:"parse_mode"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
}

// InlineQueryResultGIF is an inline query response GIF.
type InlineQueryResultGIF struct {
	InlineQueryResult
	URL                   string `json:"gif_url"` // required
	Width                 int    `json:"gif_width"`
	Height                int    `json:"gif_height"`
	ThumbURL              string `json:"thumb_url"`
	Title                 string `json:"title"`
	Caption               string `json:"caption"`
	MessageText           string `json:"message_text"`
	ParseMode             string `json:"parse_mode"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
}

// InlineQueryResultMPEG4GIF is an inline query response MPEG4 GIF.
type InlineQueryResultMPEG4GIF struct {
	InlineQueryResult
	URL                   string `json:"mpeg4_url"` // required
	Width                 int    `json:"mpeg4_width"`
	Height                int    `json:"mpeg4_height"`
	ThumbURL              string `json:"thumb_url"`
	Title                 string `json:"title"`
	Caption               string `json:"caption"`
	MessageText           string `json:"message_text"`
	ParseMode             string `json:"parse_mode"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
}

// InlineQueryResultVideo is an inline query response video.
type InlineQueryResultVideo struct {
	InlineQueryResult
	URL                   string `json:"video_url"`    // required
	MimeType              string `json:"mime_type"`    // required
	MessageText           string `json:"message_text"` // required
	ParseMode             string `json:"parse_mode"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
	Width                 int    `json:"video_width"`
	Height                int    `json:"video_height"`
	ThumbURL              string `json:"thumb_url"`
	Title                 string `json:"title"`
	Description           string `json:"description"`
}
