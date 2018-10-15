package tgbotapi_test

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"testing"
)

func TestNewInlineQueryResultArticle(t *testing.T) {
	result := tgbotapi.NewInlineQueryResultArticle("id", "title", "message")

	if result.Type != "article" ||
		result.ID != "id" ||
		result.Title != "title" ||
		result.InputMessageContent.(tgbotapi.InputTextMessageContent).Text != "message" {
		t.Fail()
	}
}

func TestNewInlineQueryResultArticleMarkdown(t *testing.T) {
	result := tgbotapi.NewInlineQueryResultArticleMarkdown("id", "title", "*message*")

	if result.Type != "article" ||
		result.ID != "id" ||
		result.Title != "title" ||
		result.InputMessageContent.(tgbotapi.InputTextMessageContent).Text != "*message*" ||
		result.InputMessageContent.(tgbotapi.InputTextMessageContent).ParseMode != "Markdown" {
		t.Fail()
	}
}

func TestNewInlineQueryResultArticleHTML(t *testing.T) {
	result := tgbotapi.NewInlineQueryResultArticleHTML("id", "title", "<b>message</b>")

	if result.Type != "article" ||
		result.ID != "id" ||
		result.Title != "title" ||
		result.InputMessageContent.(tgbotapi.InputTextMessageContent).Text != "<b>message</b>" ||
		result.InputMessageContent.(tgbotapi.InputTextMessageContent).ParseMode != "HTML" {
		t.Fail()
	}
}

func TestNewInlineQueryResultGIF(t *testing.T) {
	result := tgbotapi.NewInlineQueryResultGIF("id", "google.com")

	if result.Type != "gif" ||
		result.ID != "id" ||
		result.URL != "google.com" {
		t.Fail()
	}
}

func TestNewInlineQueryResultMPEG4GIF(t *testing.T) {
	result := tgbotapi.NewInlineQueryResultMPEG4GIF("id", "google.com")

	if result.Type != "mpeg4_gif" ||
		result.ID != "id" ||
		result.URL != "google.com" {
		t.Fail()
	}
}

func TestNewInlineQueryResultPhoto(t *testing.T) {
	result := tgbotapi.NewInlineQueryResultPhoto("id", "google.com")

	if result.Type != "photo" ||
		result.ID != "id" ||
		result.URL != "google.com" {
		t.Fail()
	}
}

func TestNewInlineQueryResultPhotoWithThumb(t *testing.T) {
	result := tgbotapi.NewInlineQueryResultPhotoWithThumb("id", "google.com", "thumb.com")

	if result.Type != "photo" ||
		result.ID != "id" ||
		result.URL != "google.com" ||
		result.ThumbURL != "thumb.com" {
		t.Fail()
	}
}

func TestNewInlineQueryResultVideo(t *testing.T) {
	result := tgbotapi.NewInlineQueryResultVideo("id", "google.com")

	if result.Type != "video" ||
		result.ID != "id" ||
		result.URL != "google.com" {
		t.Fail()
	}
}

func TestNewInlineQueryResultAudio(t *testing.T) {
	result := tgbotapi.NewInlineQueryResultAudio("id", "google.com", "title")

	if result.Type != "audio" ||
		result.ID != "id" ||
		result.URL != "google.com" ||
		result.Title != "title" {
		t.Fail()
	}
}

func TestNewInlineQueryResultVoice(t *testing.T) {
	result := tgbotapi.NewInlineQueryResultVoice("id", "google.com", "title")

	if result.Type != "voice" ||
		result.ID != "id" ||
		result.URL != "google.com" ||
		result.Title != "title" {
		t.Fail()
	}
}

func TestNewInlineQueryResultDocument(t *testing.T) {
	result := tgbotapi.NewInlineQueryResultDocument("id", "google.com", "title", "mime/type")

	if result.Type != "document" ||
		result.ID != "id" ||
		result.URL != "google.com" ||
		result.Title != "title" ||
		result.MimeType != "mime/type" {
		t.Fail()
	}
}

func TestNewInlineQueryResultLocation(t *testing.T) {
	result := tgbotapi.NewInlineQueryResultLocation("id", "name", 40, 50)

	if result.Type != "location" ||
		result.ID != "id" ||
		result.Title != "name" ||
		result.Latitude != 40 ||
		result.Longitude != 50 {
		t.Fail()
	}
}

func TestNewEditMessageText(t *testing.T) {
	edit := tgbotapi.NewEditMessageText(ChatID, ReplyToMessageID, "new text")

	if edit.Text != "new text" ||
		edit.BaseEdit.ChatID != ChatID ||
		edit.BaseEdit.MessageID != ReplyToMessageID {
		t.Fail()
	}
}

func TestNewEditMessageCaption(t *testing.T) {
	edit := tgbotapi.NewEditMessageCaption(ChatID, ReplyToMessageID, "new caption")

	if edit.Caption != "new caption" ||
		edit.BaseEdit.ChatID != ChatID ||
		edit.BaseEdit.MessageID != ReplyToMessageID {
		t.Fail()
	}
}

func TestNewEditMessageReplyMarkup(t *testing.T) {
	markup := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.InlineKeyboardButton{Text: "test"},
			},
		},
	}

	edit := tgbotapi.NewEditMessageReplyMarkup(ChatID, ReplyToMessageID, markup)

	if edit.ReplyMarkup.InlineKeyboard[0][0].Text != "test" ||
		edit.BaseEdit.ChatID != ChatID ||
		edit.BaseEdit.MessageID != ReplyToMessageID {
		t.Fail()
	}

}
