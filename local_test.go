package tgbotapi

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type APIRequest struct {
	Body string
	Path string
}

func MockTelegramAPI() (*httptest.Server, chan APIRequest, chan error) {
	reqChan := make(chan APIRequest, 1)
	errChan := make(chan error, 1)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			errChan <- err
			return
		}

		apiReq := APIRequest{
			Body: string(body),
			Path: r.URL.Path,
		}
		reqChan <- apiReq
		w.Write([]byte(`{"ok":true}`))
	}))

	return ts, reqChan, errChan
}

func TestSendWithMessageOffline(t *testing.T) {
	ts, reqChan, errChan := MockTelegramAPI()
	defer ts.Close()

	bot := BotAPI{
		APIEndpoint: ts.URL + "/bot%s/%s",
		Client:      http.DefaultClient,
	}

	msg := NewMessage(1, "Hello, World")
	_, err := bot.Send(msg)
	if err != nil {
		t.Fatal(err)
	}

	var actual APIRequest
	expected := APIRequest{
		Body: "chat_id=1&disable_notification=false&disable_web_page_preview=false&text=Hello%2C+World",
		Path: "/bot/sendMessage",
	}

	select {
	case req := <-reqChan:
		actual = req
	case err := <-errChan:
		t.Fatal(err)
	}

	if actual != expected {
		t.Errorf("\nExpected:\n%+v\nActual:\n%+v\n", expected, actual)
	}
}
