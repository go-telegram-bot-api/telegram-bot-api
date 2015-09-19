package tgbotapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/technoweenie/multipartstreamer"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
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

// MessageConfig contains information about a SendMessage request.
type MessageConfig struct {
	ChatID                int
	Text                  string
	ParseMode             string
	DisableWebPagePreview bool
	ReplyToMessageID      int
	ReplyMarkup           interface{}
}

// ForwardConfig contains information about a ForwardMessage request.
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
	File             interface{}
	FileID           string
}

// AudioConfig contains information about a SendAudio request.
type AudioConfig struct {
	ChatID           int
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
	ChatID              int
	ReplyToMessageID    int
	ReplyMarkup         interface{}
	UseExistingDocument bool
	FilePath            string
	File                interface{}
	FileID              string
}

// StickerConfig contains information about a SendSticker request.
type StickerConfig struct {
	ChatID             int
	ReplyToMessageID   int
	ReplyMarkup        interface{}
	UseExistingSticker bool
	FilePath           string
	File               interface{}
	FileID             string
}

// VideoConfig contains information about a SendVideo request.
type VideoConfig struct {
	ChatID           int
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
	ChatID           int
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

// MakeRequest makes a request to a specific endpoint with our token.
// All requests are POSTs because Telegram doesn't care, and it's easier.
func (bot *BotAPI) MakeRequest(endpoint string, params url.Values) (APIResponse, error) {
	resp, err := bot.Client.PostForm(fmt.Sprintf(APIEndpoint, bot.Token, endpoint), params)
	if err != nil {
		return APIResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return APIResponse{}, errors.New(APIForbidden)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return APIResponse{}, err
	}

	if bot.Debug {
		log.Println(endpoint, string(bytes))
	}

	var apiResp APIResponse
	json.Unmarshal(bytes, &apiResp)

	if !apiResp.Ok {
		return APIResponse{}, errors.New(apiResp.Description)
	}

	return apiResp, nil
}

// UploadFile makes a request to the API with a file.
//
// Requires the parameter to hold the file not be in the params.
// File should be a string to a file path, a FileBytes struct, or a FileReader struct.
func (bot *BotAPI) UploadFile(endpoint string, params map[string]string, fieldname string, file interface{}) (APIResponse, error) {
	ms := multipartstreamer.New()
	ms.WriteFields(params)

	switch f := file.(type) {
	case string:
		fileHandle, err := os.Open(f)
		if err != nil {
			return APIResponse{}, err
		}
		defer fileHandle.Close()

		fi, err := os.Stat(f)
		if err != nil {
			return APIResponse{}, err
		}

		ms.WriteReader(fieldname, fileHandle.Name(), fi.Size(), fileHandle)
	case FileBytes:
		buf := bytes.NewBuffer(f.Bytes)
		ms.WriteReader(fieldname, f.Name, int64(len(f.Bytes)), buf)
	case FileReader:
		if f.Size == -1 {
			data, err := ioutil.ReadAll(f.Reader)
			if err != nil {
				return APIResponse{}, err
			}
			buf := bytes.NewBuffer(data)

			ms.WriteReader(fieldname, f.Name, int64(len(data)), buf)

			break
		}

		ms.WriteReader(fieldname, f.Name, f.Size, f.Reader)
	default:
		return APIResponse{}, errors.New("bad file type")
	}

	req, err := http.NewRequest("POST", fmt.Sprintf(APIEndpoint, bot.Token, endpoint), nil)
	ms.SetupRequest(req)
	if err != nil {
		return APIResponse{}, err
	}

	res, err := bot.Client.Do(req)
	if err != nil {
		return APIResponse{}, err
	}
	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return APIResponse{}, err
	}

	if bot.Debug {
		log.Println(string(bytes[:]))
	}

	var apiResp APIResponse
	json.Unmarshal(bytes, &apiResp)

	return apiResp, nil
}

// GetMe fetches the currently authenticated bot.
//
// There are no parameters for this method.
func (bot *BotAPI) GetMe() (User, error) {
	resp, err := bot.MakeRequest("getMe", nil)
	if err != nil {
		return User{}, err
	}

	var user User
	json.Unmarshal(resp.Result, &user)

	if bot.Debug {
		log.Printf("getMe: %+v\n", user)
	}

	return user, nil
}

// SendMessage sends a Message to a chat.
//
// Requires ChatID and Text.
// DisableWebPagePreview, ReplyToMessageID, and ReplyMarkup are optional.
func (bot *BotAPI) SendMessage(config MessageConfig) (Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(config.ChatID))
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
			return Message{}, err
		}

		v.Add("reply_markup", string(data))
	}

	resp, err := bot.MakeRequest("SendMessage", v)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.Debug {
		log.Printf("SendMessage req : %+v\n", v)
		log.Printf("SendMessage resp: %+v\n", message)
	}

	return message, nil
}

// ForwardMessage forwards a message from one chat to another.
//
// Requires ChatID (destination), FromChatID (source), and MessageID.
func (bot *BotAPI) ForwardMessage(config ForwardConfig) (Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(config.ChatID))
	v.Add("from_chat_id", strconv.Itoa(config.FromChatID))
	v.Add("message_id", strconv.Itoa(config.MessageID))

	resp, err := bot.MakeRequest("forwardMessage", v)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.Debug {
		log.Printf("forwardMessage req : %+v\n", v)
		log.Printf("forwardMessage resp: %+v\n", message)
	}

	return message, nil
}

// SendPhoto sends or uploads a photo to a chat.
//
// Requires ChatID and FileID OR File.
// Caption, ReplyToMessageID, and ReplyMarkup are optional.
// File should be either a string, FileBytes, or FileReader.
func (bot *BotAPI) SendPhoto(config PhotoConfig) (Message, error) {
	if config.UseExistingPhoto {
		v := url.Values{}
		v.Add("chat_id", strconv.Itoa(config.ChatID))
		v.Add("photo", config.FileID)
		if config.Caption != "" {
			v.Add("caption", config.Caption)
		}
		if config.ReplyToMessageID != 0 {
			v.Add("reply_to_message_id", strconv.Itoa(config.ChatID))
		}
		if config.ReplyMarkup != nil {
			data, err := json.Marshal(config.ReplyMarkup)
			if err != nil {
				return Message{}, err
			}

			v.Add("reply_markup", string(data))
		}

		resp, err := bot.MakeRequest("SendPhoto", v)
		if err != nil {
			return Message{}, err
		}

		var message Message
		json.Unmarshal(resp.Result, &message)

		if bot.Debug {
			log.Printf("SendPhoto req : %+v\n", v)
			log.Printf("SendPhoto resp: %+v\n", message)
		}

		return message, nil
	}

	params := make(map[string]string)
	params["chat_id"] = strconv.Itoa(config.ChatID)
	if config.Caption != "" {
		params["caption"] = config.Caption
	}
	if config.ReplyToMessageID != 0 {
		params["reply_to_message_id"] = strconv.Itoa(config.ReplyToMessageID)
	}
	if config.ReplyMarkup != nil {
		data, err := json.Marshal(config.ReplyMarkup)
		if err != nil {
			return Message{}, err
		}

		params["reply_markup"] = string(data)
	}

	var file interface{}
	if config.FilePath == "" {
		file = config.File
	} else {
		file = config.FilePath
	}

	resp, err := bot.UploadFile("SendPhoto", params, "photo", file)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.Debug {
		log.Printf("SendPhoto resp: %+v\n", message)
	}

	return message, nil
}

// SendAudio sends or uploads an audio clip to a chat.
// If using a file, the file must be in the .mp3 format.
//
// When the fields title and performer are both empty and
// the mime-type of the file to be sent is not audio/mpeg,
// the file must be an .ogg file encoded with OPUS.
// You may use the tgutils.EncodeAudio func to assist you with this, if needed.
//
// Requires ChatID and FileID OR File.
// ReplyToMessageID and ReplyMarkup are optional.
// File should be either a string, FileBytes, or FileReader.
func (bot *BotAPI) SendAudio(config AudioConfig) (Message, error) {
	if config.UseExistingAudio {
		v := url.Values{}
		v.Add("chat_id", strconv.Itoa(config.ChatID))
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
				return Message{}, err
			}

			v.Add("reply_markup", string(data))
		}
		if config.Performer != "" {
			v.Add("performer", config.Performer)
		}
		if config.Title != "" {
			v.Add("title", config.Title)
		}

		resp, err := bot.MakeRequest("sendAudio", v)
		if err != nil {
			return Message{}, err
		}

		var message Message
		json.Unmarshal(resp.Result, &message)

		if bot.Debug {
			log.Printf("sendAudio req : %+v\n", v)
			log.Printf("sendAudio resp: %+v\n", message)
		}

		return message, nil
	}

	params := make(map[string]string)

	params["chat_id"] = strconv.Itoa(config.ChatID)
	if config.ReplyToMessageID != 0 {
		params["reply_to_message_id"] = strconv.Itoa(config.ReplyToMessageID)
	}
	if config.Duration != 0 {
		params["duration"] = strconv.Itoa(config.Duration)
	}
	if config.ReplyMarkup != nil {
		data, err := json.Marshal(config.ReplyMarkup)
		if err != nil {
			return Message{}, err
		}

		params["reply_markup"] = string(data)
	}
	if config.Performer != "" {
		params["performer"] = config.Performer
	}
	if config.Title != "" {
		params["title"] = config.Title
	}

	var file interface{}
	if config.FilePath == "" {
		file = config.File
	} else {
		file = config.FilePath
	}

	resp, err := bot.UploadFile("sendAudio", params, "audio", file)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.Debug {
		log.Printf("sendAudio resp: %+v\n", message)
	}

	return message, nil
}

// SendDocument sends or uploads a document to a chat.
//
// Requires ChatID and FileID OR File.
// ReplyToMessageID and ReplyMarkup are optional.
// File should be either a string, FileBytes, or FileReader.
func (bot *BotAPI) SendDocument(config DocumentConfig) (Message, error) {
	if config.UseExistingDocument {
		v := url.Values{}
		v.Add("chat_id", strconv.Itoa(config.ChatID))
		v.Add("document", config.FileID)
		if config.ReplyToMessageID != 0 {
			v.Add("reply_to_message_id", strconv.Itoa(config.ReplyToMessageID))
		}
		if config.ReplyMarkup != nil {
			data, err := json.Marshal(config.ReplyMarkup)
			if err != nil {
				return Message{}, err
			}

			v.Add("reply_markup", string(data))
		}

		resp, err := bot.MakeRequest("sendDocument", v)
		if err != nil {
			return Message{}, err
		}

		var message Message
		json.Unmarshal(resp.Result, &message)

		if bot.Debug {
			log.Printf("sendDocument req : %+v\n", v)
			log.Printf("sendDocument resp: %+v\n", message)
		}

		return message, nil
	}

	params := make(map[string]string)

	params["chat_id"] = strconv.Itoa(config.ChatID)
	if config.ReplyToMessageID != 0 {
		params["reply_to_message_id"] = strconv.Itoa(config.ReplyToMessageID)
	}
	if config.ReplyMarkup != nil {
		data, err := json.Marshal(config.ReplyMarkup)
		if err != nil {
			return Message{}, err
		}

		params["reply_markup"] = string(data)
	}

	var file interface{}
	if config.FilePath == "" {
		file = config.File
	} else {
		file = config.FilePath
	}

	resp, err := bot.UploadFile("sendDocument", params, "document", file)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.Debug {
		log.Printf("sendDocument resp: %+v\n", message)
	}

	return message, nil
}

// SendVoice sends or uploads a playable voice to a chat.
// If using a file, the file must be encoded as an .ogg with OPUS.
// You may use the tgutils.EncodeAudio func to assist you with this, if needed.
//
// Requires ChatID and FileID OR File.
// ReplyToMessageID and ReplyMarkup are optional.
// File should be either a string, FileBytes, or FileReader.
func (bot *BotAPI) SendVoice(config VoiceConfig) (Message, error) {
	if config.UseExistingVoice {
		v := url.Values{}
		v.Add("chat_id", strconv.Itoa(config.ChatID))
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
				return Message{}, err
			}

			v.Add("reply_markup", string(data))
		}

		resp, err := bot.MakeRequest("sendVoice", v)
		if err != nil {
			return Message{}, err
		}

		var message Message
		json.Unmarshal(resp.Result, &message)

		if bot.Debug {
			log.Printf("SendVoice req : %+v\n", v)
			log.Printf("SendVoice resp: %+v\n", message)
		}

		return message, nil
	}

	params := make(map[string]string)

	params["chat_id"] = strconv.Itoa(config.ChatID)
	if config.ReplyToMessageID != 0 {
		params["reply_to_message_id"] = strconv.Itoa(config.ReplyToMessageID)
	}
	if config.Duration != 0 {
		params["duration"] = strconv.Itoa(config.Duration)
	}
	if config.ReplyMarkup != nil {
		data, err := json.Marshal(config.ReplyMarkup)
		if err != nil {
			return Message{}, err
		}

		params["reply_markup"] = string(data)
	}

	var file interface{}
	if config.FilePath == "" {
		file = config.File
	} else {
		file = config.FilePath
	}

	resp, err := bot.UploadFile("SendVoice", params, "voice", file)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.Debug {
		log.Printf("SendVoice resp: %+v\n", message)
	}

	return message, nil
}

// SendSticker sends or uploads a sticker to a chat.
//
// Requires ChatID and FileID OR File.
// ReplyToMessageID and ReplyMarkup are optional.
// File should be either a string, FileBytes, or FileReader.
func (bot *BotAPI) SendSticker(config StickerConfig) (Message, error) {
	if config.UseExistingSticker {
		v := url.Values{}
		v.Add("chat_id", strconv.Itoa(config.ChatID))
		v.Add("sticker", config.FileID)
		if config.ReplyToMessageID != 0 {
			v.Add("reply_to_message_id", strconv.Itoa(config.ReplyToMessageID))
		}
		if config.ReplyMarkup != nil {
			data, err := json.Marshal(config.ReplyMarkup)
			if err != nil {
				return Message{}, err
			}

			v.Add("reply_markup", string(data))
		}

		resp, err := bot.MakeRequest("sendSticker", v)
		if err != nil {
			return Message{}, err
		}

		var message Message
		json.Unmarshal(resp.Result, &message)

		if bot.Debug {
			log.Printf("sendSticker req : %+v\n", v)
			log.Printf("sendSticker resp: %+v\n", message)
		}

		return message, nil
	}

	params := make(map[string]string)

	params["chat_id"] = strconv.Itoa(config.ChatID)
	if config.ReplyToMessageID != 0 {
		params["reply_to_message_id"] = strconv.Itoa(config.ReplyToMessageID)
	}
	if config.ReplyMarkup != nil {
		data, err := json.Marshal(config.ReplyMarkup)
		if err != nil {
			return Message{}, err
		}

		params["reply_markup"] = string(data)
	}

	var file interface{}
	if config.FilePath == "" {
		file = config.File
	} else {
		file = config.FilePath
	}

	resp, err := bot.UploadFile("sendSticker", params, "sticker", file)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.Debug {
		log.Printf("sendSticker resp: %+v\n", message)
	}

	return message, nil
}

// SendVideo sends or uploads a video to a chat.
//
// Requires ChatID and FileID OR File.
// ReplyToMessageID and ReplyMarkup are optional.
// File should be either a string, FileBytes, or FileReader.
func (bot *BotAPI) SendVideo(config VideoConfig) (Message, error) {
	if config.UseExistingVideo {
		v := url.Values{}
		v.Add("chat_id", strconv.Itoa(config.ChatID))
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
				return Message{}, err
			}

			v.Add("reply_markup", string(data))
		}

		resp, err := bot.MakeRequest("sendVideo", v)
		if err != nil {
			return Message{}, err
		}

		var message Message
		json.Unmarshal(resp.Result, &message)

		if bot.Debug {
			log.Printf("sendVideo req : %+v\n", v)
			log.Printf("sendVideo resp: %+v\n", message)
		}

		return message, nil
	}

	params := make(map[string]string)

	params["chat_id"] = strconv.Itoa(config.ChatID)
	if config.ReplyToMessageID != 0 {
		params["reply_to_message_id"] = strconv.Itoa(config.ReplyToMessageID)
	}
	if config.ReplyMarkup != nil {
		data, err := json.Marshal(config.ReplyMarkup)
		if err != nil {
			return Message{}, err
		}

		params["reply_markup"] = string(data)
	}

	var file interface{}
	if config.FilePath == "" {
		file = config.File
	} else {
		file = config.FilePath
	}

	resp, err := bot.UploadFile("sendVideo", params, "video", file)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.Debug {
		log.Printf("sendVideo resp: %+v\n", message)
	}

	return message, nil
}

// SendLocation sends a location to a chat.
//
// Requires ChatID, Latitude, and Longitude.
// ReplyToMessageID and ReplyMarkup are optional.
func (bot *BotAPI) SendLocation(config LocationConfig) (Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(config.ChatID))
	v.Add("latitude", strconv.FormatFloat(config.Latitude, 'f', 6, 64))
	v.Add("longitude", strconv.FormatFloat(config.Longitude, 'f', 6, 64))
	if config.ReplyToMessageID != 0 {
		v.Add("reply_to_message_id", strconv.Itoa(config.ReplyToMessageID))
	}
	if config.ReplyMarkup != nil {
		data, err := json.Marshal(config.ReplyMarkup)
		if err != nil {
			return Message{}, err
		}

		v.Add("reply_markup", string(data))
	}

	resp, err := bot.MakeRequest("sendLocation", v)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.Debug {
		log.Printf("sendLocation req : %+v\n", v)
		log.Printf("sendLocation resp: %+v\n", message)
	}

	return message, nil
}

// SendChatAction sets a current action in a chat.
//
// Requires ChatID and a valid Action (see Chat constants).
func (bot *BotAPI) SendChatAction(config ChatActionConfig) error {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(config.ChatID))
	v.Add("action", config.Action)

	_, err := bot.MakeRequest("sendChatAction", v)
	if err != nil {
		return err
	}

	return nil
}

// GetUserProfilePhotos gets a user's profile photos.
//
// Requires UserID.
// Offset and Limit are optional.
func (bot *BotAPI) GetUserProfilePhotos(config UserProfilePhotosConfig) (UserProfilePhotos, error) {
	v := url.Values{}
	v.Add("user_id", strconv.Itoa(config.UserID))
	if config.Offset != 0 {
		v.Add("offset", strconv.Itoa(config.Offset))
	}
	if config.Limit != 0 {
		v.Add("limit", strconv.Itoa(config.Limit))
	}

	resp, err := bot.MakeRequest("getUserProfilePhotos", v)
	if err != nil {
		return UserProfilePhotos{}, err
	}

	var profilePhotos UserProfilePhotos
	json.Unmarshal(resp.Result, &profilePhotos)

	if bot.Debug {
		log.Printf("getUserProfilePhotos req : %+v\n", v)
		log.Printf("getUserProfilePhotos resp: %+v\n", profilePhotos)
	}

	return profilePhotos, nil
}

// GetFile returns a file_id required to download a file.
//
// Requires FileID.
func (bot *BotAPI) GetFile(config FileConfig) (File, error) {
	v := url.Values{}
	v.Add("file_id", config.FileID)

	resp, err := bot.MakeRequest("getFile", v)
	if err != nil {
		return File{}, err
	}

	var file File
	json.Unmarshal(resp.Result, &file)

	if bot.Debug {
		log.Printf("getFile req : %+v\n", v)
		log.Printf("getFile resp: %+v\n", file)
	}

	return file, nil
}

// GetUpdates fetches updates.
// If a WebHook is set, this will not return any data!
//
// Offset, Limit, and Timeout are optional.
// To not get old items, set Offset to one higher than the previous item.
// Set Timeout to a large number to reduce requests and get responses instantly.
func (bot *BotAPI) GetUpdates(config UpdateConfig) ([]Update, error) {
	v := url.Values{}
	if config.Offset > 0 {
		v.Add("offset", strconv.Itoa(config.Offset))
	}
	if config.Limit > 0 {
		v.Add("limit", strconv.Itoa(config.Limit))
	}
	if config.Timeout > 0 {
		v.Add("timeout", strconv.Itoa(config.Timeout))
	}

	resp, err := bot.MakeRequest("getUpdates", v)
	if err != nil {
		return []Update{}, err
	}

	var updates []Update
	json.Unmarshal(resp.Result, &updates)

	if bot.Debug {
		log.Printf("getUpdates: %+v\n", updates)
	}

	return updates, nil
}

// SetWebhook sets a webhook.
// If this is set, GetUpdates will not get any data!
//
// Requires Url OR to set Clear to true.
func (bot *BotAPI) SetWebhook(config WebhookConfig) (APIResponse, error) {
	if config.Certificate == nil {
		v := url.Values{}
		if !config.Clear {
			v.Add("url", config.URL.String())
		}

		return bot.MakeRequest("setWebhook", v)
	}

	params := make(map[string]string)
	params["url"] = config.URL.String()

	resp, err := bot.UploadFile("setWebhook", params, "certificate", config.Certificate)
	if err != nil {
		return APIResponse{}, err
	}

	var apiResp APIResponse
	json.Unmarshal(resp.Result, &apiResp)

	if bot.Debug {
		log.Printf("setWebhook resp: %+v\n", apiResp)
	}

	return apiResp, nil
}
