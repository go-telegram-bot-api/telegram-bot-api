package tgbotapi

import (
	"io"
	"net/url"
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
	// ErrAPIForbidden happens when a token is bad
	ErrAPIForbidden = "forbidden"
)

// Constant values for ParseMode in MessageConfig
const (
	ModeMarkdown = "Markdown"
	ModeHTML     = "HTML"
)

// Library errors
const (
	// ErrBadFileType happens when you pass an unknown type
	ErrBadFileType = "bad file type"
	ErrBadURL      = "bad or empty url"
)

// Chattable is any config type that can be sent.
type Chattable interface {
	params() (Params, error)
	method() string
}

// Fileable is any config type that can be sent that includes a file.
type Fileable interface {
	Chattable
	name() string
	getFile() interface{}
	useExistingFile() bool
}

// BaseChat is base type for all chat config types.
type BaseChat struct {
	ChatID              int64 // required
	ChannelUsername     string
	ReplyToMessageID    int
	ReplyMarkup         interface{}
	DisableNotification bool
}

// params returns Params representation of BaseChat
func (chat *BaseChat) params() (Params, error) {
	v := make(Params)

	v.AddFirstValid("chat_id", chat.ChatID, chat.ChannelUsername)
	v.AddNonZero("reply_to_message_id", chat.ReplyToMessageID)
	v.AddBool("disable_notification", chat.DisableNotification)

	err := v.AddInterface("reply_markup", chat.ReplyMarkup)

	return v, err
}

// BaseFile is a base type for all file config types.
type BaseFile struct {
	BaseChat
	File        interface{}
	FileID      string
	UseExisting bool
	MimeType    string
	FileSize    int
}

// params returns a Params representation of BaseFile.
func (file BaseFile) params() (Params, error) {
	params, err := file.BaseChat.params()

	params.AddNonEmpty("mime_type", file.MimeType)
	params.AddNonZero("file_size", file.FileSize)

	return params, err
}

// getFile returns the file.
func (file BaseFile) getFile() interface{} {
	return file.File
}

// useExistingFile returns if the BaseFile has already been uploaded.
func (file BaseFile) useExistingFile() bool {
	return file.UseExisting
}

// BaseEdit is base type of all chat edits.
type BaseEdit struct {
	ChatID          int64
	ChannelUsername string
	MessageID       int
	InlineMessageID string
	ReplyMarkup     *InlineKeyboardMarkup
}

func (edit BaseEdit) params() (Params, error) {
	v := make(Params)

	if edit.InlineMessageID != "" {
		v["inline_message_id"] = edit.InlineMessageID
	} else {
		v.AddFirstValid("chat_id", edit.ChatID, edit.ChannelUsername)
		v.AddNonZero("message_id", edit.MessageID)
	}

	err := v.AddInterface("reply_markup", edit.ReplyMarkup)

	return v, err
}

// MessageConfig contains information about a SendMessage request.
type MessageConfig struct {
	BaseChat
	Text                  string
	ParseMode             string
	DisableWebPagePreview bool
}

// values returns a url.Values representation of MessageConfig.
func (config MessageConfig) params() (Params, error) {
	v, err := config.BaseChat.params()
	if err != nil {
		return v, err
	}

	v.AddNonEmpty("text", config.Text)
	v.AddBool("disable_web_page_preview", config.DisableWebPagePreview)
	v.AddNonEmpty("parse_mode", config.ParseMode)

	return v, nil
}

// method returns Telegram API method name for sending Message.
func (config MessageConfig) method() string {
	return "sendMessage"
}

// ForwardConfig contains information about a ForwardMessage request.
type ForwardConfig struct {
	BaseChat
	FromChatID          int64 // required
	FromChannelUsername string
	MessageID           int // required
}

// values returns a url.Values representation of ForwardConfig.
func (config ForwardConfig) params() (Params, error) {
	v, err := config.BaseChat.params()
	if err != nil {
		return v, err
	}

	v.AddNonZero64("from_chat_id", config.FromChatID)
	v.AddNonZero("message_id", config.MessageID)

	return v, nil
}

// method returns Telegram API method name for sending Forward.
func (config ForwardConfig) method() string {
	return "forwardMessage"
}

// PhotoConfig contains information about a SendPhoto request.
type PhotoConfig struct {
	BaseFile
	Caption   string
	ParseMode string
}

// Params returns a map[string]string representation of PhotoConfig.
func (config PhotoConfig) params() (Params, error) {
	params, err := config.BaseFile.params()

	params.AddNonEmpty(config.name(), config.FileID)
	params.AddNonEmpty("caption", config.Caption)
	params.AddNonEmpty("parse_mode", config.ParseMode)

	return params, err
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
	Caption   string
	ParseMode string
	Duration  int
	Performer string
	Title     string
}

// values returns a url.Values representation of AudioConfig.
func (config AudioConfig) params() (Params, error) {
	v, err := config.BaseChat.params()
	if err != nil {
		return v, err
	}

	v.AddNonEmpty(config.name(), config.FileID)
	v.AddNonZero("duration", config.Duration)
	v.AddNonEmpty("performer", config.Performer)
	v.AddNonEmpty("title", config.Title)
	v.AddNonEmpty("caption", config.Caption)
	v.AddNonEmpty("parse_mode", config.ParseMode)

	return v, nil
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
	Caption   string
	ParseMode string
}

// params returns a map[string]string representation of DocumentConfig.
func (config DocumentConfig) params() (Params, error) {
	params, err := config.BaseFile.params()

	params.AddNonEmpty(config.name(), config.FileID)
	params.AddNonEmpty("caption", config.Caption)
	params.AddNonEmpty("parse_mode", config.ParseMode)

	return params, err
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
func (config StickerConfig) params() (Params, error) {
	v, err := config.BaseChat.params()

	v.AddNonEmpty(config.name(), config.FileID)

	return v, err
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
	Duration  int
	Caption   string
	ParseMode string
}

// values returns a url.Values representation of VideoConfig.
func (config VideoConfig) params() (Params, error) {
	v, err := config.BaseChat.params()

	v.AddNonEmpty(config.name(), config.FileID)
	v.AddNonZero("duration", config.Duration)
	v.AddNonEmpty("caption", config.Caption)
	v.AddNonEmpty("parse_mode", config.ParseMode)

	return v, err
}

// name returns the field name for the Video.
func (config VideoConfig) name() string {
	return "video"
}

// method returns Telegram API method name for sending Video.
func (config VideoConfig) method() string {
	return "sendVideo"
}

// AnimationConfig contains information about a SendAnimation request.
type AnimationConfig struct {
	BaseFile
	Duration  int
	Caption   string
	ParseMode string
}

// values returns a Params representation of AnimationConfig.
func (config AnimationConfig) params() (Params, error) {
	v, err := config.BaseChat.params()

	v.AddNonEmpty(config.name(), config.FileID)
	v.AddNonZero("duration", config.Duration)
	v.AddNonEmpty("caption", config.Caption)
	v.AddNonEmpty("parse_mode", config.ParseMode)

	return v, err
}

// name returns the field name for the Animation.
func (config AnimationConfig) name() string {
	return "animation"
}

// method returns Telegram API method name for sending Animation.
func (config AnimationConfig) method() string {
	return "sendAnimation"
}

// VideoNoteConfig contains information about a SendVideoNote request.
type VideoNoteConfig struct {
	BaseFile
	Duration int
	Length   int
}

// values returns a url.Values representation of VideoNoteConfig.
func (config VideoNoteConfig) params() (Params, error) {
	v, err := config.BaseChat.params()

	v.AddNonEmpty(config.name(), config.FileID)
	v.AddNonZero("duration", config.Duration)
	v.AddNonZero("length", config.Length)

	return v, err
}

// name returns the field name for the VideoNote.
func (config VideoNoteConfig) name() string {
	return "video_note"
}

// method returns Telegram API method name for sending VideoNote.
func (config VideoNoteConfig) method() string {
	return "sendVideoNote"
}

// VoiceConfig contains information about a SendVoice request.
type VoiceConfig struct {
	BaseFile
	Caption   string
	ParseMode string
	Duration  int
}

// values returns a url.Values representation of VoiceConfig.
func (config VoiceConfig) params() (Params, error) {
	v, err := config.BaseChat.params()

	v.AddNonEmpty(config.name(), config.FileID)
	v.AddNonZero("duration", config.Duration)
	v.AddNonEmpty("caption", config.Caption)
	v.AddNonEmpty("parse_mode", config.ParseMode)

	return v, err
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
	Latitude   float64 // required
	Longitude  float64 // required
	LivePeriod int     // optional
}

// values returns a url.Values representation of LocationConfig.
func (config LocationConfig) params() (Params, error) {
	v, err := config.BaseChat.params()

	v.AddNonZeroFloat("latitude", config.Latitude)
	v.AddNonZeroFloat("longitude", config.Longitude)
	v.AddNonZero("live_period", config.LivePeriod)

	return v, err
}

// method returns Telegram API method name for sending Location.
func (config LocationConfig) method() string {
	return "sendLocation"
}

// EditMessageLiveLocationConfig allows you to update a live location.
type EditMessageLiveLocationConfig struct {
	BaseEdit
	Latitude  float64 // required
	Longitude float64 // required
}

// values returns a url.Values representation of EditMessageLiveLocationConfig.
func (config EditMessageLiveLocationConfig) params() (Params, error) {
	v, err := config.BaseEdit.params()

	v.AddNonZeroFloat("latitude", config.Latitude)
	v.AddNonZeroFloat("longitude", config.Longitude)

	return v, err
}

// method returns Telegram API method name for edit message Live Location.
func (config EditMessageLiveLocationConfig) method() string {
	return "editMessageLiveLocation"
}

// StopMessageLiveLocationConfig stops updating a live location.
type StopMessageLiveLocationConfig struct {
	BaseEdit
}

// values returns a url.Values representation of StopMessageLiveLocationConfig.
func (config StopMessageLiveLocationConfig) params() (Params, error) {
	return config.BaseEdit.params()
}

// method returns Telegram API method name for stop message Live Location.
func (config StopMessageLiveLocationConfig) method() string {
	return "stopMessageLiveLocation"
}

// VenueConfig contains information about a SendVenue request.
type VenueConfig struct {
	BaseChat
	Latitude     float64 // required
	Longitude    float64 // required
	Title        string  // required
	Address      string  // required
	FoursquareID string
}

func (config VenueConfig) params() (Params, error) {
	v, err := config.BaseChat.params()

	v.AddNonZeroFloat("latitude", config.Latitude)
	v.AddNonZeroFloat("longitude", config.Longitude)
	v["title"] = config.Title
	v["address"] = config.Address
	v.AddNonEmpty("foursquare_id", config.FoursquareID)

	return v, err
}

func (config VenueConfig) method() string {
	return "sendVenue"
}

// ContactConfig allows you to send a contact.
type ContactConfig struct {
	BaseChat
	PhoneNumber string
	FirstName   string
	LastName    string
}

func (config ContactConfig) params() (Params, error) {
	v, err := config.BaseChat.params()

	v["phone_number"] = config.PhoneNumber
	v["first_name"] = config.FirstName
	v["last_name"] = config.LastName

	return v, err
}

func (config ContactConfig) method() string {
	return "sendContact"
}

// GameConfig allows you to send a game.
type GameConfig struct {
	BaseChat
	GameShortName string
}

func (config GameConfig) params() (Params, error) {
	v, err := config.BaseChat.params()

	v["game_short_name"] = config.GameShortName

	return v, err
}

func (config GameConfig) method() string {
	return "sendGame"
}

// SetGameScoreConfig allows you to update the game score in a chat.
type SetGameScoreConfig struct {
	UserID             int
	Score              int
	Force              bool
	DisableEditMessage bool
	ChatID             int64
	ChannelUsername    string
	MessageID          int
	InlineMessageID    string
}

func (config SetGameScoreConfig) params() (Params, error) {
	v := make(Params)

	v.AddNonZero("user_id", config.UserID)
	v.AddNonZero("scrore", config.Score)
	v.AddBool("disable_edit_message", config.DisableEditMessage)

	if config.InlineMessageID != "" {
		v["inline_message_id"] = config.InlineMessageID
	} else {
		v.AddFirstValid("chat_id", config.ChatID, config.ChannelUsername)
		v.AddNonZero("message_id", config.MessageID)
	}

	return v, nil
}

func (config SetGameScoreConfig) method() string {
	return "setGameScore"
}

// GetGameHighScoresConfig allows you to fetch the high scores for a game.
type GetGameHighScoresConfig struct {
	UserID          int
	ChatID          int
	ChannelUsername string
	MessageID       int
	InlineMessageID string
}

func (config GetGameHighScoresConfig) params() (Params, error) {
	v := make(Params)

	v.AddNonZero("user_id", config.UserID)

	if config.InlineMessageID != "" {
		v["inline_message_id"] = config.InlineMessageID
	} else {
		v.AddFirstValid("chat_id", config.ChatID, config.ChannelUsername)
		v.AddNonZero("message_id", config.MessageID)
	}

	return v, nil
}

func (config GetGameHighScoresConfig) method() string {
	return "getGameHighScores"
}

// ChatActionConfig contains information about a SendChatAction request.
type ChatActionConfig struct {
	BaseChat
	Action string // required
}

// values returns a url.Values representation of ChatActionConfig.
func (config ChatActionConfig) params() (Params, error) {
	v, err := config.BaseChat.params()

	v["action"] = config.Action

	return v, err
}

// method returns Telegram API method name for sending ChatAction.
func (config ChatActionConfig) method() string {
	return "sendChatAction"
}

// EditMessageTextConfig allows you to modify the text in a message.
type EditMessageTextConfig struct {
	BaseEdit
	Text                  string
	ParseMode             string
	DisableWebPagePreview bool
}

func (config EditMessageTextConfig) params() (Params, error) {
	v, err := config.BaseEdit.params()

	v["text"] = config.Text
	v.AddNonEmpty("parse_mode", config.ParseMode)
	v.AddBool("disable_web_page_preview", config.DisableWebPagePreview)

	return v, err
}

func (config EditMessageTextConfig) method() string {
	return "editMessageText"
}

// EditMessageCaptionConfig allows you to modify the caption of a message.
type EditMessageCaptionConfig struct {
	BaseEdit
	Caption   string
	ParseMode string
}

func (config EditMessageCaptionConfig) params() (Params, error) {
	v, err := config.BaseEdit.params()

	v["caption"] = config.Caption
	v.AddNonEmpty("parse_mode", config.ParseMode)

	return v, err
}

func (config EditMessageCaptionConfig) method() string {
	return "editMessageCaption"
}

// EditMessageReplyMarkupConfig allows you to modify the reply markup
// of a message.
type EditMessageReplyMarkupConfig struct {
	BaseEdit
}

func (config EditMessageReplyMarkupConfig) params() (Params, error) {
	return config.BaseEdit.params()
}

func (config EditMessageReplyMarkupConfig) method() string {
	return "editMessageReplyMarkup"
}

// UserProfilePhotosConfig contains information about a
// GetUserProfilePhotos request.
type UserProfilePhotosConfig struct {
	UserID int
	Offset int
	Limit  int
}

func (UserProfilePhotosConfig) method() string {
	return "getUserProfilePhotos"
}

func (config UserProfilePhotosConfig) params() (Params, error) {
	params := make(Params)

	params.AddNonZero("user_id", config.UserID)
	params.AddNonZero("offset", config.Offset)
	params.AddNonZero("limit", config.Limit)

	return params, nil
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

func (UpdateConfig) method() string {
	return "getUpdates"
}

func (config UpdateConfig) params() (Params, error) {
	params := make(Params)

	params.AddNonZero("offset", config.Offset)
	params.AddNonZero("limit", config.Limit)
	params.AddNonZero("timeout", config.Timeout)

	return params, nil
}

// WebhookConfig contains information about a SetWebhook request.
type WebhookConfig struct {
	URL            *url.URL
	Certificate    interface{}
	MaxConnections int
}

func (config WebhookConfig) method() string {
	return "setWebhook"
}

func (config WebhookConfig) params() (Params, error) {
	params := make(Params)

	if config.URL != nil {
		params["url"] = config.URL.String()
	}

	params.AddNonZero("max_connections", config.MaxConnections)

	return params, nil
}

func (config WebhookConfig) name() string {
	return "certificate"
}

func (config WebhookConfig) getFile() interface{} {
	return config.Certificate
}

func (config WebhookConfig) useExistingFile() bool {
	return config.URL != nil
}

// RemoveWebhookConfig is a helper to remove a webhook.
type RemoveWebhookConfig struct {
}

func (config RemoveWebhookConfig) method() string {
	return "setWebhook"
}

func (config RemoveWebhookConfig) params() (Params, error) {
	return nil, nil
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
	InlineQueryID     string        `json:"inline_query_id"`
	Results           []interface{} `json:"results"`
	CacheTime         int           `json:"cache_time"`
	IsPersonal        bool          `json:"is_personal"`
	NextOffset        string        `json:"next_offset"`
	SwitchPMText      string        `json:"switch_pm_text"`
	SwitchPMParameter string        `json:"switch_pm_parameter"`
}

func (config InlineConfig) method() string {
	return "answerInlineQuery"
}

func (config InlineConfig) params() (Params, error) {
	v := make(Params)

	v["inline_query_id"] = config.InlineQueryID
	v.AddNonZero("cache_time", config.CacheTime)
	v.AddBool("is_personal", config.IsPersonal)
	v.AddNonEmpty("next_offset", config.NextOffset)
	v.AddNonEmpty("switch_pm_text", config.SwitchPMText)
	v.AddNonEmpty("switch_pm_parameter", config.SwitchPMParameter)

	if err := v.AddInterface("results", config.Results); err != nil {
		return v, err
	}

	return v, nil
}

// CallbackConfig contains information on making a CallbackQuery response.
type CallbackConfig struct {
	CallbackQueryID string `json:"callback_query_id"`
	Text            string `json:"text"`
	ShowAlert       bool   `json:"show_alert"`
	URL             string `json:"url"`
	CacheTime       int    `json:"cache_time"`
}

func (config CallbackConfig) method() string {
	return "answerCallbackQuery"
}

func (config CallbackConfig) params() (Params, error) {
	v := make(Params)

	v["callback_query_id"] = config.CallbackQueryID
	v.AddNonEmpty("text", config.Text)
	v.AddBool("show_alert", config.ShowAlert)
	v.AddNonEmpty("url", config.URL)
	v.AddNonZero("cache_time", config.CacheTime)

	return v, nil
}

// ChatMemberConfig contains information about a user in a chat for use
// with administrative functions such as kicking or unbanning a user.
type ChatMemberConfig struct {
	ChatID             int64
	SuperGroupUsername string
	ChannelUsername    string
	UserID             int
}

// UnbanChatMemberConfig allows you to unban a user.
type UnbanChatMemberConfig struct {
	ChatMemberConfig
}

func (config UnbanChatMemberConfig) method() string {
	return "unbanChatMember"
}

func (config UnbanChatMemberConfig) params() (Params, error) {
	v := make(Params)

	v.AddFirstValid("chat_id", config.ChatID, config.SuperGroupUsername, config.ChannelUsername)
	v.AddNonZero("user_id", config.UserID)

	return v, nil
}

// KickChatMemberConfig contains extra fields to kick user
type KickChatMemberConfig struct {
	ChatMemberConfig
	UntilDate int64
}

func (config KickChatMemberConfig) method() string {
	return "kickChatMember"
}

func (config KickChatMemberConfig) params() (Params, error) {
	v := make(Params)

	v.AddFirstValid("chat_id", config.ChatID, config.SuperGroupUsername)
	v.AddNonZero("user_id", config.UserID)
	v.AddNonZero64("until_date", config.UntilDate)

	return v, nil
}

// RestrictChatMemberConfig contains fields to restrict members of chat
type RestrictChatMemberConfig struct {
	ChatMemberConfig
	UntilDate             int64
	CanSendMessages       *bool
	CanSendMediaMessages  *bool
	CanSendOtherMessages  *bool
	CanAddWebPagePreviews *bool
}

func (config RestrictChatMemberConfig) method() string {
	return "restrictChatMember"
}

func (config RestrictChatMemberConfig) params() (Params, error) {
	v := make(Params)

	v.AddFirstValid("chat_id", config.ChatID, config.SuperGroupUsername, config.ChannelUsername)
	v.AddNonZero("user_id", config.UserID)

	v.AddNonNilBool("can_send_messages", config.CanSendMessages)
	v.AddNonNilBool("can_send_media_messages", config.CanSendMediaMessages)
	v.AddNonNilBool("can_send_other_messages", config.CanSendOtherMessages)
	v.AddNonNilBool("can_add_web_page_previews", config.CanAddWebPagePreviews)
	v.AddNonZero64("until_date", config.UntilDate)

	return v, nil
}

// PromoteChatMemberConfig contains fields to promote members of chat
type PromoteChatMemberConfig struct {
	ChatMemberConfig
	CanChangeInfo      *bool
	CanPostMessages    *bool
	CanEditMessages    *bool
	CanDeleteMessages  *bool
	CanInviteUsers     *bool
	CanRestrictMembers *bool
	CanPinMessages     *bool
	CanPromoteMembers  *bool
}

func (config PromoteChatMemberConfig) method() string {
	return "promoteChatMember"
}

func (config PromoteChatMemberConfig) params() (Params, error) {
	v := make(Params)

	v.AddFirstValid("chat_id", config.ChatID, config.SuperGroupUsername, config.ChannelUsername)
	v.AddNonZero("user_id", config.UserID)

	v.AddNonNilBool("can_change_info", config.CanChangeInfo)
	v.AddNonNilBool("can_post_messages", config.CanPostMessages)
	v.AddNonNilBool("can_edit_messages", config.CanEditMessages)
	v.AddNonNilBool("can_delete_messages", config.CanDeleteMessages)
	v.AddNonNilBool("can_invite_users", config.CanInviteUsers)
	v.AddNonNilBool("can_restrict_members", config.CanRestrictMembers)
	v.AddNonNilBool("can_pin_messages", config.CanPinMessages)
	v.AddNonNilBool("can_promote_members", config.CanPromoteMembers)

	return v, nil
}

// ChatConfig contains information about getting information on a chat.
type ChatConfig struct {
	ChatID             int64
	SuperGroupUsername string
}

// LeaveChatConfig allows you to leave a chat.
type LeaveChatConfig struct {
	ChatID          int64
	ChannelUsername string
}

func (config LeaveChatConfig) method() string {
	return "leaveChat"
}

func (config LeaveChatConfig) params() (Params, error) {
	v := make(Params)

	v.AddFirstValid("chat_id", config.ChatID, config.ChannelUsername)

	return v, nil
}

// ChatConfigWithUser contains information about getting information on
// a specific user within a chat.
type ChatConfigWithUser struct {
	ChatID             int64
	SuperGroupUsername string
	UserID             int
}

// InvoiceConfig contains information for sendInvoice request.
type InvoiceConfig struct {
	BaseChat
	Title               string          // required
	Description         string          // required
	Payload             string          // required
	ProviderToken       string          // required
	StartParameter      string          // required
	Currency            string          // required
	Prices              *[]LabeledPrice // required
	ProviderData        string
	PhotoURL            string
	PhotoSize           int
	PhotoWidth          int
	PhotoHeight         int
	NeedName            bool
	NeedPhoneNumber     bool
	NeedEmail           bool
	NeedShippingAddress bool
	IsFlexible          bool
}

func (config InvoiceConfig) params() (Params, error) {
	v, err := config.BaseChat.params()
	if err != nil {
		return v, err
	}

	v["title"] = config.Title
	v["description"] = config.Description
	v["payload"] = config.Payload
	v["provider_token"] = config.ProviderToken
	v["start_parameter"] = config.StartParameter
	v["currency"] = config.Currency

	if err = v.AddInterface("prices", config.Prices); err != nil {
		return v, err
	}

	v.AddNonEmpty("provider_data", config.ProviderData)
	v.AddNonEmpty("photo_url", config.PhotoURL)
	v.AddNonZero("photo_size", config.PhotoSize)
	v.AddNonZero("photo_width", config.PhotoWidth)
	v.AddNonZero("photo_height", config.PhotoHeight)
	v.AddBool("need_name", config.NeedName)
	v.AddBool("need_phone_number", config.NeedPhoneNumber)
	v.AddBool("need_email", config.NeedEmail)
	v.AddBool("need_shipping_address", config.NeedShippingAddress)
	v.AddBool("is_flexible", config.IsFlexible)

	return v, nil
}

func (config InvoiceConfig) method() string {
	return "sendInvoice"
}

// ShippingConfig contains information for answerShippingQuery request.
type ShippingConfig struct {
	ShippingQueryID string // required
	OK              bool   // required
	ShippingOptions *[]ShippingOption
	ErrorMessage    string
}

// PreCheckoutConfig conatins information for answerPreCheckoutQuery request.
type PreCheckoutConfig struct {
	PreCheckoutQueryID string // required
	OK                 bool   // required
	ErrorMessage       string
}

// DeleteMessageConfig contains information of a message in a chat to delete.
type DeleteMessageConfig struct {
	ChatID    int64
	MessageID int
}

func (config DeleteMessageConfig) method() string {
	return "deleteMessage"
}

func (config DeleteMessageConfig) params() (Params, error) {
	v := make(Params)

	v.AddNonZero64("chat_id", config.ChatID)
	v.AddNonZero("message_id", config.MessageID)

	return v, nil
}

// PinChatMessageConfig contains information of a message in a chat to pin.
type PinChatMessageConfig struct {
	ChatID              int64
	ChannelUsername     string
	MessageID           int
	DisableNotification bool
}

func (config PinChatMessageConfig) method() string {
	return "pinChatMessage"
}

func (config PinChatMessageConfig) params() (Params, error) {
	v := make(Params)

	v.AddFirstValid("chat_id", config.ChatID, config.ChannelUsername)
	v.AddNonZero("message_id", config.MessageID)
	v.AddBool("disable_notification", config.DisableNotification)

	return v, nil
}

// UnpinChatMessageConfig contains information of chat to unpin.
type UnpinChatMessageConfig struct {
	ChatID          int64
	ChannelUsername string
}

func (config UnpinChatMessageConfig) method() string {
	return "unpinChatMessage"
}

func (config UnpinChatMessageConfig) params() (Params, error) {
	v := make(Params)

	v.AddFirstValid("chat_id", config.ChatID, config.ChannelUsername)

	return v, nil
}

// SetChatPhotoConfig allows you to set a group, supergroup, or channel's photo.
type SetChatPhotoConfig struct {
	BaseFile
}

func (config SetChatPhotoConfig) method() string {
	return "setChatPhoto"
}

func (config SetChatPhotoConfig) name() string {
	return "photo"
}

func (config SetChatPhotoConfig) getFile() interface{} {
	return config.File
}

func (config SetChatPhotoConfig) useExistingFile() bool {
	return config.UseExisting
}

// DeleteChatPhotoConfig allows you to delete a group, supergroup, or channel's photo.
type DeleteChatPhotoConfig struct {
	ChatID          int64
	ChannelUsername string
}

func (config DeleteChatPhotoConfig) method() string {
	return "deleteChatPhoto"
}

func (config DeleteChatPhotoConfig) params() (Params, error) {
	v := make(Params)

	v.AddFirstValid("chat_id", config.ChatID, config.ChannelUsername)

	return v, nil
}

// SetChatTitleConfig allows you to set the title of something other than a private chat.
type SetChatTitleConfig struct {
	ChatID          int64
	ChannelUsername string

	Title string
}

func (config SetChatTitleConfig) method() string {
	return "setChatTitle"
}

func (config SetChatTitleConfig) params() (Params, error) {
	v := make(Params)

	v.AddFirstValid("chat_id", config.ChatID, config.ChannelUsername)
	v["title"] = config.Title

	return v, nil
}

// SetChatDescriptionConfig allows you to set the description of a supergroup or channel.
type SetChatDescriptionConfig struct {
	ChatID          int64
	ChannelUsername string

	Description string
}

func (config SetChatDescriptionConfig) method() string {
	return "setChatDescription"
}

func (config SetChatDescriptionConfig) params() (Params, error) {
	v := make(Params)

	v.AddFirstValid("chat_id", config.ChatID, config.ChannelUsername)
	v["description"] = config.Description

	return v, nil
}

// GetStickerSetConfig allows you to get the stickers in a set.
type GetStickerSetConfig struct {
	Name string
}

func (config GetStickerSetConfig) method() string {
	return "getStickerSet"
}

func (config GetStickerSetConfig) params() (Params, error) {
	v := make(Params)

	v["name"] = config.Name

	return v, nil
}

// UploadStickerConfig allows you to upload a sticker for use in a set later.
type UploadStickerConfig struct {
	UserID     int64
	PNGSticker interface{}
}

func (config UploadStickerConfig) method() string {
	return "uploadStickerFile"
}

func (config UploadStickerConfig) params() (Params, error) {
	v := make(Params)

	v.AddNonZero64("user_id", config.UserID)

	return v, nil
}

func (config UploadStickerConfig) name() string {
	return "png_sticker"
}

func (config UploadStickerConfig) getFile() interface{} {
	return config.PNGSticker
}

func (config UploadStickerConfig) useExistingFile() bool {
	return false
}

// NewStickerSetConfig allows creating a new sticker set.
type NewStickerSetConfig struct {
	UserID        int64
	Name          string
	Title         string
	PNGSticker    interface{}
	Emojis        string
	ContainsMasks bool
	MaskPosition  *MaskPosition
}

func (config NewStickerSetConfig) method() string {
	return "createNewStickerSet"
}

func (config NewStickerSetConfig) params() (Params, error) {
	v := make(Params)

	v.AddNonZero64("user_id", config.UserID)
	v["name"] = config.Name
	v["title"] = config.Title

	if sticker, ok := config.PNGSticker.(string); ok {
		v[config.name()] = sticker
	}

	v["emojis"] = config.Emojis

	v.AddBool("contains_masks", config.ContainsMasks)

	err := v.AddInterface("mask_position", config.MaskPosition)

	return v, err
}

func (config NewStickerSetConfig) getFile() interface{} {
	return config.PNGSticker
}

func (config NewStickerSetConfig) name() string {
	return "png_sticker"
}

func (config NewStickerSetConfig) useExistingFile() bool {
	_, ok := config.PNGSticker.(string)

	return ok
}

// AddStickerConfig allows you to add a sticker to a set.
type AddStickerConfig struct {
	UserID       int64
	Name         string
	PNGSticker   interface{}
	Emojis       string
	MaskPosition *MaskPosition
}

func (config AddStickerConfig) method() string {
	return "addStickerToSet"
}

func (config AddStickerConfig) params() (Params, error) {
	v := make(Params)

	v.AddNonZero64("user_id", config.UserID)
	v["name"] = config.Name
	v["emojis"] = config.Emojis

	if sticker, ok := config.PNGSticker.(string); ok {
		v[config.name()] = sticker
	}

	err := v.AddInterface("mask_position", config.MaskPosition)

	return v, err
}

func (config AddStickerConfig) name() string {
	return "png_sticker"
}

func (config AddStickerConfig) getFile() interface{} {
	return config.PNGSticker
}

func (config AddStickerConfig) useExistingFile() bool {
	return false
}

// SetStickerPositionConfig allows you to change the position of a sticker in a set.
type SetStickerPositionConfig struct {
	Sticker  string
	Position int
}

func (config SetStickerPositionConfig) method() string {
	return "setStickerPositionInSet"
}

func (config SetStickerPositionConfig) params() (Params, error) {
	v := make(Params)

	v["sticker"] = config.Sticker
	v.AddNonZero("position", config.Position)

	return v, nil
}

// DeleteStickerConfig allows you to delete a sticker from a set.
type DeleteStickerConfig struct {
	Sticker string
}

func (config DeleteStickerConfig) method() string {
	return "deleteStickerFromSet"
}

func (config DeleteStickerConfig) params() (Params, error) {
	v := make(Params)

	v["sticker"] = config.Sticker

	return v, nil
}

// SetChatStickerSetConfig allows you to set the sticker set for a supergroup.
type SetChatStickerSetConfig struct {
	ChatID             int64
	SuperGroupUsername string

	StickerSetName string
}

func (config SetChatStickerSetConfig) method() string {
	return "setChatStickerSet"
}

func (config SetChatStickerSetConfig) params() (Params, error) {
	v := make(Params)

	v.AddFirstValid("chat_id", config.ChatID, config.SuperGroupUsername)
	v["sticker_set_name"] = config.StickerSetName

	return v, nil
}

// DeleteChatStickerSetConfig allows you to remove a supergroup's sticker set.
type DeleteChatStickerSetConfig struct {
	ChatID             int64
	SuperGroupUsername string
}

func (config DeleteChatStickerSetConfig) method() string {
	return "deleteChatStickerSet"
}

func (config DeleteChatStickerSetConfig) params() (Params, error) {
	v := make(Params)

	v.AddFirstValid("chat_id", config.ChatID, config.SuperGroupUsername)

	return v, nil
}

// MediaGroupConfig allows you to send a group of media.
//
// Media consist of InputMedia items (InputMediaPhoto, InputMediaVideo).
type MediaGroupConfig struct {
	ChatID          int64
	ChannelUsername string

	Media               []interface{}
	DisableNotification bool
	ReplyToMessageID    int
}

func (config MediaGroupConfig) method() string {
	return "sendMediaGroup"
}

func (config MediaGroupConfig) params() (Params, error) {
	v := make(Params)

	v.AddFirstValid("chat_id", config.ChatID, config.ChannelUsername)
	if err := v.AddInterface("media", config.Media); err != nil {
		return v, nil
	}
	v.AddBool("disable_notification", config.DisableNotification)
	v.AddNonZero("reply_to_message_id", config.ReplyToMessageID)

	return v, nil
}
