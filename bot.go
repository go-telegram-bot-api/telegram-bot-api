// Package tgbotapi has functions and types used for interacting with
// the Telegram Bot API.
package tgbotapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/technoweenie/multipartstreamer"
)

// BotAPI allows you to interact with the Telegram Bot API.
type BotAPI struct {
	Token  string `json:"token"`
	Debug  bool   `json:"debug"`
	Buffer int    `json:"buffer"`

	Self   User         `json:"-"`
	Client *http.Client `json:"-"`
}

// NewBotAPI creates a new BotAPI instance.
//
// It requires a token, provided by @BotFather on Telegram.
func NewBotAPI(token string) (*BotAPI, error) {
	return NewBotAPIWithClient(token, &http.Client{})
}

// NewBotAPIWithClient creates a new BotAPI instance
// and allows you to pass a http.Client.
//
// It requires a token, provided by @BotFather on Telegram.
func NewBotAPIWithClient(token string, client *http.Client) (*BotAPI, error) {
	bot := &BotAPI{
		Token:  token,
		Client: client,
		Buffer: 100,
	}

	self, err := bot.GetMe()
	if err != nil {
		return nil, err
	}

	bot.Self = self

	return bot, nil
}

// MakeRequest makes a request to a specific endpoint with our token.
func (bot *BotAPI) MakeRequest(endpoint string, params url.Values) (APIResponse, error) {
	method := fmt.Sprintf(APIEndpoint, bot.Token, endpoint)

	resp, err := bot.Client.PostForm(method, params)
	if err != nil {
		return APIResponse{}, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	bytes, err := bot.decodeAPIResponse(resp.Body, &apiResp)
	if err != nil {
		return apiResp, err
	}

	if bot.Debug {
		log.Printf("%s resp: %s", endpoint, bytes)
	}

	if !apiResp.Ok {
		return apiResp, errors.New(apiResp.Description)
	}

	return apiResp, nil
}

// decodeAPIResponse decode response and return slice of bytes if debug enabled.
// If debug disabled, just decode http.Response.Body stream to APIResponse struct
// for efficient memory usage
func (bot *BotAPI) decodeAPIResponse(responseBody io.Reader, resp *APIResponse) (_ []byte, err error) {
	if !bot.Debug {
		dec := json.NewDecoder(responseBody)
		err = dec.Decode(resp)
		return
	}

	// if debug, read reponse body
	data, err := ioutil.ReadAll(responseBody)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, resp)
	if err != nil {
		return
	}

	return data, nil
}

// UploadFile makes a request to the API with a file.
//
// Requires the parameter to hold the file not be in the params.
// File should be a string to a file path, a FileBytes struct,
// a FileReader struct, or a url.URL.
//
// Note that if your FileReader has a size set to -1, it will read
// the file into memory to calculate a size.
func (bot *BotAPI) UploadFile(endpoint string, params map[string]string, fieldname string, file interface{}) (APIResponse, error) {
	ms := multipartstreamer.New()

	switch f := file.(type) {
	case string:
		ms.WriteFields(params)

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
		ms.WriteFields(params)

		buf := bytes.NewBuffer(f.Bytes)
		ms.WriteReader(fieldname, f.Name, int64(len(f.Bytes)), buf)
	case FileReader:
		ms.WriteFields(params)

		if f.Size != -1 {
			ms.WriteReader(fieldname, f.Name, f.Size, f.Reader)

			break
		}

		data, err := ioutil.ReadAll(f.Reader)
		if err != nil {
			return APIResponse{}, err
		}

		buf := bytes.NewBuffer(data)

		ms.WriteReader(fieldname, f.Name, int64(len(data)), buf)
	case url.URL:
		params[fieldname] = f.String()

		ms.WriteFields(params)
	default:
		return APIResponse{}, errors.New(ErrBadFileType)
	}

	method := fmt.Sprintf(APIEndpoint, bot.Token, endpoint)

	req, err := http.NewRequest("POST", method, nil)
	if err != nil {
		return APIResponse{}, err
	}

	ms.SetupRequest(req)

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
		log.Println(string(bytes))
	}

	var apiResp APIResponse

	err = json.Unmarshal(bytes, &apiResp)
	if err != nil {
		return APIResponse{}, err
	}

	if !apiResp.Ok {
		return APIResponse{}, errors.New(apiResp.Description)
	}

	return apiResp, nil
}

// GetFileDirectURL returns direct URL to file
//
// It requires the FileID.
func (bot *BotAPI) GetFileDirectURL(fileID string) (string, error) {
	file, err := bot.GetFile(FileConfig{fileID})

	if err != nil {
		return "", err
	}

	return file.Link(bot.Token), nil
}

// GetMe fetches the currently authenticated bot.
//
// This method is called upon creation to validate the token,
// and so you may get this data from BotAPI.Self without the need for
// another request.
func (bot *BotAPI) GetMe() (User, error) {
	resp, err := bot.MakeRequest("getMe", nil)
	if err != nil {
		return User{}, err
	}

	var user User
	json.Unmarshal(resp.Result, &user)

	bot.debugLog("getMe", nil, user)

	return user, nil
}

// IsMessageToMe returns true if message directed to this bot.
//
// It requires the Message.
func (bot *BotAPI) IsMessageToMe(message Message) bool {
	return strings.Contains(message.Text, "@"+bot.Self.UserName)
}

// Send will send a Chattable item to Telegram.
//
// It requires the Chattable to send.
func (bot *BotAPI) Send(c Chattable) (Message, error) {
	resp, err := bot.Request(c)
	if err != nil {
		return Message{}, err
	}

	var message Message
	err = json.Unmarshal(resp.Result, &message)

	return message, err
}

// Request makes a request to Telegram that returns an APIResponse, rather than
// a Message.
func (bot *BotAPI) Request(c Chattable) (APIResponse, error) {
	switch t := c.(type) {
	case Fileable:
		if t.useExistingFile() {
			v, err := t.values()
			if err != nil {
				return APIResponse{}, err
			}

			return bot.MakeRequest(t.method(), v)
		}

		p, err := t.params()
		if err != nil {
			return APIResponse{}, err
		}

		return bot.UploadFile(t.method(), p, t.name(), t.getFile())
	default:
		v, err := c.values()
		if err != nil {
			return APIResponse{}, err
		}

		return bot.MakeRequest(c.method(), v)
	}
}

// debugLog checks if the bot is currently running in debug mode, and if
// so will display information about the request and response in the
// debug log.
func (bot *BotAPI) debugLog(context string, v url.Values, message interface{}) {
	if bot.Debug {
		log.Printf("%s req : %+v\n", context, v)
		log.Printf("%s resp: %+v\n", context, message)
	}
}

// GetUserProfilePhotos gets a user's profile photos.
//
// It requires UserID.
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

	bot.debugLog("GetUserProfilePhoto", v, profilePhotos)

	return profilePhotos, nil
}

// GetFile returns a File which can download a file from Telegram.
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

	bot.debugLog("GetFile", v, file)

	return file, nil
}

// GetUpdates fetches updates.
// If a WebHook is set, this will not return any data!
//
// Offset, Limit, and Timeout are optional.
// To avoid stale items, set Offset to one higher than the previous item.
// Set Timeout to a large number to reduce requests so you can get updates
// instantly instead of having to wait between requests.
func (bot *BotAPI) GetUpdates(config UpdateConfig) ([]Update, error) {
	v := url.Values{}
	if config.Offset != 0 {
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

	bot.debugLog("getUpdates", v, updates)

	return updates, nil
}

// GetWebhookInfo allows you to fetch information about a webhook and if
// one currently is set, along with pending update count and error messages.
func (bot *BotAPI) GetWebhookInfo() (WebhookInfo, error) {
	resp, err := bot.MakeRequest("getWebhookInfo", url.Values{})
	if err != nil {
		return WebhookInfo{}, err
	}

	var info WebhookInfo
	err = json.Unmarshal(resp.Result, &info)

	return info, err
}

// GetUpdatesChan starts and returns a channel for getting updates.
func (bot *BotAPI) GetUpdatesChan(config UpdateConfig) (UpdatesChannel, error) {
	ch := make(chan Update, bot.Buffer)

	go func() {
		for {
			updates, err := bot.GetUpdates(config)
			if err != nil {
				log.Println(err)
				log.Println("Failed to get updates, retrying in 3 seconds...")
				time.Sleep(time.Second * 3)

				continue
			}

			for _, update := range updates {
				if update.UpdateID >= config.Offset {
					config.Offset = update.UpdateID + 1
					ch <- update
				}
			}
		}
	}()

	return ch, nil
}

// ListenForWebhook registers a http handler for a webhook.
func (bot *BotAPI) ListenForWebhook(pattern string) UpdatesChannel {
	ch := make(chan Update, bot.Buffer)

	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := ioutil.ReadAll(r.Body)

		var update Update
		json.Unmarshal(bytes, &update)

		ch <- update
	})

	return ch
}

// GetChat gets information about a chat.
func (bot *BotAPI) GetChat(config ChatConfig) (Chat, error) {
	v := url.Values{}

	if config.SuperGroupUsername == "" {
		v.Add("chat_id", strconv.FormatInt(config.ChatID, 10))
	} else {
		v.Add("chat_id", config.SuperGroupUsername)
	}

	resp, err := bot.MakeRequest("getChat", v)
	if err != nil {
		return Chat{}, err
	}

	var chat Chat
	err = json.Unmarshal(resp.Result, &chat)

	bot.debugLog("getChat", v, chat)

	return chat, err
}

// GetChatAdministrators gets a list of administrators in the chat.
//
// If none have been appointed, only the creator will be returned.
// Bots are not shown, even if they are an administrator.
func (bot *BotAPI) GetChatAdministrators(config ChatConfig) ([]ChatMember, error) {
	v := url.Values{}

	if config.SuperGroupUsername == "" {
		v.Add("chat_id", strconv.FormatInt(config.ChatID, 10))
	} else {
		v.Add("chat_id", config.SuperGroupUsername)
	}

	resp, err := bot.MakeRequest("getChatAdministrators", v)
	if err != nil {
		return []ChatMember{}, err
	}

	var members []ChatMember
	err = json.Unmarshal(resp.Result, &members)

	bot.debugLog("getChatAdministrators", v, members)

	return members, err
}

// GetChatMembersCount gets the number of users in a chat.
func (bot *BotAPI) GetChatMembersCount(config ChatConfig) (int, error) {
	v := url.Values{}

	if config.SuperGroupUsername == "" {
		v.Add("chat_id", strconv.FormatInt(config.ChatID, 10))
	} else {
		v.Add("chat_id", config.SuperGroupUsername)
	}

	resp, err := bot.MakeRequest("getChatMembersCount", v)
	if err != nil {
		return -1, err
	}

	var count int
	err = json.Unmarshal(resp.Result, &count)

	bot.debugLog("getChatMembersCount", v, count)

	return count, err
}

// GetChatMember gets a specific chat member.
func (bot *BotAPI) GetChatMember(config ChatConfigWithUser) (ChatMember, error) {
	v := url.Values{}

	if config.SuperGroupUsername == "" {
		v.Add("chat_id", strconv.FormatInt(config.ChatID, 10))
	} else {
		v.Add("chat_id", config.SuperGroupUsername)
	}
	v.Add("user_id", strconv.Itoa(config.UserID))

	resp, err := bot.MakeRequest("getChatMember", v)
	if err != nil {
		return ChatMember{}, err
	}

	var member ChatMember
	err = json.Unmarshal(resp.Result, &member)

	bot.debugLog("getChatMember", v, member)

	return member, err
}

// GetGameHighScores allows you to get the high scores for a game.
func (bot *BotAPI) GetGameHighScores(config GetGameHighScoresConfig) ([]GameHighScore, error) {
	v, _ := config.values()

	resp, err := bot.MakeRequest(config.method(), v)
	if err != nil {
		return []GameHighScore{}, err
	}

	var highScores []GameHighScore
	err = json.Unmarshal(resp.Result, &highScores)

	return highScores, err
}

// GetInviteLink get InviteLink for a chat
func (bot *BotAPI) GetInviteLink(config ChatConfig) (string, error) {
	v := url.Values{}

	if config.SuperGroupUsername == "" {
		v.Add("chat_id", strconv.FormatInt(config.ChatID, 10))
	} else {
		v.Add("chat_id", config.SuperGroupUsername)
	}

	resp, err := bot.MakeRequest("exportChatInviteLink", v)
	if err != nil {
		return "", err
	}

	var inviteLink string
	err = json.Unmarshal(resp.Result, &inviteLink)

	return inviteLink, err
}

// GetStickerSet returns a StickerSet.
func (bot *BotAPI) GetStickerSet(config GetStickerSetConfig) (StickerSet, error) {
	v, err := config.values()
	if err != nil {
		return StickerSet{}, nil
	}

	resp, err := bot.MakeRequest(config.method(), v)
	if err != nil {
		return StickerSet{}, nil
	}

	var stickers StickerSet
	err = json.Unmarshal(resp.Result, &stickers)

	return stickers, err
}
