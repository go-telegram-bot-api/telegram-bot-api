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
	Result      json.RawMessage     `json:"result,omitempty"`
	ErrorCode   int                 `json:"error_code,omitempty"`
	Description string              `json:"description,omitempty"`
	Parameters  *ResponseParameters `json:"parameters,omitempty"`
}

// ResponseParameters are various errors that can be returned in APIResponse.
type ResponseParameters struct {
	MigrateToChatID int64 `json:"migrate_to_chat_id,omitempty"` // optional
	RetryAfter      int   `json:"retry_after,omitempty"`        // optional
}

// Update is an update response, from GetUpdates.
type Update struct {
	UpdateID           int                 `json:"update_id"`
	Message            *Message            `json:"message,omitempty"`
	EditedMessage      *Message            `json:"edited_message,omitempty"`
	ChannelPost        *Message            `json:"channel_post,omitempty"`
	EditedChannelPost  *Message            `json:"edited_channel_post,omitempty"`
	InlineQuery        *InlineQuery        `json:"inline_query,omitempty"`
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result,omitempty"`
	CallbackQuery      *CallbackQuery      `json:"callback_query,omitempty"`
	ShippingQuery      *ShippingQuery      `json:"shipping_query,omitempty"`
	PreCheckoutQuery   *PreCheckoutQuery   `json:"pre_checkout_query,omitempty"`
	Poll               *Poll               `json:"poll,omitempty"`
	PollAnswer         *PollAnswer         `json:"poll_answer,omitempty"`
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
	ID                      int    `json:"id"`
	FirstName               string `json:"first_name"`
	LastName                string `json:"last_name,omitempty"`                   // optional
	UserName                string `json:"username,omitempty"`                    // optional
	LanguageCode            string `json:"language_code,omitempty"`               // optional
	IsBot                   bool   `json:"is_bot,omitempty"`                      // optional
	CanJoinGroups           bool   `json:"can_join_groups,omitempty"`             // optional
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages,omitempty"` // optional
	SupportsInlineQueries   bool   `json:"supports_inline_queries,omitempty"`     // optional
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

// ChatPhoto represents a chat photo.
type ChatPhoto struct {
	SmallFileID       string `json:"small_file_id"`
	SmallFileUniqueID string `json:"small_file_unique_id"`
	BigFileID         string `json:"big_file_id"`
	BigFileUniqueID   string `json:"big_file_unique_id"`
}

// ChatPermissions describes actions that a non-administrator user is
// allowed to take in a chat. All fields are optional.
type ChatPermissions struct {
	CanSendMessages       bool `json:"can_send_messages"`
	CanSendMediaMessages  bool `json:"can_send_media_messages"`
	CanSendPolls          bool `json:"can_send_polls"`
	CanSendOtherMessages  bool `json:"can_send_other_messages"`
	CanAddWebPagePreviews bool `json:"can_add_web_page_previews"`
	CanChangeInfo         bool `json:"can_change_info"`
	CanInviteUsers        bool `json:"can_invite_users"`
	CanPinMessages        bool `json:"can_pin_messages"`
}

// Chat contains information about the place a message was sent.
type Chat struct {
	ID                  int64            `json:"id"`
	Type                string           `json:"type"`
	Title               string           `json:"title,omitempty"`                          // optional
	UserName            string           `json:"username,omitempty"`                       // optional
	FirstName           string           `json:"first_name,omitempty"`                     // optional
	LastName            string           `json:"last_name,omitempty"`                      // optional
	AllMembersAreAdmins bool             `json:"all_members_are_administrators,omitempty"` // deprecated, optional
	Photo               *ChatPhoto       `json:"photo,omitempty"`                          // optional
	Description         string           `json:"description,omitempty"`                    // optional
	InviteLink          string           `json:"invite_link,omitempty"`                    // optional
	PinnedMessage       *Message         `json:"pinned_message,omitempty"`                 // optional
	Permissions         *ChatPermissions `json:"permissions,omitempty"`                    // optional
	SlowModeDelay       int              `json:"slow_mode_delay,omitempty"`                // optional
	StickerSetName      string           `json:"sticker_set_name,omitempty"`               // optional
	CanSetStickerSet    bool             `json:"can_set_sticker_set,omitempty"`            // optional
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
	MessageID             int                   `json:"message_id"`
	From                  *User                 `json:"from,omitempty"` // optional
	Date                  int                   `json:"date"`
	Chat                  *Chat                 `json:"chat"`
	ForwardFrom           *User                 `json:"forward_from,omitempty"`            // optional
	ForwardFromChat       *Chat                 `json:"forward_from_chat,omitempty"`       // optional
	ForwardFromMessageID  int                   `json:"forward_from_message_id,omitempty"` // optional
	ForwardSignature      string                `json:"forward_signature,omitempty"`       // optional
	ForwardSenderName     string                `json:"forward_sender_name,omitempty"`     // optional
	ForwardDate           int                   `json:"forward_date,omitempty"`            // optional
	ReplyToMessage        *Message              `json:"reply_to_message,omitempty"`        // optional
	EditDate              int                   `json:"edit_date,omitempty"`               // optional
	MediaGroupID          string                `json:"media_group_id,omitempty"`          // optional
	AuthorSignature       string                `json:"author_signature,omitempty"`        // optional
	Text                  string                `json:"text,omitempty"`                    // optional
	Entities              []MessageEntity       `json:"entities,omitempty"`                // optional
	CaptionEntities       []MessageEntity       `json:"caption_entities,omitempty"`        // optional
	Audio                 *Audio                `json:"audio,omitempty"`                   // optional
	Document              *Document             `json:"document,omitempty"`                // optional
	Animation             *ChatAnimation        `json:"animation,omitempty"`               // optional
	Game                  *Game                 `json:"game,omitempty"`                    // optional
	Photo                 []PhotoSize           `json:"photo,omitempty"`                   // optional
	Sticker               *Sticker              `json:"sticker,omitempty"`                 // optional
	Video                 *Video                `json:"video,omitempty"`                   // optional
	VideoNote             *VideoNote            `json:"video_note,omitempty"`              // optional
	Voice                 *Voice                `json:"voice,omitempty"`                   // optional
	Caption               string                `json:"caption,omitempty"`                 // optional
	Contact               *Contact              `json:"contact,omitempty"`                 // optional
	Location              *Location             `json:"location,omitempty"`                // optional
	Venue                 *Venue                `json:"venue,omitempty"`                   // optional
	Poll                  *Poll                 `json:"poll,omitempty"`                    // optional
	Dice                  *Dice                 `json:"dice,omitempty"`                    // optional
	NewChatMembers        []User                `json:"new_chat_members,omitempty"`        // optional
	LeftChatMember        *User                 `json:"left_chat_member,omitempty"`        // optional
	NewChatTitle          string                `json:"new_chat_title,omitempty"`          // optional
	NewChatPhoto          []PhotoSize           `json:"new_chat_photo,omitempty"`          // optional
	DeleteChatPhoto       bool                  `json:"delete_chat_photo,omitempty"`       // optional
	GroupChatCreated      bool                  `json:"group_chat_created,omitempty"`      // optional
	SuperGroupChatCreated bool                  `json:"supergroup_chat_created,omitempty"` // optional
	ChannelChatCreated    bool                  `json:"channel_chat_created,omitempty"`    // optional
	MigrateToChatID       int64                 `json:"migrate_to_chat_id,omitempty"`      // optional
	MigrateFromChatID     int64                 `json:"migrate_from_chat_id,omitempty"`    // optional
	PinnedMessage         *Message              `json:"pinned_message,omitempty"`          // optional
	Invoice               *Invoice              `json:"invoice,omitempty"`                 // optional
	SuccessfulPayment     *SuccessfulPayment    `json:"successful_payment,omitempty"`      // optional
	ConnectedWebsite      string                `json:"connected_website,omitempty"`       // optional
	PassportData          *PassportData         `json:"passport_data,omitempty"`           // optional
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`            // optional
}

// Time converts the message timestamp into a Time.
func (m *Message) Time() time.Time {
	return time.Unix(int64(m.Date), 0)
}

// IsCommand returns true if message starts with a "bot_command" entity.
func (m *Message) IsCommand() bool {
	if m.Entities == nil || len(m.Entities) == 0 {
		return false
	}

	entity := m.Entities[0]
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
	entity := m.Entities[0]
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
	entity := m.Entities[0]

	if len(m.Text) == entity.Length {
		return "" // The command makes up the whole message
	}

	return m.Text[entity.Length+1:]
}

// MessageEntity contains information about data in a Message.
type MessageEntity struct {
	Type     string `json:"type"`
	Offset   int    `json:"offset"`
	Length   int    `json:"length"`
	URL      string `json:"url,omitempty"`      // optional
	User     *User  `json:"user,omitempty"`     // optional
	Language string `json:"language,omitempty"` // optional
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
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int    `json:"file_size,omitempty"` // optional
}

// Audio contains information about audio.
type Audio struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Duration     int    `json:"duration"`
	Performer    string `json:"performer,omitempty"` // optional
	Title        string `json:"title,omitempty"`     // optional
	MimeType     string `json:"mime_type,omitempty"` // optional
	FileSize     int    `json:"file_size,omitempty"` // optional
}

// Document contains information about a document.
type Document struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Thumbnail    *PhotoSize `json:"thumb,omitempty"`     // optional
	FileName     string     `json:"file_name,omitempty"` // optional
	MimeType     string     `json:"mime_type,omitempty"` // optional
	FileSize     int        `json:"file_size,omitempty"` // optional
}

// Sticker contains information about a sticker.
type Sticker struct {
	FileID       string       `json:"file_id"`
	FileUniqueID string       `json:"file_unique_id"`
	Width        int          `json:"width"`
	Height       int          `json:"height"`
	IsAnimated   bool         `json:"is_animated"`
	Thumbnail    *PhotoSize   `json:"thumb,omitempty"`         // optional
	Emoji        string       `json:"emoji,omitempty"`         // optional
	SetName      string       `json:"set_name,omitempty"`      // optional
	MaskPosition MaskPosition `json:"mask_position,omitempty"` //optional
	FileSize     int          `json:"file_size,omitempty"`     // optional
}

// MaskPosition is the position of a mask.
type MaskPosition struct {
	Point     string     `json:"point"`
	XShift    float32    `json:"x_shift"`
	YShift    float32    `json:"y_shift"`
	Scale     float32    `json:"scale"`
	FileID    string     `json:"file_id"`
	Width     int        `json:"width"`
	Height    int        `json:"height"`
	Thumbnail *PhotoSize `json:"thumb,omitempty"`     // optional
	Emoji     string     `json:"emoji,omitempty"`     // optional
	FileSize  int        `json:"file_size,omitempty"` // optional
	SetName   string     `json:"set_name,omitempty"`  // optional
}

// ChatAnimation contains information about an animation.
type ChatAnimation struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Width        int        `json:"width"`
	Height       int        `json:"height"`
	Duration     int        `json:"duration"`
	Thumbnail    *PhotoSize `json:"thumb,omitempty"`     // optional
	FileName     string     `json:"file_name,omitempty"` // optional
	MimeType     string     `json:"mime_type,omitempty"` // optional
	FileSize     int        `json:"file_size,omitempty"` // optional
}

// Video contains information about a video.
type Video struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Width        int        `json:"width"`
	Height       int        `json:"height"`
	Duration     int        `json:"duration"`
	Thumbnail    *PhotoSize `json:"thumb,omitempty"`     // optional
	MimeType     string     `json:"mime_type,omitempty"` // optional
	FileSize     int        `json:"file_size,omitempty"` // optional
}

// VideoNote contains information about a video.
type VideoNote struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Length       int        `json:"length"`
	Duration     int        `json:"duration"`
	Thumbnail    *PhotoSize `json:"thumb,omitempty"`     // optional
	FileSize     int        `json:"file_size,omitempty"` // optional
}

// Voice contains information about a voice.
type Voice struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Duration     int    `json:"duration"`
	MimeType     string `json:"mime_type,omitempty"` // optional
	FileSize     int    `json:"file_size,omitempty"` // optional
}

// Contact contains information about a contact.
//
// Note that LastName and UserID may be empty.
type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name,omitempty"` // optional
	UserID      int    `json:"user_id,omitempty"`   // optional
	VCard       string `json:"vcard,omitempty"`     // optional
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
	FoursquareID string   `json:"foursquare_id,omitempty"` // optional
}

// PollOption contains information about one answer option in a poll.
type PollOption struct {
	Text       string `json:"text"`
	VoterCount int    `json:"voter_count"`
}

// PollAnswer represents an answer of a user in a non-anonymous poll.
type PollAnswer struct {
	PollID    string `json:"poll_id"`
	User      User   `json:"user"`
	OptionIDs []int  `json:"option_ids"`
}

// Poll contains information about a poll.
type Poll struct {
	ID                    string          `json:"id"`
	Question              string          `json:"question"`
	Options               []PollOption    `json:"options"`
	IsClosed              bool            `json:"is_closed"`
	IsAnonymous           bool            `json:"is_anonymous"`
	Type                  string          `json:"type"`
	AllowsMultipleAnswers bool            `json:"allows_multiple_answers"`
	CorrectOptionID       int             `json:"correct_option_id,omitempty"`    // optional
	Explanation           string          `json:"explanation,omitempty"`          // optional
	ExplanationEntities   []MessageEntity `json:"explanation_entities,omitempty"` // optional
	OpenPeriod            int             `json:"open_period,omitempty"`          // optional
	CloseDate             int             `json:"close_date,omitempty"`           // optional
}

// Dice represents a single dice value.
type Dice struct {
	Emoji string `json:"emoji"`
	Value int    `json:"value"`
}

// UserProfilePhotos contains a set of user profile photos.
type UserProfilePhotos struct {
	TotalCount int           `json:"total_count"`
	Photos     [][]PhotoSize `json:"photos"`
}

// File contains information about a file to download from Telegram.
type File struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int    `json:"file_size,omitempty"` // optional
	FilePath     string `json:"file_path,omitempty"` // optional
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
	ResizeKeyboard  bool               `json:"resize_keyboard,omitempty"`   // optional
	OneTimeKeyboard bool               `json:"one_time_keyboard,omitempty"` // optional
	Selective       bool               `json:"selective,omitempty"`         // optional
}

// KeyboardButton is a button within a custom keyboard.
type KeyboardButton struct {
	Text            string                  `json:"text"`
	RequestContact  bool                    `json:"request_contact"`
	RequestLocation bool                    `json:"request_location"`
	RequestPoll     *KeyboardButtonPollType `json:"request_poll,omitempty"`
}

// KeyboardButtonPollType represents type of a poll, which is allowed to
// be created and sent when the corresponding button is pressed.
type KeyboardButtonPollType struct {
	Type string `json:"type"`
}

// ReplyKeyboardHide allows the Bot to hide a custom keyboard.
type ReplyKeyboardHide struct {
	HideKeyboard bool `json:"hide_keyboard"`
	Selective    bool `json:"selective,omitempty"` // optional
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
	LoginURL                     *LoginURL     `json:"login_url,omitempty"`                        // optional
	CallbackData                 *string       `json:"callback_data,omitempty"`                    // optional
	SwitchInlineQuery            *string       `json:"switch_inline_query,omitempty"`              // optional
	SwitchInlineQueryCurrentChat *string       `json:"switch_inline_query_current_chat,omitempty"` // optional
	CallbackGame                 *CallbackGame `json:"callback_game,omitempty"`                    // optional
	Pay                          bool          `json:"pay,omitempty"`                              // optional
}

// LoginURL is the parameters for the login inline keyboard button type.
type LoginURL struct {
	URL                string `json:"url"`
	ForwardText        string `json:"forward_text"`
	BotUsername        string `json:"bot_username"`
	RequestWriteAccess bool   `json:"request_write_access"`
}

// CallbackQuery is data sent when a keyboard button with callback data
// is clicked.
type CallbackQuery struct {
	ID              string   `json:"id"`
	From            *User    `json:"from"`
	Message         *Message `json:"message,omitempty"`           // optional
	InlineMessageID string   `json:"inline_message_id,omitempty"` // optional
	ChatInstance    string   `json:"chat_instance"`
	Data            string   `json:"data,omitempty"`            // optional
	GameShortName   string   `json:"game_short_name,omitempty"` // optional
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
	CanPostMessages       bool   `json:"can_post_messages,omitempty"`         // optional
	CanEditMessages       bool   `json:"can_edit_messages,omitempty"`         // optional
	CanDeleteMessages     bool   `json:"can_delete_messages,omitempty"`       // optional
	CanRestrictMembers    bool   `json:"can_restrict_members,omitempty"`      // optional
	CanPromoteMembers     bool   `json:"can_promote_members,omitempty"`       // optional
	CanChangeInfo         bool   `json:"can_change_info,omitempty"`           // optional
	CanInviteUsers        bool   `json:"can_invite_users,omitempty"`          // optional
	CanPinMessages        bool   `json:"can_pin_messages,omitempty"`          // optional
	IsChatMember          bool   `json:"is_member,omitempty"`                 // optional
	CanSendMessages       bool   `json:"can_send_messages,omitempty"`         // optional
	CanSendMediaMessages  bool   `json:"can_send_media_messages,omitempty"`   // optional
	CanSendPolls          bool   `json:"can_send_polls,omitempty"`            // optional
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
	FileID       string    `json:"file_id"`
	FileUniqueID string    `json:"file_unique_id"`
	Thumb        PhotoSize `json:"thumb"`
	FileName     string    `json:"file_name"`
	MimeType     string    `json:"mime_type"`
	FileSize     int       `json:"file_size"`
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
	LastErrorDate        int    `json:"last_error_date,omitempty"`    // optional
	LastErrorMessage     string `json:"last_error_message,omitempty"` // optional
}

// IsSet returns true if a webhook is currently set.
func (info WebhookInfo) IsSet() bool {
	return info.URL != ""
}

// InlineQuery is a Query from Telegram for an inline request.
type InlineQuery struct {
	ID       string    `json:"id"`
	From     *User     `json:"from"`
	Location *Location `json:"location,omitempty"` // optional
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
	Type                string                `json:"type"`        // required
	ID                  string                `json:"id"`          // required
	Latitude            float64               `json:"latitude"`    // required
	Longitude           float64               `json:"longitude"`   // required
	LivePeriod          int                   `json:"live_period"` // optional
	Title               string                `json:"title"`       // required
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent interface{}           `json:"input_message_content,omitempty"`
	ThumbURL            string                `json:"thumb_url"`
	ThumbWidth          int                   `json:"thumb_width"`
	ThumbHeight         int                   `json:"thumb_height"`
}

// InlineQueryResultContact is an inline query response contact.
type InlineQueryResultContact struct {
	Type                string                `json:"type"`         // required
	ID                  string                `json:"id"`           // required
	PhoneNumber         string                `json:"phone_number"` // required
	FirstName           string                `json:"first_name"`   // required
	LastName            string                `json:"last_name"`
	VCard               string                `json:"vcard"`
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
	VCard       string `json:"vcard"`
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
	ID     string         `json:"id"`
	Title  string         `json:"title"`
	Prices []LabeledPrice `json:"prices"`
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

// StickerSet is a collection of stickers.
type StickerSet struct {
	Name          string     `json:"name"`
	Title         string     `json:"title"`
	IsAnimated    bool       `json:"is_animated"`
	ContainsMasks bool       `json:"contains_masks"`
	Stickers      []Sticker  `json:"stickers"`
	Thumb         *PhotoSize `json:"thumb"`
}

// BotCommand represents Telegram's understanding of a command.
type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

// BaseInputMedia is a base type for the InputMedia types.
type BaseInputMedia struct {
	Type      string `json:"type"`
	Media     string `json:"media"`
	Caption   string `json:"caption"`
	ParseMode string `json:"parse_mode"`
}

// InputMediaPhoto is a photo to send as part of a media group.
type InputMediaPhoto struct {
	BaseInputMedia
}

// InputMediaVideo is a video to send as part of a media group.
type InputMediaVideo struct {
	BaseInputMedia
	Width             int  `json:"width"`
	Height            int  `json:"height"`
	Duration          int  `json:"duration"`
	SupportsStreaming bool `json:"supports_streaming"`
}

// InputMediaAnimation is an animation to send as part of a media group.
type InputMediaAnimation struct {
	BaseInputMedia
	Width    int `json:"width"`
	Height   int `json:"height"`
	Duration int `json:"duration"`
}

// InputMediaAudio is a audio to send as part of a media group.
type InputMediaAudio struct {
	BaseInputMedia
	Duration  int    `json:"duration"`
	Performer string `json:"performer"`
	Title     string `json:"title"`
}

// InputMediaDocument is a audio to send as part of a media group.
type InputMediaDocument struct {
	BaseInputMedia
}

// Error is an error containing extra information returned by the Telegram API.
type Error struct {
	Code    int
	Message string
	ResponseParameters
}

// Error message string.
func (e Error) Error() string {
	return e.Message
}
