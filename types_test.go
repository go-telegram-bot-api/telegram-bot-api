package tgbotapi_test

import (
	"testing"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func TestUserStringWith(t *testing.T) {
	user := tgbotapi.User{
		ID:           0,
		FirstName:    "Test",
		LastName:     "Test",
		UserName:     "",
		LanguageCode: "en",
		IsBot:        false,
	}

	if user.String() != "Test Test" {
		t.Fail()
	}
}

func TestUserStringWithUserName(t *testing.T) {
	user := tgbotapi.User{
		ID:           0,
		FirstName:    "Test",
		LastName:     "Test",
		UserName:     "@test",
		LanguageCode: "en",
	}

	if user.String() != "@test" {
		t.Fail()
	}
}

func TestMessageTime(t *testing.T) {
	message := tgbotapi.Message{Date: 0}

	date := time.Unix(0, 0)
	if message.Time() != date {
		t.Fail()
	}
}

func TestMessageIsCommandWithCommand(t *testing.T) {
	message := tgbotapi.Message{Text: "/command"}
	message.Entities = []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 8}}

	if !message.IsCommand() {
		t.Fail()
	}
}

func TestIsCommandWithText(t *testing.T) {
	message := tgbotapi.Message{Text: "some text"}

	if message.IsCommand() {
		t.Fail()
	}
}

func TestIsCommandWithEmptyText(t *testing.T) {
	message := tgbotapi.Message{Text: ""}

	if message.IsCommand() {
		t.Fail()
	}
}

func TestCommandWithCommand(t *testing.T) {
	message := tgbotapi.Message{Text: "/command"}
	message.Entities = []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 8}}

	if message.Command() != "command" {
		t.Fail()
	}
}

func TestCommandWithEmptyText(t *testing.T) {
	message := tgbotapi.Message{Text: ""}

	if message.Command() != "" {
		t.Fail()
	}
}

func TestCommandWithNonCommand(t *testing.T) {
	message := tgbotapi.Message{Text: "test text"}

	if message.Command() != "" {
		t.Fail()
	}
}

func TestCommandWithBotName(t *testing.T) {
	message := tgbotapi.Message{Text: "/command@testbot"}
	message.Entities = []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 16}}

	if message.Command() != "command" {
		t.Fail()
	}
}

func TestCommandWithAtWithBotName(t *testing.T) {
	message := tgbotapi.Message{Text: "/command@testbot"}
	message.Entities = []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 16}}

	if message.CommandWithAt() != "command@testbot" {
		t.Fail()
	}
}

func TestMessageCommandArgumentsWithArguments(t *testing.T) {
	message := tgbotapi.Message{Text: "/command with arguments"}
	message.Entities = []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 8}}
	if message.CommandArguments() != "with arguments" {
		t.Fail()
	}
}

func TestMessageCommandArgumentsWithMalformedArguments(t *testing.T) {
	message := tgbotapi.Message{Text: "/command-without argument space"}
	message.Entities = []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 8}}
	if message.CommandArguments() != "without argument space" {
		t.Fail()
	}
}

func TestMessageCommandArgumentsWithoutArguments(t *testing.T) {
	message := tgbotapi.Message{Text: "/command"}
	if message.CommandArguments() != "" {
		t.Fail()
	}
}

func TestMessageCommandArgumentsForNonCommand(t *testing.T) {
	message := tgbotapi.Message{Text: "test text"}
	if message.CommandArguments() != "" {
		t.Fail()
	}
}

func TestMessageEntityParseURLGood(t *testing.T) {
	entity := tgbotapi.MessageEntity{URL: "https://www.google.com"}

	if _, err := entity.ParseURL(); err != nil {
		t.Fail()
	}
}

func TestMessageEntityParseURLBad(t *testing.T) {
	entity := tgbotapi.MessageEntity{URL: ""}

	if _, err := entity.ParseURL(); err == nil {
		t.Fail()
	}
}

func TestChatIsPrivate(t *testing.T) {
	chat := tgbotapi.Chat{ID: 10, Type: "private"}

	if !chat.IsPrivate() {
		t.Fail()
	}
}

func TestChatIsGroup(t *testing.T) {
	chat := tgbotapi.Chat{ID: 10, Type: "group"}

	if !chat.IsGroup() {
		t.Fail()
	}
}

func TestChatIsChannel(t *testing.T) {
	chat := tgbotapi.Chat{ID: 10, Type: "channel"}

	if !chat.IsChannel() {
		t.Fail()
	}
}

func TestChatIsSuperGroup(t *testing.T) {
	chat := tgbotapi.Chat{ID: 10, Type: "supergroup"}

	if !chat.IsSuperGroup() {
		t.Fail()
	}
}

func TestFileLink(t *testing.T) {
	file := tgbotapi.File{FilePath: "test/test.txt"}

	if file.Link("token") != "https://api.telegram.org/file/bottoken/test/test.txt" {
		t.Fail()
	}
}

// Ensure all configs are sendable
var (
	_ tgbotapi.Chattable = tgbotapi.AnimationConfig{}
	_ tgbotapi.Chattable = tgbotapi.AudioConfig{}
	_ tgbotapi.Chattable = tgbotapi.CallbackConfig{}
	_ tgbotapi.Chattable = tgbotapi.ChatActionConfig{}
	_ tgbotapi.Chattable = tgbotapi.ContactConfig{}
	_ tgbotapi.Chattable = tgbotapi.DeleteChatPhotoConfig{}
	_ tgbotapi.Chattable = tgbotapi.DeleteChatStickerSetConfig{}
	_ tgbotapi.Chattable = tgbotapi.DeleteMessageConfig{}
	_ tgbotapi.Chattable = tgbotapi.DocumentConfig{}
	_ tgbotapi.Chattable = tgbotapi.EditMessageCaptionConfig{}
	_ tgbotapi.Chattable = tgbotapi.EditMessageLiveLocationConfig{}
	_ tgbotapi.Chattable = tgbotapi.EditMessageReplyMarkupConfig{}
	_ tgbotapi.Chattable = tgbotapi.EditMessageTextConfig{}
	_ tgbotapi.Chattable = tgbotapi.ForwardConfig{}
	_ tgbotapi.Chattable = tgbotapi.GameConfig{}
	_ tgbotapi.Chattable = tgbotapi.GetGameHighScoresConfig{}
	_ tgbotapi.Chattable = tgbotapi.InlineConfig{}
	_ tgbotapi.Chattable = tgbotapi.InvoiceConfig{}
	_ tgbotapi.Chattable = tgbotapi.KickChatMemberConfig{}
	_ tgbotapi.Chattable = tgbotapi.LocationConfig{}
	_ tgbotapi.Chattable = tgbotapi.MediaGroupConfig{}
	_ tgbotapi.Chattable = tgbotapi.MessageConfig{}
	_ tgbotapi.Chattable = tgbotapi.PhotoConfig{}
	_ tgbotapi.Chattable = tgbotapi.PinChatMessageConfig{}
	_ tgbotapi.Chattable = tgbotapi.SetChatDescriptionConfig{}
	_ tgbotapi.Chattable = tgbotapi.SetChatPhotoConfig{}
	_ tgbotapi.Chattable = tgbotapi.SetChatTitleConfig{}
	_ tgbotapi.Chattable = tgbotapi.SetGameScoreConfig{}
	_ tgbotapi.Chattable = tgbotapi.StickerConfig{}
	_ tgbotapi.Chattable = tgbotapi.UnpinChatMessageConfig{}
	_ tgbotapi.Chattable = tgbotapi.UpdateConfig{}
	_ tgbotapi.Chattable = tgbotapi.UserProfilePhotosConfig{}
	_ tgbotapi.Chattable = tgbotapi.VenueConfig{}
	_ tgbotapi.Chattable = tgbotapi.VideoConfig{}
	_ tgbotapi.Chattable = tgbotapi.VideoNoteConfig{}
	_ tgbotapi.Chattable = tgbotapi.VoiceConfig{}
	_ tgbotapi.Chattable = tgbotapi.WebhookConfig{}
)
