// Package tgbotapi has functions and types used for interacting with
// the Telegram Bot API.
package tgbotapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// HTTPClient is the type needed for the bot to perform HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// BotAPI allows you to interact with the Telegram Bot API.
type BotAPI struct {
	BotAPIConfig

	Self User
}

// BotAPIConfig contains information about how the Bot API should operate.
type BotAPIConfig struct {
	Token       string
	APIEndpoint string

	Debug  bool
	Buffer int

	Client HTTPClient
	Logger BotLogger
}

// BotLogger is an interface that represents the required methods to log data.
//
// Instead of requiring the standard logger, we can just specify the methods we
// use and allow users to pass anything that implements these.
type BotLogger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

// NewBotAPI creates a new BotAPI instance.
//
// It requires a token, provided by @BotFather on Telegram.
func NewBotAPI(token string) (*BotAPI, error) {
	return NewBotAPIWithConfig(BotAPIConfig{
		Token: token,
	})
}

func NewBotAPIWithConfig(config BotAPIConfig) (*BotAPI, error) {
	bot := &BotAPI{
		BotAPIConfig: config,
	}

	if bot.APIEndpoint == "" {
		bot.APIEndpoint = APIEndpoint
	}

	if bot.Client == nil {
		bot.Client = &http.Client{}
	}

	if bot.Logger == nil {
		bot.Logger = &log.Logger{}
	}

	self, err := bot.GetMe()
	if err != nil {
		return nil, err
	}

	bot.Self = self

	return bot, nil
}

// NewBotAPIWithAPIEndpoint creates a new BotAPI instance
// and allows you to pass API endpoint.
//
// It requires a token, provided by @BotFather on Telegram and API endpoint.
func NewBotAPIWithAPIEndpoint(token, apiEndpoint string) (*BotAPI, error) {
	return NewBotAPIWithConfig(BotAPIConfig{
		Token:       token,
		APIEndpoint: apiEndpoint,
	})
}

// NewBotAPIWithClient creates a new BotAPI instance
// and allows you to pass a http.Client.
//
// It requires a token, provided by @BotFather on Telegram and API endpoint.
func NewBotAPIWithClient(token, apiEndpoint string, client HTTPClient) (*BotAPI, error) {
	return NewBotAPIWithConfig(BotAPIConfig{
		Token:       token,
		APIEndpoint: apiEndpoint,
		Client:      client,
	})
}

func buildParams(in Params) (out url.Values) {
	if in == nil {
		return url.Values{}
	}

	out = url.Values{}

	for key, value := range in {
		out.Set(key, value)
	}

	return
}

// MakeRequest makes a request to a specific endpoint with our token.
func (bot *BotAPI) MakeRequest(endpoint string, params Params) (*APIResponse, error) {
	return bot.MakeRequestWithContext(context.Background(), endpoint, params)
}

func (bot *BotAPI) MakeRequestWithContext(ctx context.Context, endpoint string, params Params) (*APIResponse, error) {
	if bot.Debug {
		bot.Logger.Printf("endpoint: %s, params: %v\n", endpoint, params)
	}

	method := fmt.Sprintf(bot.APIEndpoint, bot.Token, endpoint)

	values := buildParams(params)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, method, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := bot.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	bytes, err := bot.decodeAPIResponse(resp.Body, &apiResp)
	if err != nil {
		return &apiResp, err
	}

	if bot.Debug {
		bot.Logger.Printf("endpoint: %s, response: %s\n", endpoint, string(bytes))
	}

	if !apiResp.Ok {
		var parameters ResponseParameters

		if apiResp.Parameters != nil {
			parameters = *apiResp.Parameters
		}

		return &apiResp, &Error{
			Code:               apiResp.ErrorCode,
			Message:            apiResp.Description,
			ResponseParameters: parameters,
		}
	}

	return &apiResp, nil
}

// decodeAPIResponse decode response and return slice of bytes if debug enabled.
// If debug disabled, just decode http.Response.Body stream to APIResponse struct
// for efficient memory usage.
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

// UploadFiles makes a request to the API with files.
func (bot *BotAPI) UploadFiles(endpoint string, params Params, files []RequestFile) (*APIResponse, error) {
	return bot.UploadFilesWithContext(context.Background(), endpoint, params, files)
}

func (bot *BotAPI) UploadFilesWithContext(ctx context.Context, endpoint string, params Params, files []RequestFile) (*APIResponse, error) {
	r, w := io.Pipe()
	m := multipart.NewWriter(w)

	// This code modified from the very helpful @HirbodBehnam
	// https://github.com/go-telegram-bot-api/telegram-bot-api/issues/354#issuecomment-663856473
	go func() {
		defer w.Close()
		defer m.Close()

		for field, value := range params {
			if err := m.WriteField(field, value); err != nil {
				_ = w.CloseWithError(err)
				return
			}
		}

		for _, file := range files {
			switch f := file.File.(type) {
			case string:
				fileHandle, err := os.Open(f)
				if err != nil {
					_ = w.CloseWithError(err)
					return
				}
				defer fileHandle.Close()

				part, err := m.CreateFormFile(file.Name, fileHandle.Name())
				if err != nil {
					_ = w.CloseWithError(err)
					return
				}

				if _, err := io.Copy(part, fileHandle); err != nil {
					_ = w.CloseWithError(err)
					return
				}
			case FileBytes:
				part, err := m.CreateFormFile(file.Name, f.Name)
				if err != nil {
					_ = w.CloseWithError(err)
					return
				}

				buf := bytes.NewBuffer(f.Bytes)
				if _, err := io.Copy(part, buf); err != nil {
					_ = w.CloseWithError(err)
					return
				}
			case FileReader:
				part, err := m.CreateFormFile(file.Name, f.Name)
				if err != nil {
					_ = w.CloseWithError(err)
					return
				}

				if _, err := io.Copy(part, f.Reader); err != nil {
					_ = w.CloseWithError(err)
					return
				}
			case FileURL:
				val := string(f)
				if err := m.WriteField(file.Name, val); err != nil {
					_ = w.CloseWithError(err)
					return
				}
			case FileID:
				val := string(f)
				if err := m.WriteField(file.Name, val); err != nil {
					_ = w.CloseWithError(err)
					return
				}
			default:
				_ = w.CloseWithError(ErrBadFileType)
				return
			}
		}
	}()

	if bot.Debug {
		bot.Logger.Printf("endpoint: %s, params: %v, with %d files\n", endpoint, params, len(files))
	}

	method := fmt.Sprintf(bot.APIEndpoint, bot.Token, endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, method, r)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", m.FormDataContentType())

	resp, err := bot.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	bytes, err := bot.decodeAPIResponse(resp.Body, &apiResp)
	if err != nil {
		return &apiResp, err
	}

	if bot.Debug {
		bot.Logger.Printf("endpoint: %s, response: %s\n", endpoint, string(bytes))
	}

	if !apiResp.Ok {
		var parameters ResponseParameters

		if apiResp.Parameters != nil {
			parameters = *apiResp.Parameters
		}

		return &apiResp, &Error{
			Message:            apiResp.Description,
			ResponseParameters: parameters,
		}
	}

	return &apiResp, nil
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
	err = json.Unmarshal(resp.Result, &user)

	return user, err
}

// IsMessageToMe returns true if message directed to this bot.
//
// It requires the Message.
func (bot *BotAPI) IsMessageToMe(message Message) bool {
	return strings.Contains(message.Text, "@"+bot.Self.UserName)
}

func hasFilesNeedingUpload(files []RequestFile) bool {
	for _, file := range files {
		switch file.File.(type) {
		case string, FileBytes, FileReader:
			return true
		}
	}

	return false
}

// Request sends a Chattable to Telegram, and returns the APIResponse.
func (bot *BotAPI) Request(c Chattable) (*APIResponse, error) {
	return bot.RequestWithContext(context.Background(), c)
}

func (bot *BotAPI) RequestWithContext(ctx context.Context, c Chattable) (*APIResponse, error) {
	params, err := c.params()
	if err != nil {
		return nil, err
	}

	if t, ok := c.(Fileable); ok {
		files := t.files()

		// If we have files that need to be uploaded, we should delegate the
		// request to UploadFile.
		if hasFilesNeedingUpload(files) {
			return bot.UploadFiles(t.method(), params, files)
		}

		// However, if there are no files to be uploaded, there's likely things
		// that need to be turned into params instead.
		for _, file := range files {
			var s string

			switch f := file.File.(type) {
			case string:
				s = f
			case FileID:
				s = string(f)
			case FileURL:
				s = string(f)
			default:
				return nil, ErrBadFileType
			}

			params[file.Name] = s
		}
	}

	return bot.MakeRequestWithContext(ctx, c.method(), params)
}

// Send will send a Chattable item to Telegram and provides the
// returned Message.
func (bot *BotAPI) Send(c Chattable) (Message, error) {
	return bot.SendWithContext(context.Background(), c)
}

func (bot *BotAPI) SendWithContext(ctx context.Context, c Chattable) (Message, error) {
	resp, err := bot.RequestWithContext(ctx, c)
	if err != nil {
		return Message{}, err
	}

	var message Message
	err = json.Unmarshal(resp.Result, &message)

	return message, err
}

// SendMediaGroup sends a media group and returns the resulting messages.
func (bot *BotAPI) SendMediaGroup(config MediaGroupConfig) ([]Message, error) {
	resp, err := bot.Request(config)
	if err != nil {
		return nil, err
	}

	var messages []Message
	err = json.Unmarshal(resp.Result, &messages)

	return messages, err
}

// GetUserProfilePhotos gets a user's profile photos.
//
// It requires UserID.
// Offset and Limit are optional.
func (bot *BotAPI) GetUserProfilePhotos(config UserProfilePhotosConfig) (UserProfilePhotos, error) {
	resp, err := bot.Request(config)
	if err != nil {
		return UserProfilePhotos{}, err
	}

	var profilePhotos UserProfilePhotos
	err = json.Unmarshal(resp.Result, &profilePhotos)

	return profilePhotos, err
}

// GetFile returns a File which can download a file from Telegram.
//
// Requires FileID.
func (bot *BotAPI) GetFile(config FileConfig) (File, error) {
	resp, err := bot.Request(config)
	if err != nil {
		return File{}, err
	}

	var file File
	err = json.Unmarshal(resp.Result, &file)

	return file, err
}

// GetUpdates fetches updates.
// If a WebHook is set, this will not return any data!
//
// Offset, Limit, Timeout, and AllowedUpdates are optional.
// To avoid stale items, set Offset to one higher than the previous item.
// Set Timeout to a large number to reduce requests so you can get updates
// instantly instead of having to wait between requests.
func (bot *BotAPI) GetUpdates(config UpdateConfig) ([]Update, error) {
	return bot.GetUpdatesWithContext(context.Background(), config)
}

func (bot *BotAPI) GetUpdatesWithContext(ctx context.Context, config UpdateConfig) ([]Update, error) {
	resp, err := bot.RequestWithContext(ctx, config)
	if err != nil {
		return nil, err
	}

	var updates []Update
	err = json.Unmarshal(resp.Result, &updates)

	return updates, err
}

// GetWebhookInfo allows you to fetch information about a webhook and if
// one currently is set, along with pending update count and error messages.
func (bot *BotAPI) GetWebhookInfo() (WebhookInfo, error) {
	resp, err := bot.MakeRequest("getWebhookInfo", nil)
	if err != nil {
		return WebhookInfo{}, err
	}

	var info WebhookInfo
	err = json.Unmarshal(resp.Result, &info)

	return info, err
}

// GetUpdatesChan starts and returns a channel for getting updates.
func (bot *BotAPI) GetUpdatesChan(config UpdateConfig) UpdatesChannel {
	return bot.GetUpdatesChanWithContext(context.Background(), config)
}

func (bot *BotAPI) GetUpdatesChanWithContext(ctx context.Context, config UpdateConfig) UpdatesChannel {
	ch := make(chan Update, bot.Buffer)

	go func() {
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				updates, err := bot.GetUpdatesWithContext(ctx, config)

				if err != nil {
					bot.Logger.Printf("failed to get updates: %s\n", err)
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
		}
	}()

	return ch
}

// ListenForWebhook registers a http handler for a webhook.
func (bot *BotAPI) ListenForWebhook(pattern string) UpdatesChannel {
	ch := make(chan Update, bot.Buffer)

	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		update, err := bot.HandleUpdate(r)
		if err != nil {
			errMsg, _ := json.Marshal(map[string]string{"error": err.Error()})
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(errMsg)
			return
		}

		ch <- *update
	})

	return ch
}

// HandleUpdate parses and returns update received via webhook.
func (bot *BotAPI) HandleUpdate(r *http.Request) (*Update, error) {
	if r.Method != http.MethodPost {
		return nil, fmt.Errorf("got %s request: %w", r.Method, ErrWrongMethod)
	}

	var update Update
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		return nil, err
	}

	return &update, nil
}

// WriteToHTTPResponse writes the request to the HTTP ResponseWriter.
//
// It doesn't support uploading files.
//
// See https://core.telegram.org/bots/api#making-requests-when-getting-updates
// for details.
func WriteToHTTPResponse(w http.ResponseWriter, c Chattable) error {
	params, err := c.params()
	if err != nil {
		return err
	}

	if t, ok := c.(Fileable); ok {
		if hasFilesNeedingUpload(t.files()) {
			return ErrDisallowedUploads
		}
	}

	values := buildParams(params)
	values.Set("method", c.method())

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	_, err = w.Write([]byte(values.Encode()))
	return err
}

// GetChat gets information about a chat.
func (bot *BotAPI) GetChat(config ChatInfoConfig) (Chat, error) {
	resp, err := bot.Request(config)
	if err != nil {
		return Chat{}, err
	}

	var chat Chat
	err = json.Unmarshal(resp.Result, &chat)

	return chat, err
}

// GetChatAdministrators gets a list of administrators in the chat.
//
// If none have been appointed, only the creator will be returned.
// Bots are not shown, even if they are an administrator.
func (bot *BotAPI) GetChatAdministrators(config ChatAdministratorsConfig) ([]ChatMember, error) {
	resp, err := bot.Request(config)
	if err != nil {
		return []ChatMember{}, err
	}

	var members []ChatMember
	err = json.Unmarshal(resp.Result, &members)

	return members, err
}

// GetChatMembersCount gets the number of users in a chat.
func (bot *BotAPI) GetChatMembersCount(config ChatMemberCountConfig) (int, error) {
	resp, err := bot.Request(config)
	if err != nil {
		return -1, err
	}

	var count int
	err = json.Unmarshal(resp.Result, &count)

	return count, err
}

// GetChatMember gets a specific chat member.
func (bot *BotAPI) GetChatMember(config GetChatMemberConfig) (ChatMember, error) {
	resp, err := bot.Request(config)
	if err != nil {
		return ChatMember{}, err
	}

	var member ChatMember
	err = json.Unmarshal(resp.Result, &member)

	return member, err
}

// GetGameHighScores allows you to get the high scores for a game.
func (bot *BotAPI) GetGameHighScores(config GetGameHighScoresConfig) ([]GameHighScore, error) {
	resp, err := bot.Request(config)
	if err != nil {
		return []GameHighScore{}, err
	}

	var highScores []GameHighScore
	err = json.Unmarshal(resp.Result, &highScores)

	return highScores, err
}

// GetInviteLink get InviteLink for a chat.
func (bot *BotAPI) GetInviteLink(config ChatInviteLinkConfig) (string, error) {
	resp, err := bot.Request(config)
	if err != nil {
		return "", err
	}

	var inviteLink string
	err = json.Unmarshal(resp.Result, &inviteLink)

	return inviteLink, err
}

// GetStickerSet returns a StickerSet.
func (bot *BotAPI) GetStickerSet(config GetStickerSetConfig) (StickerSet, error) {
	resp, err := bot.Request(config)
	if err != nil {
		return StickerSet{}, err
	}

	var stickers StickerSet
	err = json.Unmarshal(resp.Result, &stickers)

	return stickers, err
}

// StopPoll stops a poll and returns the result.
func (bot *BotAPI) StopPoll(config StopPollConfig) (Poll, error) {
	resp, err := bot.Request(config)
	if err != nil {
		return Poll{}, err
	}

	var poll Poll
	err = json.Unmarshal(resp.Result, &poll)

	return poll, err
}

// GetMyCommands gets the currently registered commands.
func (bot *BotAPI) GetMyCommands() ([]BotCommand, error) {
	config := GetMyCommandsConfig{}

	resp, err := bot.Request(config)
	if err != nil {
		return nil, err
	}

	var commands []BotCommand
	err = json.Unmarshal(resp.Result, &commands)

	return commands, err
}

// CopyMessage copy messages of any kind. The method is analogous to the method
// forwardMessage, but the copied message doesn't have a link to the original
// message. Returns the MessageID of the sent message on success.
func (bot *BotAPI) CopyMessage(config CopyMessageConfig) (MessageID, error) {
	params, err := config.params()
	if err != nil {
		return MessageID{}, err
	}

	resp, err := bot.MakeRequest(config.method(), params)
	if err != nil {
		return MessageID{}, err
	}

	var messageID MessageID
	err = json.Unmarshal(resp.Result, &messageID)

	return messageID, err
}
