// Package tgbotapi has functions and types used for interacting with
// the Telegram Bot API.
package tgbotapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// HTTPClient is the type needed for the bot to perform HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// BotAPI allows you to interact with the Telegram Bot API.
type BotAPI struct {
	Token  string `json:"token"`
	Debug  bool   `json:"debug"`
	Buffer int    `json:"buffer"`

	Self   User       `json:"-"`
	Client HTTPClient `json:"-"`

	apiEndpoint string

	stoppers []context.CancelFunc
	mu       sync.RWMutex
}

// NewBotAPI creates a new BotAPI instance.
//
// It requires a token, provided by @BotFather on Telegram.
func NewBotAPI(token string) (*BotAPI, error) {
	return NewBotAPIWithClient(token, APIEndpoint, &http.Client{})
}

// NewBotAPIWithAPIEndpoint creates a new BotAPI instance
// and allows you to pass API endpoint.
//
// It requires a token, provided by @BotFather on Telegram and API endpoint.
func NewBotAPIWithAPIEndpoint(token, apiEndpoint string) (*BotAPI, error) {
	return NewBotAPIWithClient(token, apiEndpoint, &http.Client{})
}

// NewBotAPIWithClient creates a new BotAPI instance
// and allows you to pass a http.Client.
//
// It requires a token, provided by @BotFather on Telegram and API endpoint.
func NewBotAPIWithClient(token, apiEndpoint string, client HTTPClient) (*BotAPI, error) {
	bot := &BotAPI{
		Token:  token,
		Client: client,
		Buffer: 100,

		apiEndpoint: apiEndpoint,
	}

	self, err := bot.GetMe()
	if err != nil {
		return nil, err
	}

	bot.Self = self

	return bot, nil
}

// SetAPIEndpoint changes the Telegram Bot API endpoint used by the instance.
func (bot *BotAPI) SetAPIEndpoint(apiEndpoint string) {
	bot.apiEndpoint = apiEndpoint
}

func buildParams(in Params) url.Values {
	if in == nil {
		return url.Values{}
	}

	out := url.Values{}

	for key, value := range in {
		out.Set(key, value)
	}

	return out
}

// MakeRequestWithCtx makes a request to a specific endpoint with our token with context.
func (bot *BotAPI) MakeRequestWithCtx(ctx context.Context, endpoint string, params Params) (*APIResponse, error) {
	if bot.Debug {
		log.Printf("Endpoint: %s, params: %v\n", endpoint, params)
	}

	method := fmt.Sprintf(bot.apiEndpoint, bot.Token, endpoint)

	values := buildParams(params)

	req, err := http.NewRequestWithContext(ctx, "POST", method, strings.NewReader(values.Encode()))
	if err != nil {
		return &APIResponse{}, err
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
		log.Printf("Endpoint: %s, response: %s\n", endpoint, string(bytes))
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

// MakeRequest makes a request to a specific endpoint with our token.
func (bot *BotAPI) MakeRequest(endpoint string, params Params) (*APIResponse, error) {
	return bot.MakeRequestWithCtx(context.Background(), endpoint, params)
}

// decodeAPIResponse decode response and return slice of bytes if debug enabled.
// If debug disabled, just decode http.Response.Body stream to APIResponse struct
// for efficient memory usage
func (bot *BotAPI) decodeAPIResponse(responseBody io.Reader, resp *APIResponse) ([]byte, error) {
	if !bot.Debug {
		dec := json.NewDecoder(responseBody)
		err := dec.Decode(resp)
		return nil, err
	}

	// if debug, read response body
	data, err := io.ReadAll(responseBody)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (bot *BotAPI) UploadFilesWithCtx(ctx context.Context, endpoint string, params Params, files []RequestFile) (*APIResponse, error) {
	r, w := io.Pipe()
	m := multipart.NewWriter(w)

	// This code modified from the very helpful @HirbodBehnam
	// https://github.com/go-telegram-bot-api/telegram-bot-api/issues/354#issuecomment-663856473
	go func() {
		defer w.Close()
		defer m.Close()

		for field, value := range params {
			if err := m.WriteField(field, value); err != nil {
				w.CloseWithError(err)
				return
			}
		}

		for _, file := range files {
			if file.Data.NeedsUpload() {
				name, reader, err := file.Data.UploadData()
				if err != nil {
					w.CloseWithError(err)
					return
				}

				part, err := m.CreateFormFile(file.Name, name)
				if err != nil {
					w.CloseWithError(err)
					return
				}

				if _, err := io.Copy(part, reader); err != nil {
					w.CloseWithError(err)
					return
				}

				if closer, ok := reader.(io.ReadCloser); ok {
					if err = closer.Close(); err != nil {
						w.CloseWithError(err)
						return
					}
				}
			} else {
				value := file.Data.SendData()

				if err := m.WriteField(file.Name, value); err != nil {
					w.CloseWithError(err)
					return
				}
			}
		}
	}()

	if bot.Debug {
		log.Printf("Endpoint: %s, params: %v, with %d files\n", endpoint, params, len(files))
	}

	method := fmt.Sprintf(bot.apiEndpoint, bot.Token, endpoint)

	req, err := http.NewRequestWithContext(ctx, "POST", method, r)
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
		log.Printf("Endpoint: %s, response: %s\n", endpoint, string(bytes))
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

// UploadFiles makes a request to the API with files.
func (bot *BotAPI) UploadFiles(endpoint string, params Params, files []RequestFile) (*APIResponse, error) {
	return bot.UploadFilesWithCtx(context.Background(), endpoint, params, files)
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

func (bot *BotAPI) GetMeWithCtx(ctx context.Context) (User, error) {
	resp, err := bot.MakeRequestWithCtx(ctx, "getMe", nil)
	if err != nil {
		return User{}, err
	}

	var user User
	err = json.Unmarshal(resp.Result, &user)

	return user, err
}

// GetMe fetches the currently authenticated bot.
//
// This method is called upon creation to validate the token,
// and so you may get this data from BotAPI.Self without the need for
// another request.
func (bot *BotAPI) GetMe() (User, error) {
	return bot.GetMeWithCtx(context.Background())
}

// IsMessageToMe returns true if message directed to this bot.
//
// It requires the Message.
func (bot *BotAPI) IsMessageToMe(message Message) bool {
	return strings.Contains(message.Text, "@"+bot.Self.UserName)
}

func hasFilesNeedingUpload(files []RequestFile) bool {
	for _, file := range files {
		if file.Data.NeedsUpload() {
			return true
		}
	}

	return false
}

func (bot *BotAPI) RequestWithCtx(ctx context.Context, c Chattable) (*APIResponse, error) {
	params, err := c.params()
	if err != nil {
		return nil, err
	}

	if t, ok := c.(Fileable); ok {
		files := t.files()

		// If we have files that need to be uploaded, we should delegate the
		// request to UploadFile.
		if hasFilesNeedingUpload(files) {
			return bot.UploadFilesWithCtx(ctx, t.method(), params, files)
		}

		// However, if there are no files to be uploaded, there's likely things
		// that need to be turned into params instead.
		for _, file := range files {
			params[file.Name] = file.Data.SendData()
		}
	}

	return bot.MakeRequestWithCtx(ctx, c.method(), params)
}

// Request sends a Chattable to Telegram, and returns the APIResponse.
func (bot *BotAPI) Request(c Chattable) (*APIResponse, error) {
	return bot.RequestWithCtx(context.Background(), c)
}

func (bot *BotAPI) requestUnmarshal(ctx context.Context, c Chattable, recv interface{}) error {
	resp, err := bot.RequestWithCtx(ctx, c)
	if err != nil {
		return err
	}
	return json.Unmarshal(resp.Result, &recv)
}

func (bot *BotAPI) SendWithCtx(ctx context.Context, c Chattable) (res Message, err error) {
	err = bot.requestUnmarshal(ctx, c, &res)
	return
}

// Send will send a Chattable item to Telegram and provides the
// returned Message.
func (bot *BotAPI) Send(c Chattable) (Message, error) {
	return bot.SendWithCtx(context.Background(), c)
}

func (bot *BotAPI) SendMediaGroupWithCtx(ctx context.Context, c MediaGroupConfig) (res []Message, err error) {
	err = bot.requestUnmarshal(ctx, c, &res)
	return
}

// SendMediaGroup sends a media group and returns the resulting messages.
func (bot *BotAPI) SendMediaGroup(cfg MediaGroupConfig) ([]Message, error) {
	return bot.SendMediaGroupWithCtx(context.Background(), cfg)
}

// GetUserProfilePhotosWithCtx gets a user's profile photos.
//
// It requires UserID.
// Offset and Limit are optional.
func (bot *BotAPI) GetUserProfilePhotosWithCtx(ctx context.Context, cfg UserProfilePhotosConfig) (res UserProfilePhotos, err error) {
	err = bot.requestUnmarshal(ctx, cfg, &res)
	return
}

func (bot *BotAPI) GetUserProfilePhotos(cfg UserProfilePhotosConfig) (UserProfilePhotos, error) {
	return bot.GetUserProfilePhotosWithCtx(context.Background(), cfg)
}

func (bot *BotAPI) GetFileWithCtx(ctx context.Context, cfg FileConfig) (res File, err error) {
	err = bot.requestUnmarshal(ctx, cfg, &res)
	return
}

// GetFile returns a File which can download a file from Telegram.
//
// Requires FileID.
func (bot *BotAPI) GetFile(cfg FileConfig) (File, error) {
	return bot.GetFileWithCtx(context.Background(), cfg)
}

// GetUpdates fetches updates.
// If a WebHook is set, this will not return any data!
//
// Offset, Limit, Timeout, and AllowedUpdates are optional.
// To avoid stale items, set Offset to one higher than the previous item.
// Set Timeout to a large number to reduce requests, so you can get updates
// instantly instead of having to wait between requests.
func (bot *BotAPI) GetUpdates(config UpdateConfig) ([]Update, error) {
	return bot.GetUpdatesWithCtx(context.Background(), config)
}

func (bot *BotAPI) GetUpdatesWithCtx(ctx context.Context, cfg UpdateConfig) (res []Update, err error) {
	err = bot.requestUnmarshal(ctx, cfg, &res)
	return
}

func (bot *BotAPI) GetWebhookInfoWithCtx(ctx context.Context) (res WebhookInfo, err error) {
	err = bot.requestUnmarshal(ctx, newChattable(nil, "getWebhookInfo"), &res)
	return
}

// GetWebhookInfo allows you to fetch information about a webhook and if
// one currently is set, along with pending update count and error messages.
func (bot *BotAPI) GetWebhookInfo() (WebhookInfo, error) {
	return bot.GetWebhookInfoWithCtx(context.Background())
}

func (bot *BotAPI) GetUpdatesChanWithCtx(ctx context.Context, config UpdateConfig) UpdatesChannel {
	ch := make(chan Update, bot.Buffer)

	go func() {
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			updates, err := bot.GetUpdatesWithCtx(ctx, config)
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

	return ch
}

// GetUpdatesChan starts and returns a channel for getting updates.
func (bot *BotAPI) GetUpdatesChan(config UpdateConfig) UpdatesChannel {
	ctx, cancel := context.WithCancel(context.Background())
	bot.mu.Lock()
	bot.stoppers = append(bot.stoppers, cancel)
	bot.mu.Unlock()

	return bot.GetUpdatesChanWithCtx(ctx, config)
}

// StopReceivingUpdates stops the go routine which receives updates
func (bot *BotAPI) StopReceivingUpdates() {
	bot.mu.Lock()
	defer bot.mu.Unlock()

	for _, stopper := range bot.stoppers {
		stopper()
	}
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

// ListenForWebhookRespReqFormat registers a http handler for a single incoming webhook.
func (bot *BotAPI) ListenForWebhookRespReqFormat(w http.ResponseWriter, r *http.Request) UpdatesChannel {
	ch := make(chan Update, bot.Buffer)

	func(w http.ResponseWriter, r *http.Request) {
		defer close(ch)

		update, err := bot.HandleUpdate(r)
		if err != nil {
			errMsg, _ := json.Marshal(map[string]string{"error": err.Error()})
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(errMsg)
			return
		}

		ch <- *update
	}(w, r)

	return ch
}

// HandleUpdate parses and returns update received via webhook
func (bot *BotAPI) HandleUpdate(r *http.Request) (*Update, error) {
	if r.Method != http.MethodPost {
		err := errors.New("wrong HTTP method required POST")
		return nil, err
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
			return errors.New("unable to use http response to upload files")
		}
	}

	values := buildParams(params)
	values.Set("method", c.method())

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	_, err = w.Write([]byte(values.Encode()))
	return err
}

func (bot *BotAPI) GetChatWithCtx(ctx context.Context, cfg ChatInfoConfig) (res Chat, err error) {
	err = bot.requestUnmarshal(ctx, cfg, &res)
	return
}

// GetChat gets information about a chat.
func (bot *BotAPI) GetChat(cfg ChatInfoConfig) (Chat, error) {
	return bot.GetChatWithCtx(context.Background(), cfg)
}

func (bot *BotAPI) GetChatAdministratorsWithCtx(ctx context.Context, cfg ChatAdministratorsConfig) (res []ChatMember, err error) {
	err = bot.requestUnmarshal(ctx, cfg, &res)
	return
}

// GetChatAdministrators gets a list of administrators in the chat.
//
// If none have been appointed, only the creator will be returned.
// Bots are not shown, even if they are an administrator.
func (bot *BotAPI) GetChatAdministrators(cfg ChatAdministratorsConfig) ([]ChatMember, error) {
	return bot.GetChatAdministratorsWithCtx(context.Background(), cfg)
}

func (bot *BotAPI) GetChatMembersCountWithCtx(ctx context.Context, cfg ChatMemberCountConfig) (res int, err error) {
	err = bot.requestUnmarshal(ctx, cfg, &res)
	return
}

// GetChatMembersCount gets the number of users in a chat.
func (bot *BotAPI) GetChatMembersCount(cfg ChatMemberCountConfig) (int, error) {
	return bot.GetChatMembersCountWithCtx(context.Background(), cfg)
}

func (bot *BotAPI) GetChatMemberWithCtx(ctx context.Context, cfg GetChatMemberConfig) (res ChatMember, err error) {
	err = bot.requestUnmarshal(ctx, cfg, &res)
	return
}

// GetChatMember gets a specific chat member.
func (bot *BotAPI) GetChatMember(cfg GetChatMemberConfig) (ChatMember, error) {
	return bot.GetChatMemberWithCtx(context.Background(), cfg)
}

func (bot *BotAPI) GetGameHighScoresWithContext(ctx context.Context, cfg GetGameHighScoresConfig) (res []GameHighScore, err error) {
	err = bot.requestUnmarshal(ctx, cfg, &res)
	return
}

// GetGameHighScores allows you to get the high scores for a game.
func (bot *BotAPI) GetGameHighScores(cfg GetGameHighScoresConfig) ([]GameHighScore, error) {
	return bot.GetGameHighScoresWithContext(context.Background(), cfg)
}

func (bot *BotAPI) GetInviteLinkWithCtx(ctx context.Context, cfg ChatInviteLinkConfig) (res string, err error) {
	err = bot.requestUnmarshal(ctx, cfg, &res)
	return
}

// GetInviteLink get InviteLink for a chat
func (bot *BotAPI) GetInviteLink(cfg ChatInviteLinkConfig) (string, error) {
	return bot.GetInviteLinkWithCtx(context.Background(), cfg)
}

func (bot *BotAPI) GetStickerSetWithCtx(ctx context.Context, cfg GetStickerSetConfig) (res StickerSet, err error) {
	err = bot.requestUnmarshal(ctx, cfg, &res)
	return
}

// GetStickerSet returns a StickerSet.
func (bot *BotAPI) GetStickerSet(cfg GetStickerSetConfig) (StickerSet, error) {
	return bot.GetStickerSetWithCtx(context.Background(), cfg)
}

func (bot *BotAPI) StopPollWithCtx(ctx context.Context, cfg StopPollConfig) (res Poll, err error) {
	err = bot.requestUnmarshal(ctx, cfg, &res)
	return
}

// StopPoll stops a poll and returns the result.
func (bot *BotAPI) StopPoll(cfg StopPollConfig) (Poll, error) {
	return bot.StopPollWithCtx(context.Background(), cfg)
}

func (bot *BotAPI) GetMyCommandsWithCtx(ctx context.Context) ([]BotCommand, error) {
	return bot.GetMyCommandsWithConfigWithCtx(ctx, GetMyCommandsConfig{})
}

// GetMyCommands gets the currently registered commands.
func (bot *BotAPI) GetMyCommands() ([]BotCommand, error) {
	return bot.GetMyCommandsWithConfig(GetMyCommandsConfig{})
}

func (bot *BotAPI) GetMyCommandsWithConfigWithCtx(ctx context.Context, cfg GetMyCommandsConfig) (res []BotCommand, err error) {
	err = bot.requestUnmarshal(ctx, cfg, &res)
	return
}

// GetMyCommandsWithConfig gets the currently registered commands with a config.
func (bot *BotAPI) GetMyCommandsWithConfig(cfg GetMyCommandsConfig) ([]BotCommand, error) {
	return bot.GetMyCommandsWithConfigWithCtx(context.Background(), cfg)
}

func (bot *BotAPI) CopyMessageWithCtx(ctx context.Context, cfg CopyMessageConfig) (res MessageID, err error) {
	err = bot.requestUnmarshal(ctx, cfg, &res)
	return
}

// CopyMessage copy messages of any kind. The method is analogous to the method
// forwardMessage, but the copied message doesn't have a link to the original
// message. Returns the MessageID of the sent message on success.
func (bot *BotAPI) CopyMessage(cfg CopyMessageConfig) (MessageID, error) {
	return bot.CopyMessageWithCtx(context.Background(), cfg)
}

func (bot *BotAPI) AnswerWebAppQueryWithCtx(ctx context.Context, cfg AnswerWebAppQueryConfig) (res SentWebAppMessage, err error) {
	err = bot.requestUnmarshal(ctx, cfg, &res)
	return
}

// AnswerWebAppQuery sets the result of an interaction with a Web App and send a
// corresponding message on behalf of the user to the chat from which the query originated.
func (bot *BotAPI) AnswerWebAppQuery(cfg AnswerWebAppQueryConfig) (SentWebAppMessage, error) {
	return bot.AnswerWebAppQueryWithCtx(context.Background(), cfg)
}

func (bot *BotAPI) GetMyDefaultAdministratorRightsWithCtx(ctx context.Context, cfg GetMyDefaultAdministratorRightsConfig) (res ChatAdministratorRights, err error) {
	err = bot.requestUnmarshal(ctx, cfg, &res)
	return
}

// GetMyDefaultAdministratorRights gets the current default administrator rights of the bot.
func (bot *BotAPI) GetMyDefaultAdministratorRights(cfg GetMyDefaultAdministratorRightsConfig) (ChatAdministratorRights, error) {
	return bot.GetMyDefaultAdministratorRightsWithCtx(context.Background(), cfg)
}

// EscapeText takes an input text and escape Telegram markup symbols.
// In this way we can send a text without being afraid of having to escape the characters manually.
// Note that you don't have to include the formatting style in the input text, or it will be escaped too.
// If there is an error, an empty string will be returned.
//
// parseMode is the text formatting mode (ModeMarkdown, ModeMarkdownV2 or ModeHTML)
// text is the input string that will be escaped
func EscapeText(parseMode string, text string) string {
	var replacer *strings.Replacer

	if parseMode == ModeHTML {
		replacer = strings.NewReplacer("<", "&lt;", ">", "&gt;", "&", "&amp;")
	} else if parseMode == ModeMarkdown {
		replacer = strings.NewReplacer("_", "\\_", "*", "\\*", "`", "\\`", "[", "\\[")
	} else if parseMode == ModeMarkdownV2 {
		replacer = strings.NewReplacer(
			"_", "\\_", "*", "\\*", "[", "\\[", "]", "\\]", "(",
			"\\(", ")", "\\)", "~", "\\~", "`", "\\`", ">", "\\>",
			"#", "\\#", "+", "\\+", "-", "\\-", "=", "\\=", "|",
			"\\|", "{", "\\{", "}", "\\}", ".", "\\.", "!", "\\!",
		)
	} else {
		return ""
	}

	return replacer.Replace(text)
}
