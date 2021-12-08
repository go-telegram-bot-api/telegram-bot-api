package tgbotapi

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	APIToken = "123abc"

	BotID int64 = 12375639
)

type TestLogger struct {
	t *testing.T
}

func (logger TestLogger) Println(v ...interface{}) {
	logger.t.Log(v...)
}

func (logger TestLogger) Printf(format string, v ...interface{}) {
	logger.t.Logf(format, v...)
}

func getServerURL(tsURL string) string {
	return tsURL + "/bot%s/%s"
}

// BotRequestFile is a file that should be included as part of a request to the
// Telegram Bot API.
type BotRequestFile struct {
	// Name is the name of the field.
	Name string
	// Data is the contents of the file.
	Data []byte
}

// BotRequest is information about a mocked request to the Telegram Bot API.
type BotRequest struct {
	// Endpoint is the method that should be called on the Bot API.
	Endpoint string
	// RequestData is the data that is sent to the endpoint.
	RequestData url.Values

	// Files are files that should exist as part of the request. If this is nil,
	// it is assumed the request should be treated as x-www-form-urlencoded,
	// otherwise it will be treated as multipart data.
	Files []BotRequestFile

	// ResponseData is the response body.
	ResponseData string
}

var getMeRequest = BotRequest{
	Endpoint:     "getMe",
	ResponseData: fmt.Sprintf(`{"ok": true, "result": {"id": %d, "is_bot": true, "first_name": "Test", "username": "test_bot", "can_join_groups": true, "can_read_all_group_messages": false, "supports_inline_queries": true}}`, BotID),
}

func assertBotRequests(t *testing.T, token string, botID int64, requests []BotRequest) (*httptest.Server, *BotAPI) {
	logger := TestLogger{t}
	SetLogger(logger)

	currentRequest := 0

	lock := sync.Mutex{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lock.Lock()
		defer lock.Unlock()

		t.Logf("Got request to %s", r.URL.String())

		assert.Greater(t, len(requests), currentRequest, "got more requests than were configured")

		nextRequest := requests[currentRequest]
		expected := fmt.Sprintf("/bot%s/%s", token, nextRequest.Endpoint)
		t.Logf("Expecting request to %s", expected)

		url := r.URL.String()
		t.Logf("Got request to %s", url)

		switch url {
		case expected:
			if nextRequest.Files == nil {
				assert.NoError(t, r.ParseForm())
				assert.Equal(t, len(nextRequest.RequestData), len(r.Form), "request must have same number of values")

				for expectedValueKey, expectedValue := range nextRequest.RequestData {
					t.Logf("Checking if %s contains %v", expectedValueKey, expectedValue)

					assert.Len(t, expectedValue, 1, "each expected key should only have one value")

					foundValue := r.Form[expectedValueKey]
					t.Logf("Form contains %+v", foundValue)

					assert.Len(t, foundValue, 1, "each key should have exactly one value")
					assert.Equal(t, expectedValue[0], foundValue[0])
				}
			} else if nextRequest.Files != nil {
				assert.NoError(t, r.ParseMultipartForm(1024*1024*50), "request must be valid multipart form")
				assert.Equal(t, len(nextRequest.RequestData), len(r.MultipartForm.Value), "request must have correct number of values")
				assert.Equal(t, len(nextRequest.Files), len(r.MultipartForm.File), "request must have correct number of files")

				for expectedValueKey, expectedValue := range nextRequest.RequestData {
					t.Logf("Checking if %s contains %v", expectedValueKey, expectedValue)

					assert.Len(t, expectedValue, 1, "each expected key should only have one value")

					foundValue := r.MultipartForm.Value[expectedValueKey]
					assert.Len(t, foundValue, 1, "each key should have exactly one value")
					assert.Equal(t, expectedValue[0], foundValue[0])
				}

				for _, expectedFile := range nextRequest.Files {
					t.Logf("Checking if %s is set correctly", expectedFile.Name)

					foundFile := r.MultipartForm.File[expectedFile.Name]
					assert.Len(t, foundFile, 1, "each file must appear exactly once")

					f, err := foundFile[0].Open()
					assert.NoError(t, err, "should be able to open file")
					data, err := io.ReadAll(f)
					assert.NoError(t, err, "should be able to read file")
					assert.Equal(t, expectedFile.Data, data, "uploaded file should be the same")
				}
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, nextRequest.ResponseData)
		default:
			t.Errorf("expected request to %s but got request to %s", expected, url)
		}

		currentRequest += 1
		if currentRequest > len(requests) {
			t.Errorf("expected %d requests but got %d", len(requests), currentRequest)
		}
	}))

	bot, err := NewBotAPIWithAPIEndpoint(APIToken, getServerURL(ts.URL))
	if err != nil {
		t.Error(err)
	}

	bot.Debug = true
	return ts, bot
}

func TestNewBotAPI(t *testing.T) {
	ts, bot := assertBotRequests(t, APIToken, BotID, []BotRequest{getMeRequest})
	defer ts.Close()

	assert.Equal(t, BotID, bot.Self.ID, "Bot ID should be expected ID")
}

func TestMakeRequest(t *testing.T) {
	values := url.Values{}
	values.Set("param_name", "param_value")

	goodReq := BotRequest{
		Endpoint:    "testEndpoint",
		RequestData: values,

		ResponseData: `{"ok": true, "result": true}`,
	}

	badReq := BotRequest{
		Endpoint:    "testEndpoint",
		RequestData: values,

		ResponseData: `{"ok": false, "description": "msg", "error_code": 12343}`,
	}

	badReqWithParams := BotRequest{
		Endpoint:    "testEndpoint",
		RequestData: values,

		ResponseData: `{"ok": false, "description": "msg", "error_code": 12343, "parameters": {"retry_after": 52, "migrate_to_chat_id": 245}}`,
	}

	ts, bot := assertBotRequests(t, APIToken, BotID, []BotRequest{getMeRequest, goodReq, badReq, badReqWithParams})
	defer ts.Close()

	params := make(Params)
	params["param_name"] = "param_value"

	resp, err := bot.MakeRequest("testEndpoint", params)
	assert.NoError(t, err, "bot should be able to make request without errors")
	assert.Equal(t, &APIResponse{
		Ok:     true,
		Result: []byte("true"),
	}, resp)

	resp, err = bot.MakeRequest("testEndpoint", params)
	assert.Error(t, err)
	assert.Equal(t, &Error{
		Code:    12343,
		Message: "msg",
	}, err)
	assert.NotNil(t, resp)

	resp, err = bot.MakeRequest("testEndpoint", params)
	assert.Error(t, err)
	assert.Equal(t, &Error{
		Code:    12343,
		Message: "msg",
		ResponseParameters: ResponseParameters{
			MigrateToChatID: 245,
			RetryAfter:      52,
		},
	}, err)
	assert.NotNil(t, resp)
	assert.Equal(t, &ResponseParameters{
		MigrateToChatID: 245,
		RetryAfter:      52,
	}, resp.Parameters)
}

func TestUploadFilesBasic(t *testing.T) {
	values := url.Values{}
	values.Set("param_name", "param_value")

	goodReq := BotRequest{
		Endpoint:    "testEndpoint",
		RequestData: values,
		Files: []BotRequestFile{{
			Name: "file1",
			Data: []byte("data1"),
		}},

		ResponseData: `{"ok": true, "result": true}`,
	}

	badReq := BotRequest{
		Endpoint:    "testEndpoint",
		RequestData: values,
		Files: []BotRequestFile{{
			Name: "file1",
			Data: []byte("data1"),
		}},

		ResponseData: `{"ok": false, "description": "msg", "error_code": 12343}`,
	}

	ts, bot := assertBotRequests(t, APIToken, BotID, []BotRequest{getMeRequest, goodReq, badReq})
	defer ts.Close()

	params := make(Params)
	params["param_name"] = "param_value"

	files := []RequestFile{{
		Name: "file1",
		Data: FileBytes{Name: "file1", Bytes: []byte("data1")},
	}}

	resp, err := bot.UploadFiles("testEndpoint", params, files)
	assert.NoError(t, err, "bot should be able to make request without errors")
	assert.Equal(t, &APIResponse{
		Ok:     true,
		Result: []byte("true"),
	}, resp)

	resp, err = bot.UploadFiles("testEndpoint", params, files)
	assert.Error(t, err)
	assert.Equal(t, &Error{
		Code:    12343,
		Message: "msg",
	}, err)
	assert.NotNil(t, resp)
}

func TestUploadFilesAllTypes(t *testing.T) {
	values := url.Values{}
	values.Set("param_name", "param_value")
	values.Set("file-url", "url")
	values.Set("file-id", "id")
	values.Set("file-attach", "attach")

	req := BotRequest{
		Endpoint:    "uploadFiles",
		RequestData: values,
		Files: []BotRequestFile{
			{Name: "file-bytes", Data: []byte("byte-data")},
			{Name: "file-reader", Data: []byte("reader-data")},
			{Name: "file-path", Data: []byte("path-data\n")},
		},

		ResponseData: `{"ok": true, "result": true}`,
	}

	ts, bot := assertBotRequests(t, APIToken, BotID, []BotRequest{getMeRequest, req})
	defer ts.Close()

	params := make(Params)
	params["param_name"] = "param_value"

	files := []RequestFile{{
		Name: "file-bytes",
		Data: FileBytes{
			Name:  "file-bytes-name",
			Bytes: []byte("byte-data"),
		},
	}, {
		Name: "file-reader",
		Data: FileReader{
			Name:   "file-reader-name",
			Reader: bytes.NewReader([]byte("reader-data")),
		},
	}, {
		Name: "file-path",
		Data: FilePath("tests/file-path"),
	}, {
		Name: "file-url",
		Data: FileURL("url"),
	}, {
		Name: "file-id",
		Data: FileID("id"),
	}, {
		Name: "file-attach",
		Data: fileAttach("attach"),
	}}

	resp, err := bot.UploadFiles("uploadFiles", params, files)
	assert.NoError(t, err, "bot should be able to make request without errors")
	assert.Equal(t, &APIResponse{
		Ok:     true,
		Result: []byte("true"),
	}, resp)
}

func TestGetMe(t *testing.T) {
	goodReq := BotRequest{
		Endpoint:     "getMe",
		ResponseData: `{"ok": true, "result": {}}`,
	}

	badReq := BotRequest{
		Endpoint:     "getMe",
		ResponseData: `{"ok": false}`,
	}

	ts, bot := assertBotRequests(t, APIToken, BotID, []BotRequest{getMeRequest, goodReq, badReq})
	defer ts.Close()

	_, err := bot.GetMe()
	assert.NoError(t, err)

	_, err = bot.GetMe()
	assert.Error(t, err)
}

func TestIsMessageToMe(t *testing.T) {
	testCases := []struct {
		Text    string
		Caption string

		Entities        []MessageEntity
		CaptionEntities []MessageEntity

		IsMention bool
	}{{
		Text:      "asdf",
		Entities:  []MessageEntity{},
		IsMention: false,
	}, {
		Text: "@test_bot",
		Entities: []MessageEntity{{
			Type:   "mention",
			Offset: 0,
			Length: 9,
		}},
		IsMention: true,
	}, {
		Text: "prefix @test_bot suffix",
		Entities: []MessageEntity{{
			Type:   "mention",
			Offset: 7,
			Length: 9,
		}},
		IsMention: true,
	}, {
		Text: "prefix @test_bot suffix",
		Entities: []MessageEntity{{
			Type:   "link",
			Offset: 7,
			Length: 9,
		}},
		IsMention: false,
	}, {
		Text: "prefix @test_bot suffix",
		Entities: []MessageEntity{{
			Type:   "link",
			Offset: 0,
			Length: 6,
		}, {
			Type:   "mention",
			Offset: 7,
			Length: 9,
		}},
		IsMention: true,
	}, {
		Text:      "prefix @test_bot suffix",
		IsMention: false,
	}, {
		Caption: "prefix @test_bot suffix",
		CaptionEntities: []MessageEntity{{
			Type:   "mention",
			Offset: 7,
			Length: 9,
		}},
		IsMention: true,
	}}

	bot := BotAPI{
		Self: User{
			UserName: "test_bot",
		},
	}

	for _, test := range testCases {
		assert.Equal(t, test.IsMention, bot.IsMessageToMe(Message{
			Text:            test.Text,
			Caption:         test.Caption,
			Entities:        test.Entities,
			CaptionEntities: test.CaptionEntities,
		}))
	}
}

func TestRequest(t *testing.T) {
	values1 := url.Values{}
	values1.Set("chat_id", "12356")
	values1.Set("thumb", "url")

	req1 := BotRequest{
		Endpoint:    "sendPhoto",
		RequestData: values1,

		Files: []BotRequestFile{{
			Name: "photo",
			Data: []byte("photo-data"),
		}},
		ResponseData: `{"ok": true, "result": {"message_id": "asdf"}}`,
	}

	values2 := url.Values{}
	values2.Set("chat_id", "12356")
	values2.Set("photo", "id")

	req2 := BotRequest{
		Endpoint:     "sendPhoto",
		RequestData:  values2,
		ResponseData: `{"ok": true, "result": {"message_id": 123}}`,
	}

	ts, bot := assertBotRequests(t, APIToken, BotID, []BotRequest{getMeRequest, req1, req2})
	defer ts.Close()

	resp, err := bot.Request(PhotoConfig{
		BaseFile: BaseFile{
			BaseChat: BaseChat{
				ChatID: 12356,
			},
			File: FileBytes{
				Name:  "photo.jpg",
				Bytes: []byte("photo-data"),
			},
		},
		Thumb: FileURL("url"),
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	resp, err = bot.Request(PhotoConfig{
		BaseFile: BaseFile{
			BaseChat: BaseChat{
				ChatID: 12356,
			},
			File: FileID("id"),
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestSend(t *testing.T) {
	values := url.Values{}
	values.Set("chat_id", "12356")
	values.Set("text", "text")

	req1 := BotRequest{
		Endpoint:     "sendMessage",
		RequestData:  values,
		ResponseData: `{"ok": true, "result": {"message_id": "asdf"}}`,
	}

	req2 := BotRequest{
		Endpoint:     "sendMessage",
		RequestData:  values,
		ResponseData: `{"ok": true, "result": {"message_id": 123}}`,
	}

	ts, bot := assertBotRequests(t, APIToken, BotID, []BotRequest{getMeRequest, req1, req2})
	defer ts.Close()

	msg, err := bot.Send(MessageConfig{
		BaseChat: BaseChat{
			ChatID: 12356,
		},
		Text: "text",
	})
	assert.Error(t, err)
	assert.Empty(t, msg)

	msg, err = bot.Send(MessageConfig{
		BaseChat: BaseChat{
			ChatID: 12356,
		},
		Text: "text",
	})
	assert.NoError(t, err)
	assert.Equal(t, Message{MessageID: 123}, msg)
}

func TestSendMediaGroup(t *testing.T) {
	values := url.Values{}
	values.Set("chat_id", "125")
	values.Set("media", `[{"type":"photo","media":"attach://file-0"},{"type":"photo","media":"file-id"}]`)

	req1 := BotRequest{
		Endpoint:    "sendMediaGroup",
		RequestData: values,
		Files: []BotRequestFile{{
			Name: "file-0",
			Data: []byte("path-data\n"),
		}},
		ResponseData: `{"ok": true, "result": [{"message_id": 123643}, {"message_id": 53452}]}`,
	}

	ts, bot := assertBotRequests(t, APIToken, BotID, []BotRequest{getMeRequest, req1})
	defer ts.Close()

	group, err := bot.SendMediaGroup(MediaGroupConfig{})
	assert.ErrorIs(t, err, ErrEmptyMediaGroup)
	assert.Nil(t, group)

	group, err = bot.SendMediaGroup(MediaGroupConfig{
		ChatID: 125,
		Media: []interface{}{
			NewInputMediaPhoto(FilePath("tests/file-path")),
			NewInputMediaPhoto(FileID("file-id")),
		},
	})
	assert.NoError(t, err)
	assert.Len(t, group, 2)
}

func TestGetUserProfilePhotos(t *testing.T) {
	values := url.Values{}
	values.Set("user_id", "5426")

	req1 := BotRequest{
		Endpoint:     "getUserProfilePhotos",
		RequestData:  values,
		ResponseData: `{"ok": true, "result": {"total_count": 24, "photos": [[{"file_id": "abc123-file"}]]}}`,
	}

	ts, bot := assertBotRequests(t, APIToken, BotID, []BotRequest{getMeRequest, req1})
	defer ts.Close()

	profilePhotos, err := bot.GetUserProfilePhotos(UserProfilePhotosConfig{
		UserID: 5426,
	})
	assert.NoError(t, err)
	assert.Equal(t, UserProfilePhotos{
		TotalCount: 24,
		Photos: [][]PhotoSize{
			{{
				FileID: "abc123-file",
			}},
		},
	}, profilePhotos)
}
