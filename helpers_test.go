package tgbotapi

import (
	"testing"
)

func TestNewWebhook(t *testing.T) {
	result, err := NewWebhook("https://example.com/token")

	if err != nil ||
		result.URL.String() != "https://example.com/token" ||
		result.Certificate != interface{}(nil) ||
		result.MaxConnections != 0 ||
		len(result.AllowedUpdates) != 0 {
		t.Fail()
	}
}

func TestNewWebhookWithCert(t *testing.T) {
	exampleFile := FileID("123")
	result, err := NewWebhookWithCert("https://example.com/token", exampleFile)

	if err != nil ||
		result.URL.String() != "https://example.com/token" ||
		result.Certificate != exampleFile ||
		result.MaxConnections != 0 ||
		len(result.AllowedUpdates) != 0 {
		t.Fail()
	}
}

func TestNewInlineQueryResultArticle(t *testing.T) {
	result := NewInlineQueryResultArticle("id", "title", "message")

	if result.Type != "article" ||
		result.ID != "id" ||
		result.Title != "title" ||
		result.InputMessageContent.(InputTextMessageContent).Text != "message" {
		t.Fail()
	}
}

func TestNewInlineQueryResultArticleMarkdown(t *testing.T) {
	result := NewInlineQueryResultArticleMarkdown("id", "title", "*message*")

	if result.Type != "article" ||
		result.ID != "id" ||
		result.Title != "title" ||
		result.InputMessageContent.(InputTextMessageContent).Text != "*message*" ||
		result.InputMessageContent.(InputTextMessageContent).ParseMode != "Markdown" {
		t.Fail()
	}
}

func TestNewInlineQueryResultArticleHTML(t *testing.T) {
	result := NewInlineQueryResultArticleHTML("id", "title", "<b>message</b>")

	if result.Type != "article" ||
		result.ID != "id" ||
		result.Title != "title" ||
		result.InputMessageContent.(InputTextMessageContent).Text != "<b>message</b>" ||
		result.InputMessageContent.(InputTextMessageContent).ParseMode != "HTML" {
		t.Fail()
	}
}

func TestNewInlineQueryResultGIF(t *testing.T) {
	result := NewInlineQueryResultGIF("id", "google.com")

	if result.Type != "gif" ||
		result.ID != "id" ||
		result.URL != "google.com" {
		t.Fail()
	}
}

func TestNewInlineQueryResultMPEG4GIF(t *testing.T) {
	result := NewInlineQueryResultMPEG4GIF("id", "google.com")

	if result.Type != "mpeg4_gif" ||
		result.ID != "id" ||
		result.URL != "google.com" {
		t.Fail()
	}
}

func TestNewInlineQueryResultPhoto(t *testing.T) {
	result := NewInlineQueryResultPhoto("id", "google.com")

	if result.Type != "photo" ||
		result.ID != "id" ||
		result.URL != "google.com" {
		t.Fail()
	}
}

func TestNewInlineQueryResultPhotoWithThumb(t *testing.T) {
	result := NewInlineQueryResultPhotoWithThumb("id", "google.com", "thumb.com")

	if result.Type != "photo" ||
		result.ID != "id" ||
		result.URL != "google.com" ||
		result.ThumbURL != "thumb.com" {
		t.Fail()
	}
}

func TestNewInlineQueryResultVideo(t *testing.T) {
	result := NewInlineQueryResultVideo("id", "google.com")

	if result.Type != "video" ||
		result.ID != "id" ||
		result.URL != "google.com" {
		t.Fail()
	}
}

func TestNewInlineQueryResultAudio(t *testing.T) {
	result := NewInlineQueryResultAudio("id", "google.com", "title")

	if result.Type != "audio" ||
		result.ID != "id" ||
		result.URL != "google.com" ||
		result.Title != "title" {
		t.Fail()
	}
}

func TestNewInlineQueryResultVoice(t *testing.T) {
	result := NewInlineQueryResultVoice("id", "google.com", "title")

	if result.Type != "voice" ||
		result.ID != "id" ||
		result.URL != "google.com" ||
		result.Title != "title" {
		t.Fail()
	}
}

func TestNewInlineQueryResultDocument(t *testing.T) {
	result := NewInlineQueryResultDocument("id", "google.com", "title", "mime/type")

	if result.Type != "document" ||
		result.ID != "id" ||
		result.URL != "google.com" ||
		result.Title != "title" ||
		result.MimeType != "mime/type" {
		t.Fail()
	}
}

func TestNewInlineQueryResultLocation(t *testing.T) {
	result := NewInlineQueryResultLocation("id", "name", 40, 50)

	if result.Type != "location" ||
		result.ID != "id" ||
		result.Title != "name" ||
		result.Latitude != 40 ||
		result.Longitude != 50 {
		t.Fail()
	}
}

func TestNewInlineKeyboardButtonLoginURL(t *testing.T) {
	result := NewInlineKeyboardButtonLoginURL("text", LoginURL{
		URL:                "url",
		ForwardText:        "ForwardText",
		BotUsername:        "username",
		RequestWriteAccess: false,
	})

	if result.Text != "text" ||
		result.LoginURL.URL != "url" ||
		result.LoginURL.ForwardText != "ForwardText" ||
		result.LoginURL.BotUsername != "username" ||
		result.LoginURL.RequestWriteAccess != false {
		t.Fail()
	}
}

func TestNewEditMessageText(t *testing.T) {
	edit := NewEditMessageText(ChatID, ReplyToMessageID, "new text")

	if edit.Text != "new text" ||
		edit.BaseEdit.ChatID != ChatID ||
		edit.BaseEdit.MessageID != ReplyToMessageID {
		t.Fail()
	}
}

func TestNewEditMessageCaption(t *testing.T) {
	edit := NewEditMessageCaption(ChatID, ReplyToMessageID, "new caption")

	if edit.Caption != "new caption" ||
		edit.BaseEdit.ChatID != ChatID ||
		edit.BaseEdit.MessageID != ReplyToMessageID {
		t.Fail()
	}
}

func TestNewEditMessageReplyMarkup(t *testing.T) {
	markup := InlineKeyboardMarkup{
		InlineKeyboard: [][]InlineKeyboardButton{
			{
				{Text: "test"},
			},
		},
	}

	edit := NewEditMessageReplyMarkup(ChatID, ReplyToMessageID, markup)

	if edit.ReplyMarkup.InlineKeyboard[0][0].Text != "test" ||
		edit.BaseEdit.ChatID != ChatID ||
		edit.BaseEdit.MessageID != ReplyToMessageID {
		t.Fail()
	}

}

func TestNewDice(t *testing.T) {
	dice := NewDice(42)

	if dice.ChatID != 42 ||
		dice.Emoji != "" {
		t.Fail()
	}
}

func TestNewDiceWithEmoji(t *testing.T) {
	dice := NewDiceWithEmoji(42, "🏀")

	if dice.ChatID != 42 ||
		dice.Emoji != "🏀" {
		t.Fail()
	}
}

func TestValidateWebAppData(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		token := "5473903189:AAFnHnISQMP5UQQ5MEaoEWvxeiwNgz2CN2U"
		initData := "query_id=AAG1bpMJAAAAALVukwmZ_H2t&user=%7B%22id%22%3A160657077%2C%22first_name%22%3A%22Yury%20R%22%2C%22last_name%22%3A%22%22%2C%22username%22%3A%22crashiura%22%2C%22language_code%22%3A%22en%22%7D&auth_date=1656804462&hash=8d6960760a573d3212deb05e20d1a34959c83d24c1bc44bb26dde49a42aa9b34"
		result, err := ValidateWebAppData(token, initData)
		if err != nil {
			t.Fail()
		}
		if !result {
			t.Fail()
		}
	})

	t.Run("error", func(t *testing.T) {
		token := "5473903189:AAFnHnISQMP5UQQ5MEaoEWvxeiwNgz2CN2U"
		initData := "asdfasdfasdfasdfasdf"
		result, err := ValidateWebAppData(token, initData)
		if err == nil {
			t.Fail()
		}
		if result {
			t.Fail()
		}
	})
}
