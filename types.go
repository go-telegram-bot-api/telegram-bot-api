package tgbotapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// APIResponse is a response from the Telegram API with the result
// stored raw.
type APIResponse struct {
	Ok          bool                `json:"ok"`
	Result      json.RawMessage     `json:"result"`
	ErrorCode   int                 `json:"error_code"`
	Description string              `json:"description"`
	Parameters  *ResponseParameters `json:"parameters"`
}

// ResponseParameters are various errors that can be returned in APIResponse.
type ResponseParameters struct {
	MigrateToChatID int64 `json:"migrate_to_chat_id"` // optional
	RetryAfter      int   `json:"retry_after"`        // optional
}

// Update is an update response, from GetUpdates.
type Update struct {
	UpdateID           int                 `json:"update_id"`
	Message            *Message            `json:"message"`
	EditedMessage      *Message            `json:"edited_message"`
	ChannelPost        *Message            `json:"channel_post"`
	EditedChannelPost  *Message            `json:"edited_channel_post"`
	InlineQuery        *InlineQuery        `json:"inline_query"`
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result"`
	CallbackQuery      *CallbackQuery      `json:"callback_query"`
	ShippingQuery      *ShippingQuery      `json:"shipping_query"`
	PreCheckoutQuery   *PreCheckoutQuery   `json:"pre_checkout_query"`
}

// UpdatesChannel is the channel for getting updates.
type UpdatesChannel <-chan Update

// Clear discards all unprocessed incoming updates.
func (ch UpdatesChannel) Clear() {
	for len(ch) != 0 {
		<-ch
	}
}

// User is a user on Telegram.
type User struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`     // optional
	UserName     string `json:"username"`      // optional
	LanguageCode string `json:"language_code"` // optional
	IsBot        bool   `json:"is_bot"`        // optional
}

// String displays a simple text version of a user.
//
// It is normally a user's username, but falls back to a first/last
// name as available.
func (u *User) String() string {
	if u == nil {
		return ""
	}
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

// ChatPhoto represents a chat photo.
type ChatPhoto struct {
	SmallFileID string `json:"small_file_id"`
	BigFileID   string `json:"big_file_id"`
}

// Chat contains information about the place a message was sent.
type Chat struct {
	ID                  int64      `json:"id"`
	Type                string     `json:"type"`
	Title               string     `json:"title"`                          // optional
	UserName            string     `json:"username"`                       // optional
	FirstName           string     `json:"first_name"`                     // optional
	LastName            string     `json:"last_name"`                      // optional
	AllMembersAreAdmins bool       `json:"all_members_are_administrators"` // optional
	Photo               *ChatPhoto `json:"photo"`
	Description         string     `json:"description,omitempty"` // optional
	InviteLink          string     `json:"invite_link,omitempty"` // optional
	PinnedMessage       *Message   `json:"pinned_message"`        // optional
}

// IsPrivate returns if the Chat is a private conversation.
func (c Chat) IsPrivate() bool {
	return c.Type == "private"
}

// IsGroup returns if the Chat is a group.
func (c Chat) IsGroup() bool {
	return c.Type == "group"
}

// IsSuperGroup returns if the Chat is a supergroup.
func (c Chat) IsSuperGroup() bool {
	return c.Type == "supergroup"
}

// IsChannel returns if the Chat is a channel.
func (c Chat) IsChannel() bool {
	return c.Type == "channel"
}

// ChatConfig returns a ChatConfig struct for chat related methods.
func (c Chat) ChatConfig() ChatConfig {
	return ChatConfig{ChatID: c.ID}
}

// Message is returned by almost every request, and contains data about
// almost anything.
type Message struct {
	// MessageID is a unique message identifier inside this chat
	MessageID int `json:"message_id"`
	// From is a sender, empty for messages sent to channels;
	// optional
	From *User `json:"from"`
	// Date of the message was sent in Unix time
	Date int `json:"date"`
	// Chat is the conversation the message belongs to
	Chat *Chat `json:"chat"`
	// ForwardFrom for forwarded messages, sender of the original message;
	// optional
	ForwardFrom *User `json:"forward_from"`
	// ForwardFromChat for messages forwarded from channels,
	// information about the original channel;
	// optional
	ForwardFromChat *Chat `json:"forward_from_chat"`
	// ForwardFromMessageID for messages forwarded from channels,
	// identifier of the original message in the channel;
	// optional
	ForwardFromMessageID int `json:"forward_from_message_id"`
	// ForwardDate for forwarded messages, date the original message was sent in Unix time;
	// optional
	ForwardDate int `json:"forward_date"`
	// ReplyToMessage for replies, the original message.
	// Note that the Message object in this field will not contain further ReplyToMessage fields
	// even if it itself is a reply;
	// optional
	ReplyToMessage *Message `json:"reply_to_message"`
	// ViaBot through which the message was sent;
	// optional
	ViaBot *User `json:"via_bot"`
	// EditDate of the message was last edited in Unix time;
	// optional
	EditDate int `json:"edit_date"`
	// MediaGroupID is the unique identifier of a media message group this message belongs to;
	// optional
	MediaGroupID string `json:"media_group_id"`
	// AuthorSignature is the signature of the post author for messages in channels;
	// optional
	AuthorSignature string `json:"author_signature"`
	// Text is for text messages, the actual UTF-8 text of the message, 0-4096 characters;
	// optional
	Text string `json:"text"`
	// Entities is for text messages, special entities like usernames,
	// URLs, bot commands, etc. that appear in the text;
	// optional
	Entities *[]MessageEntity `json:"entities"`
	// CaptionEntities;
	// optional
	CaptionEntities *[]MessageEntity `json:"caption_entities"`
	// Audio message is an audio file, information about the file;
	// optional
	Audio *Audio `json:"audio"`
	// Document message is a general file, information about the file;
	// optional
	Document *Document `json:"document"`
	// Animation message is an animation, information about the animation.
	// For backward compatibility, when this field is set, the document field will also be set;
	// optional
	Animation *ChatAnimation `json:"animation"`
	// Game message is a game, information about the game;
	// optional
	Game *Game `json:"game"`
	// Photo message is a photo, available sizes of the photo;
	// optional
	Photo *[]PhotoSize `json:"photo"`
	// Sticker message is a sticker, information about the sticker;
	// optional
	Sticker *Sticker `json:"sticker"`
	// Video message is a video, information about the video;
	// optional
	Video *Video `json:"video"`
	// VideoNote message is a video note, information about the video message;
	// optional
	VideoNote *VideoNote `json:"video_note"`
	// Voice message is a voice message, information about the file;
	// optional
	Voice *Voice `json:"voice"`
	// Caption for the animation, audio, document, photo, video or voice, 0-1024 characters;
	// optional
	Caption string `json:"caption"`
	// Contact message is a shared contact, information about the contact;
	// optional
	Contact *Contact `json:"contact"`
	// Location message is a shared location, information about the location;
	// optional
	Location *Location `json:"location"`
	// Venue message is a venue, information about the venue.
	// For backward compatibility, when this field is set, the location field will also be set;
	// optional
	Venue *Venue `json:"venue"`
	// NewChatMembers that were added to the group or supergroup
	// and information about them (the bot itself may be one of these members);
	// optional
	NewChatMembers *[]User `json:"new_chat_members"`
	// LeftChatMember is a member was removed from the group,
	// information about them (this member may be the bot itself);
	// optional
	LeftChatMember *User `json:"left_chat_member"`
	// NewChatTitle is a chat title was changed to this value;
	// optional
	NewChatTitle string `json:"new_chat_title"`
	// NewChatPhoto is a chat photo was change to this value;
	// optional
	NewChatPhoto *[]PhotoSize `json:"new_chat_photo"`
	// DeleteChatPhoto is a service message: the chat photo was deleted;
	// optional
	DeleteChatPhoto bool `json:"delete_chat_photo"`
	// GroupChatCreated is a service message: the group has been created;
	// optional
	GroupChatCreated bool `json:"group_chat_created"`
	// SuperGroupChatCreated is a service message: the supergroup has been created.
	// This field can't be received in a message coming through updates,
	// because bot can't be a member of a supergroup when it is created.
	// It can only be found in ReplyToMessage if someone replies to a very first message
	// in a directly created supergroup;
	// optional
	SuperGroupChatCreated bool `json:"supergroup_chat_created"`
	// ChannelChatCreated is a service message: the channel has been created.
	// This field can't be received in a message coming through updates,
	// because bot can't be a member of a channel when it is created.
	// It can only be found in ReplyToMessage
	// if someone replies to a very first message in a channel;
	// optional
	ChannelChatCreated bool `json:"channel_chat_created"`
	// MigrateToChatID is the group has been migrated to a supergroup with the specified identifier.
	// This number may be greater than 32 bits and some programming languages
	// may have difficulty/silent defects in interpreting it.
	// But it is smaller than 52 bits, so a signed 64 bit integer
	// or double-precision float type are safe for storing this identifier;
	// optional
	MigrateToChatID int64 `json:"migrate_to_chat_id"`
	// MigrateFromChatID is the supergroup has been migrated from a group with the specified identifier.
	// This number may be greater than 32 bits and some programming languages
	// may have difficulty/silent defects in interpreting it.
	// But it is smaller than 52 bits, so a signed 64 bit integer
	// or double-precision float type are safe for storing this identifier;
	// optional
	MigrateFromChatID int64 `json:"migrate_from_chat_id"`
	// PinnedMessage is a specified message was pinned.
	// Note that the Message object in this field will not contain further ReplyToMessage
	// fields even if it is itself a reply;
	// optional
	PinnedMessage *Message `json:"pinned_message"`
	// Invoice message is an invoice for a payment;
	// optional
	Invoice *Invoice `json:"invoice"`
	// SuccessfulPayment message is a service message about a successful payment,
	// information about the payment;
	// optional
	SuccessfulPayment *SuccessfulPayment `json:"successful_payment"`
	// PassportData is a Telegram Passport data;
	// optional
	PassportData *PassportData `json:"passport_data,omitempty"`
}

// Time converts the message timestamp into a Time.
func (m *Message) Time() time.Time {
	return time.Unix(int64(m.Date), 0)
}

// IsCommand returns true if message starts with a "bot_command" entity.
func (m *Message) IsCommand() bool {
	if m.Entities == nil || len(*m.Entities) == 0 {
		return false
	}

	entity := (*m.Entities)[0]
	return entity.Offset == 0 && entity.IsCommand()
}

// Command checks if the message was a command and if it was, returns the
// command. If the Message was not a command, it returns an empty string.
//
// If the command contains the at name syntax, it is removed. Use
// CommandWithAt() if you do not want that.
func (m *Message) Command() string {
	command := m.CommandWithAt()

	if i := strings.Index(command, "@"); i != -1 {
		command = command[:i]
	}

	return command
}

// CommandWithAt checks if the message was a command and if it was, returns the
// command. If the Message was not a command, it returns an empty string.
//
// If the command contains the at name syntax, it is not removed. Use Command()
// if you want that.
func (m *Message) CommandWithAt() string {
	if !m.IsCommand() {
		return ""
	}

	// IsCommand() checks that the message begins with a bot_command entity
	entity := (*m.Entities)[0]
	return m.Text[1:entity.Length]
}

// CommandArguments checks if the message was a command and if it was,
// returns all text after the command name. If the Message was not a
// command, it returns an empty string.
//
// Note: The first character after the command name is omitted:
// - "/foo bar baz" yields "bar baz", not " bar baz"
// - "/foo-bar baz" yields "bar baz", too
// Even though the latter is not a command conforming to the spec, the API
// marks "/foo" as command entity.
func (m *Message) CommandArguments() string {
	if !m.IsCommand() {
		return ""
	}

	// IsCommand() checks that the message begins with a bot_command entity
	entity := (*m.Entities)[0]
	if len(m.Text) == entity.Length {
		return "" // The command makes up the whole message
	}

	return m.Text[entity.Length+1:]
}

// MessageEntity contains information about data in a Message.
type MessageEntity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	URL    string `json:"url"`  // optional
	User   *User  `json:"user"` // optional
}

// ParseURL attempts to parse a URL contained within a MessageEntity.
func (e MessageEntity) ParseURL() (*url.URL, error) {
	if e.URL == "" {
		return nil, errors.New(ErrBadURL)
	}

	return url.Parse(e.URL)
}

// IsMention returns true if the type of the message entity is "mention" (@username).
func (e MessageEntity) IsMention() bool {
	return e.Type == "mention"
}

// IsHashtag returns true if the type of the message entity is "hashtag".
func (e MessageEntity) IsHashtag() bool {
	return e.Type == "hashtag"
}

// IsCommand returns true if the type of the message entity is "bot_command".
func (e MessageEntity) IsCommand() bool {
	return e.Type == "bot_command"
}

// IsUrl returns true if the type of the message entity is "url".
func (e MessageEntity) IsUrl() bool {
	return e.Type == "url"
}

// IsEmail returns true if the type of the message entity is "email".
func (e MessageEntity) IsEmail() bool {
	return e.Type == "email"
}

// IsBold returns true if the type of the message entity is "bold" (bold text).
func (e MessageEntity) IsBold() bool {
	return e.Type == "bold"
}

// IsItalic returns true if the type of the message entity is "italic" (italic text).
func (e MessageEntity) IsItalic() bool {
	return e.Type == "italic"
}

// IsCode returns true if the type of the message entity is "code" (monowidth string).
func (e MessageEntity) IsCode() bool {
	return e.Type == "code"
}

// IsPre returns true if the type of the message entity is "pre" (monowidth block).
func (e MessageEntity) IsPre() bool {
	return e.Type == "pre"
}

// IsTextLink returns true if the type of the message entity is "text_link" (clickable text URL).
func (e MessageEntity) IsTextLink() bool {
	return e.Type == "text_link"
}

// PhotoSize contains information about photos.
type PhotoSize struct {
	// FileID identifier for this file, which can be used to download or reuse the file
	FileID string `json:"file_id"`
	// Width photo width
	Width int `json:"width"`
	// Height photo height
	Height int `json:"height"`
	// FileSize file size
	// optional
	FileSize int `json:"file_size"`
}

// Audio contains information about audio.
type Audio struct {
	// FileID is an identifier for this file, which can be used to download or reuse the file
	FileID string `json:"file_id"`
	// Duration of the audio in seconds as defined by sender
	Duration int `json:"duration"`
	// Performer of the audio as defined by sender or by audio tags
	// optional
	Performer string `json:"performer"`
	// Title of the audio as defined by sender or by audio tags
	// optional
	Title string `json:"title"`
	// MimeType of the file as defined by sender
	// optional
	MimeType string `json:"mime_type"`
	// FileSize file size
	// optional
	FileSize int `json:"file_size"`
}

// Document contains information about a document.
type Document struct {
	// FileID is a identifier for this file, which can be used to download or reuse the file
	FileID string `json:"file_id"`
	// Thumbnail document thumbnail as defined by sender
	// optional
	Thumbnail *PhotoSize `json:"thumb"`
	// FileName original filename as defined by sender
	// optional
	FileName string `json:"file_name"`
	// MimeType  of the file as defined by sender
	// optional
	MimeType string `json:"mime_type"`
	// FileSize file size
	// optional
	FileSize int `json:"file_size"`
}

// Sticker contains information about a sticker.
type Sticker struct {
	// FileUniqueID is an unique identifier for this file,
	// which is supposed to be the same over time and for different bots.
	// Can't be used to download or reuse the file.
	FileUniqueID string `json:"file_unique_id"`
	// FileID is an identifier for this file, which can be used to download or reuse the file
	FileID string `json:"file_id"`
	// Width sticker width
	Width int `json:"width"`
	// Height sticker height
	Height int `json:"height"`
	// Thumbnail sticker thumbnail in the .WEBP or .JPG format
	// optional
	Thumbnail *PhotoSize `json:"thumb"`
	// Emoji associated with the sticker
	// optional
	Emoji string `json:"emoji"`
	// FileSize
	// optional
	FileSize int `json:"file_size"`
	// SetName of the sticker set to which the sticker belongs
	// optional
	SetName string `json:"set_name"`
	// IsAnimated true, if the sticker is animated
	// optional
	IsAnimated bool `json:"is_animated"`
}

// StickerSet contains information about an sticker set.
type StickerSet struct {
	// Name sticker set name
	Name string `json:"name"`
	// Title sticker set title
	Title string `json:"title"`
	// IsAnimated true, if the sticker set contains animated stickers
	IsAnimated bool `json:"is_animated"`
	// ContainsMasks true, if the sticker set contains masks
	ContainsMasks bool `json:"contains_masks"`
	// Stickers list of all set stickers
	Stickers []Sticker `json:"stickers"`
}

// ChatAnimation contains information about an animation.
type ChatAnimation struct {
	// FileID odentifier for this file, which can be used to download or reuse the file
	FileID string `json:"file_id"`
	// Width video width as defined by sender
	Width int `json:"width"`
	// Height video height as defined by sender
	Height int `json:"height"`
	// Duration of the video in seconds as defined by sender
	Duration int `json:"duration"`
	// Thumbnail animation thumbnail as defined by sender
	// optional
	Thumbnail *PhotoSize `json:"thumb"`
	// FileName original animation filename as defined by sender
	// optional
	FileName string `json:"file_name"`
	// MimeType of the file as defined by sender
	// optional
	MimeType string `json:"mime_type"`
	// FileSize file size
	// optional
	FileSize int `json:"file_size"`
}

// Video contains information about a video.
type Video struct {
	// FileID identifier for this file, which can be used to download or reuse the file
	FileID string `json:"file_id"`
	// Width video width as defined by sender
	Width int `json:"width"`
	// Height video height as defined by sender
	Height int `json:"height"`
	// Duration of the video in seconds as defined by sender
	Duration int `json:"duration"`
	// Thumbnail video thumbnail
	// optional
	Thumbnail *PhotoSize `json:"thumb"`
	// MimeType of a file as defined by sender
	// optional
	MimeType string `json:"mime_type"`
	// FileSize file size
	// optional
	FileSize int `json:"file_size"`
}

// VideoNote contains information about a video.
type VideoNote struct {
	// FileID identifier for this file, which can be used to download or reuse the file
	FileID string `json:"file_id"`
	// Length video width and height (diameter of the video message) as defined by sender
	Length int `json:"length"`
	// Duration of the video in seconds as defined by sender
	Duration int `json:"duration"`
	// Thumbnail video thumbnail
	// optional
	Thumbnail *PhotoSize `json:"thumb"`
	// FileSize file size
	// optional
	FileSize int `json:"file_size"`
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
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

// Venue contains information about a venue, including its Location.
type Venue struct {
	Location     Location `json:"location"`
	Title        string   `json:"title"`
	Address      string   `json:"address"`
	FoursquareID string   `json:"foursquare_id"` // optional
}

// UserProfilePhotos contains a set of user profile photos.
type UserProfilePhotos struct {
	TotalCount int           `json:"total_count"`
	Photos     [][]PhotoSize `json:"photos"`
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
	Keyboard        [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool               `json:"resize_keyboard"`   // optional
	OneTimeKeyboard bool               `json:"one_time_keyboard"` // optional
	Selective       bool               `json:"selective"`         // optional
}

// KeyboardButton is a button within a custom keyboard.
type KeyboardButton struct {
	Text            string `json:"text"`
	RequestContact  bool   `json:"request_contact"`
	RequestLocation bool   `json:"request_location"`
}

// ReplyKeyboardHide allows the Bot to hide a custom keyboard.
type ReplyKeyboardHide struct {
	HideKeyboard bool `json:"hide_keyboard"`
	Selective    bool `json:"selective"` // optional
}

// ReplyKeyboardRemove allows the Bot to hide a custom keyboard.
type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`
	Selective      bool `json:"selective"`
}

// InlineKeyboardMarkup is a custom keyboard presented for an inline bot.
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

// InlineKeyboardButton is a button within a custom keyboard for
// inline query responses.
//
// Note that some values are references as even an empty string
// will change behavior.
//
// CallbackGame, if set, MUST be first button in first row.
type InlineKeyboardButton struct {
	Text                         string        `json:"text"`
	URL                          *string       `json:"url,omitempty"`                              // optional
	CallbackData                 *string       `json:"callback_data,omitempty"`                    // optional
	SwitchInlineQuery            *string       `json:"switch_inline_query,omitempty"`              // optional
	SwitchInlineQueryCurrentChat *string       `json:"switch_inline_query_current_chat,omitempty"` // optional
	CallbackGame                 *CallbackGame `json:"callback_game,omitempty"`                    // optional
	Pay                          bool          `json:"pay,omitempty"`                              // optional
}

// CallbackQuery is data sent when a keyboard button with callback data
// is clicked.
type CallbackQuery struct {
	ID              string   `json:"id"`
	From            *User    `json:"from"`
	Message         *Message `json:"message"`           // optional
	InlineMessageID string   `json:"inline_message_id"` // optional
	ChatInstance    string   `json:"chat_instance"`
	Data            string   `json:"data"`            // optional
	GameShortName   string   `json:"game_short_name"` // optional
}

// ForceReply allows the Bot to have users directly reply to it without
// additional interaction.
type ForceReply struct {
	ForceReply bool `json:"force_reply"`
	Selective  bool `json:"selective"` // optional
}

// ChatMember is information about a member in a chat.
type ChatMember struct {
	User                  *User  `json:"user"`
	Status                string `json:"status"`
	CustomTitle           string `json:"custom_title,omitempty"`              // optional
	UntilDate             int64  `json:"until_date,omitempty"`                // optional
	CanBeEdited           bool   `json:"can_be_edited,omitempty"`             // optional
	CanChangeInfo         bool   `json:"can_change_info,omitempty"`           // optional
	CanPostMessages       bool   `json:"can_post_messages,omitempty"`         // optional
	CanEditMessages       bool   `json:"can_edit_messages,omitempty"`         // optional
	CanDeleteMessages     bool   `json:"can_delete_messages,omitempty"`       // optional
	CanInviteUsers        bool   `json:"can_invite_users,omitempty"`          // optional
	CanRestrictMembers    bool   `json:"can_restrict_members,omitempty"`      // optional
	CanPinMessages        bool   `json:"can_pin_messages,omitempty"`          // optional
	CanPromoteMembers     bool   `json:"can_promote_members,omitempty"`       // optional
	CanSendMessages       bool   `json:"can_send_messages,omitempty"`         // optional
	CanSendMediaMessages  bool   `json:"can_send_media_messages,omitempty"`   // optional
	CanSendOtherMessages  bool   `json:"can_send_other_messages,omitempty"`   // optional
	CanAddWebPagePreviews bool   `json:"can_add_web_page_previews,omitempty"` // optional
}

// IsCreator returns if the ChatMember was the creator of the chat.
func (chat ChatMember) IsCreator() bool { return chat.Status == "creator" }

// IsAdministrator returns if the ChatMember is a chat administrator.
func (chat ChatMember) IsAdministrator() bool { return chat.Status == "administrator" }

// IsMember returns if the ChatMember is a current member of the chat.
func (chat ChatMember) IsMember() bool { return chat.Status == "member" }

// HasLeft returns if the ChatMember left the chat.
func (chat ChatMember) HasLeft() bool { return chat.Status == "left" }

// WasKicked returns if the ChatMember was kicked from the chat.
func (chat ChatMember) WasKicked() bool { return chat.Status == "kicked" }

// Game is a game within Telegram.
type Game struct {
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	Photo        []PhotoSize     `json:"photo"`
	Text         string          `json:"text"`
	TextEntities []MessageEntity `json:"text_entities"`
	Animation    Animation       `json:"animation"`
}

// Animation is a GIF animation demonstrating the game.
type Animation struct {
	FileID   string    `json:"file_id"`
	Thumb    PhotoSize `json:"thumb"`
	FileName string    `json:"file_name"`
	MimeType string    `json:"mime_type"`
	FileSize int       `json:"file_size"`
}

// GameHighScore is a user's score and position on the leaderboard.
type GameHighScore struct {
	Position int  `json:"position"`
	User     User `json:"user"`
	Score    int  `json:"score"`
}

// CallbackGame is for starting a game in an inline keyboard button.
type CallbackGame struct{}

// WebhookInfo is information about a currently set webhook.
type WebhookInfo struct {
	URL                  string `json:"url"`
	HasCustomCertificate bool   `json:"has_custom_certificate"`
	PendingUpdateCount   int    `json:"pending_update_count"`
	LastErrorDate        int    `json:"last_error_date"`    // optional
	LastErrorMessage     string `json:"last_error_message"` // optional
	MaxConnections       int    `json:"max_connections"`    // optional
}

// IsSet returns true if a webhook is currently set.
func (info WebhookInfo) IsSet() bool {
	return info.URL != ""
}

// InputMediaPhoto contains a photo for displaying as part of a media group.
type InputMediaPhoto struct {
	Type      string `json:"type"`
	Media     string `json:"media"`
	Caption   string `json:"caption"`
	ParseMode string `json:"parse_mode"`
}

// InputMediaVideo contains a video for displaying as part of a media group.
type InputMediaVideo struct {
	Type  string `json:"type"`
	Media string `json:"media"`
	// thumb intentionally missing as it is not currently compatible
	Caption           string `json:"caption"`
	ParseMode         string `json:"parse_mode"`
	Width             int    `json:"width"`
	Height            int    `json:"height"`
	Duration          int    `json:"duration"`
	SupportsStreaming bool   `json:"supports_streaming"`
}

// InlineQuery is a Query from Telegram for an inline request.
type InlineQuery struct {
	ID       string    `json:"id"`
	From     *User     `json:"from"`
	Location *Location `json:"location"` // optional
	Query    string    `json:"query"`
	Offset   string    `json:"offset"`
}

// InlineQueryResultArticle is an inline query response article.
type InlineQueryResultArticle struct {
	Type                string                `json:"type"`                            // required
	ID                  string                `json:"id"`                              // required
	Title               string                `json:"title"`                           // required
	InputMessageContent interface{}           `json:"input_message_content,omitempty"` // required
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	URL                 string                `json:"url"`
	HideURL             bool                  `json:"hide_url"`
	Description         string                `json:"description"`
	ThumbURL            string                `json:"thumb_url"`
	ThumbWidth          int                   `json:"thumb_width"`
	ThumbHeight         int                   `json:"thumb_height"`
}

// InlineQueryResultPhoto is an inline query response photo.
type InlineQueryResultPhoto struct {
	Type                string                `json:"type"`      // required
	ID                  string                `json:"id"`        // required
	URL                 string                `json:"photo_url"` // required
	MimeType            string                `json:"mime_type"`
	Width               int                   `json:"photo_width"`
	Height              int                   `json:"photo_height"`
	ThumbURL            string                `json:"thumb_url"`
	Title               string                `json:"title"`
	Description         string                `json:"description"`
	Caption             string                `json:"caption"`
	ParseMode           string                `json:"parse_mode"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent interface{}           `json:"input_message_content,omitempty"`
}

// InlineQueryResultCachedPhoto is an inline query response with cached photo.
type InlineQueryResultCachedPhoto struct {
	Type                string                `json:"type"`          // required
	ID                  string                `json:"id"`            // required
	PhotoID             string                `json:"photo_file_id"` // required
	Title               string                `json:"title"`
	Description         string                `json:"description"`
	Caption             string                `json:"caption"`
	ParseMode           string                `json:"parse_mode"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent interface{}           `json:"input_message_content,omitempty"`
}

// InlineQueryResultGIF is an inline query response GIF.
type InlineQueryResultGIF struct {
	Type                string                `json:"type"`      // required
	ID                  string                `json:"id"`        // required
	URL                 string                `json:"gif_url"`   // required
	ThumbURL            string                `json:"thumb_url"` // required
	Width               int                   `json:"gif_width,omitempty"`
	Height              int                   `json:"gif_height,omitempty"`
	Duration            int                   `json:"gif_duration,omitempty"`
	Title               string                `json:"title,omitempty"`
	Caption             string                `json:"caption,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent interface{}           `json:"input_message_content,omitempty"`
}

// InlineQueryResultCachedGIF is an inline query response with cached gif.
type InlineQueryResultCachedGIF struct {
	Type                string                `json:"type"`        // required
	ID                  string                `json:"id"`          // required
	GifID               string                `json:"gif_file_id"` // required
	Title               string                `json:"title"`
	Caption             string                `json:"caption"`
	ParseMode           string                `json:"parse_mode"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent interface{}           `json:"input_message_content,omitempty"`
}

// InlineQueryResultMPEG4GIF is an inline query response MPEG4 GIF.
type InlineQueryResultMPEG4GIF struct {
	Type                string                `json:"type"`      // required
	ID                  string                `json:"id"`        // required
	URL                 string                `json:"mpeg4_url"` // required
	Width               int                   `json:"mpeg4_width"`
	Height              int                   `json:"mpeg4_height"`
	Duration            int                   `json:"mpeg4_duration"`
	ThumbURL            string                `json:"thumb_url"`
	Title               string                `json:"title"`
	Caption             string                `json:"caption"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent interface{}           `json:"input_message_content,omitempty"`
}

// InlineQueryResultCachedMpeg4Gif is an inline query response with cached
// H.264/MPEG-4 AVC video without sound gif.
type InlineQueryResultCachedMpeg4Gif struct {
	Type                string                `json:"type"`          // required
	ID                  string                `json:"id"`            // required
	MGifID              string                `json:"mpeg4_file_id"` // required
	Title               string                `json:"title"`
	Caption             string                `json:"caption"`
	ParseMode           string                `json:"parse_mode"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent interface{}           `json:"input_message_content,omitempty"`
}

// InlineQueryResultVideo is an inline query response video.
type InlineQueryResultVideo struct {
	Type                string                `json:"type"`      // required
	ID                  string                `json:"id"`        // required
	URL                 string                `json:"video_url"` // required
	MimeType            string                `json:"mime_type"` // required
	ThumbURL            string                `json:"thumb_url"`
	Title               string                `json:"title"`
	Caption             string                `json:"caption"`
	Width               int                   `json:"video_width"`
	Height              int                   `json:"video_height"`
	Duration            int                   `json:"video_duration"`
	Description         string                `json:"description"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent interface{}           `json:"input_message_content,omitempty"`
}

// InlineQueryResultCachedVideo is an inline query response with cached video.
type InlineQueryResultCachedVideo struct {
	Type                string                `json:"type"`          // required
	ID                  string                `json:"id"`            // required
	VideoID             string                `json:"video_file_id"` // required
	Title               string                `json:"title"`         // required
	Description         string                `json:"description"`
	Caption             string                `json:"caption"`
	ParseMode           string                `json:"parse_mode"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent interface{}           `json:"input_message_content,omitempty"`
}

// InlineQueryResultCachedSticker is an inline query response with cached sticker.
type InlineQueryResultCachedSticker struct {
	Type                string                `json:"type"`            // required
	ID                  string                `json:"id"`              // required
	StickerID           string                `json:"sticker_file_id"` // required
	Title               string                `json:"title"`           // required
	ParseMode           string                `json:"parse_mode"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent interface{}           `json:"input_message_content,omitempty"`
}

// InlineQueryResultAudio is an inline query response audio.
type InlineQueryResultAudio struct {
	Type                string                `json:"type"`      // required
	ID                  string                `json:"id"`        // required
	URL                 string                `json:"audio_url"` // required
	Title               string                `json:"title"`     // required
	Caption             string                `json:"caption"`
	Performer           string                `json:"performer"`
	Duration            int                   `json:"audio_duration"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent interface{}           `json:"input_message_content,omitempty"`
}

// InlineQueryResultCachedAudio is an inline query response with cached audio.
type InlineQueryResultCachedAudio struct {
	Type                string                `json:"type"`          // required
	ID                  string                `json:"id"`            // required
	AudioID             string                `json:"audio_file_id"` // required
	Caption             string                `json:"caption"`
	ParseMode           string                `json:"parse_mode"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent interface{}           `json:"input_message_content,omitempty"`
}

// InlineQueryResultVoice is an inline query response voice.
type InlineQueryResultVoice struct {
	Type                string                `json:"type"`      // required
	ID                  string                `json:"id"`        // required
	URL                 string                `json:"voice_url"` // required
	Title               string                `json:"title"`     // required
	Caption             string                `json:"caption"`
	Duration            int                   `json:"voice_duration"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent interface{}           `json:"input_message_content,omitempty"`
}

// InlineQueryResultCachedVoice is an inline query response with cached voice.
type InlineQueryResultCachedVoice struct {
	Type                string                `json:"type"`          // required
	ID                  string                `json:"id"`            // required
	VoiceID             string                `json:"voice_file_id"` // required
	Title               string                `json:"title"`         // required
	Caption             string                `json:"caption"`
	ParseMode           string                `json:"parse_mode"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent interface{}           `json:"input_message_content,omitempty"`
}

// InlineQueryResultDocument is an inline query response document.
type InlineQueryResultDocument struct {
	Type                string                `json:"type"`  // required
	ID                  string                `json:"id"`    // required
	Title               string                `json:"title"` // required
	Caption             string                `json:"caption"`
	URL                 string                `json:"document_url"` // required
	MimeType            string                `json:"mime_type"`    // required
	Description         string                `json:"description"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent interface{}           `json:"input_message_content,omitempty"`
	ThumbURL            string                `json:"thumb_url"`
	ThumbWidth          int                   `json:"thumb_width"`
	ThumbHeight         int                   `json:"thumb_height"`
}

// InlineQueryResultCachedDocument is an inline query response with cached document.
type InlineQueryResultCachedDocument struct {
	Type                string                `json:"type"`             // required
	ID                  string                `json:"id"`               // required
	DocumentID          string                `json:"document_file_id"` // required
	Title               string                `json:"title"`            // required
	Caption             string                `json:"caption"`
	Description         string                `json:"description"`
	ParseMode           string                `json:"parse_mode"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent interface{}           `json:"input_message_content,omitempty"`
}

// InlineQueryResultLocation is an inline query response location.
type InlineQueryResultLocation struct {
	Type                string                `json:"type"`      // required
	ID                  string                `json:"id"`        // required
	Latitude            float64               `json:"latitude"`  // required
	Longitude           float64               `json:"longitude"` // required
	Title               string                `json:"title"`     // required
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent interface{}           `json:"input_message_content,omitempty"`
	ThumbURL            string                `json:"thumb_url"`
	ThumbWidth          int                   `json:"thumb_width"`
	ThumbHeight         int                   `json:"thumb_height"`
}

// InlineQueryResultVenue is an inline query response venue.
type InlineQueryResultVenue struct {
	Type                string                `json:"type"`      // required
	ID                  string                `json:"id"`        // required
	Latitude            float64               `json:"latitude"`  // required
	Longitude           float64               `json:"longitude"` // required
	Title               string                `json:"title"`     // required
	Address             string                `json:"address"`   // required
	FoursquareID        string                `json:"foursquare_id"`
	FoursquareType      string                `json:"foursquare_type"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent interface{}           `json:"input_message_content,omitempty"`
	ThumbURL            string                `json:"thumb_url"`
	ThumbWidth          int                   `json:"thumb_width"`
	ThumbHeight         int                   `json:"thumb_height"`
}

// InlineQueryResultGame is an inline query response game.
type InlineQueryResultGame struct {
	Type          string                `json:"type"`
	ID            string                `json:"id"`
	GameShortName string                `json:"game_short_name"`
	ReplyMarkup   *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// ChosenInlineResult is an inline query result chosen by a User
type ChosenInlineResult struct {
	ResultID        string    `json:"result_id"`
	From            *User     `json:"from"`
	Location        *Location `json:"location"`
	InlineMessageID string    `json:"inline_message_id"`
	Query           string    `json:"query"`
}

// InputTextMessageContent contains text for displaying
// as an inline query result.
type InputTextMessageContent struct {
	Text                  string `json:"message_text"`
	ParseMode             string `json:"parse_mode"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
}

// InputLocationMessageContent contains a location for displaying
// as an inline query result.
type InputLocationMessageContent struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// InputVenueMessageContent contains a venue for displaying
// as an inline query result.
type InputVenueMessageContent struct {
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Title        string  `json:"title"`
	Address      string  `json:"address"`
	FoursquareID string  `json:"foursquare_id"`
}

// InputContactMessageContent contains a contact for displaying
// as an inline query result.
type InputContactMessageContent struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
}

// Invoice contains basic information about an invoice.
type Invoice struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	StartParameter string `json:"start_parameter"`
	Currency       string `json:"currency"`
	TotalAmount    int    `json:"total_amount"`
}

// LabeledPrice represents a portion of the price for goods or services.
type LabeledPrice struct {
	Label  string `json:"label"`
	Amount int    `json:"amount"`
}

// ShippingAddress represents a shipping address.
type ShippingAddress struct {
	CountryCode string `json:"country_code"`
	State       string `json:"state"`
	City        string `json:"city"`
	StreetLine1 string `json:"street_line1"`
	StreetLine2 string `json:"street_line2"`
	PostCode    string `json:"post_code"`
}

// OrderInfo represents information about an order.
type OrderInfo struct {
	Name            string           `json:"name,omitempty"`
	PhoneNumber     string           `json:"phone_number,omitempty"`
	Email           string           `json:"email,omitempty"`
	ShippingAddress *ShippingAddress `json:"shipping_address,omitempty"`
}

// ShippingOption represents one shipping option.
type ShippingOption struct {
	ID     string          `json:"id"`
	Title  string          `json:"title"`
	Prices *[]LabeledPrice `json:"prices"`
}

// SuccessfulPayment contains basic information about a successful payment.
type SuccessfulPayment struct {
	Currency                string     `json:"currency"`
	TotalAmount             int        `json:"total_amount"`
	InvoicePayload          string     `json:"invoice_payload"`
	ShippingOptionID        string     `json:"shipping_option_id,omitempty"`
	OrderInfo               *OrderInfo `json:"order_info,omitempty"`
	TelegramPaymentChargeID string     `json:"telegram_payment_charge_id"`
	ProviderPaymentChargeID string     `json:"provider_payment_charge_id"`
}

// ShippingQuery contains information about an incoming shipping query.
type ShippingQuery struct {
	ID              string           `json:"id"`
	From            *User            `json:"from"`
	InvoicePayload  string           `json:"invoice_payload"`
	ShippingAddress *ShippingAddress `json:"shipping_address"`
}

// PreCheckoutQuery contains information about an incoming pre-checkout query.
type PreCheckoutQuery struct {
	ID               string     `json:"id"`
	From             *User      `json:"from"`
	Currency         string     `json:"currency"`
	TotalAmount      int        `json:"total_amount"`
	InvoicePayload   string     `json:"invoice_payload"`
	ShippingOptionID string     `json:"shipping_option_id,omitempty"`
	OrderInfo        *OrderInfo `json:"order_info,omitempty"`
}

// Error is an error containing extra information returned by the Telegram API.
type Error struct {
	Code    int
	Message string
	ResponseParameters
}

func (e Error) Error() string {
	return e.Message
}

// BotCommand represents a bot command.
type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}
