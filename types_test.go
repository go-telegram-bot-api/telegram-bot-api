package tgbotapi_test

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"testing"
	"time"
)

func TestUserStringWith(t *testing.T) {
	user := tgbotapi.User{0, "Test", "Test", ""}

	if user.String() != "Test Test" {
		t.Fail()
	}
}

func TestUserStringWithUserName(t *testing.T) {
	user := tgbotapi.User{0, "Test", "Test", "@test"}

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

	if message.IsCommand() != true {
		t.Fail()
	}
}

func TestIsCommandWithText(t *testing.T) {
	message := tgbotapi.Message{Text: "some text"}

	if message.IsCommand() != false {
		t.Fail()
	}
}

func TestIsCommandWithEmptyText(t *testing.T) {
	message := tgbotapi.Message{Text: ""}

	if message.IsCommand() != false {
		t.Fail()
	}
}

func TestCommandWithCommand(t *testing.T) {
	message := tgbotapi.Message{Text: "/command"}

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

	if message.Command() != "command" {
		t.Fail()
	}
}

func TestMessageCommandArgumentsWithArguments(t *testing.T) {
	message := tgbotapi.Message{Text: "/command with arguments"}
	if message.CommandArguments() != "with arguments" {
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

	if chat.IsPrivate() != true {
		t.Fail()
	}
}

func TestChatIsGroup(t *testing.T) {
	chat := tgbotapi.Chat{ID: 10, Type: "group"}

	if chat.IsGroup() != true {
		t.Fail()
	}
}

func TestChatIsChannel(t *testing.T) {
	chat := tgbotapi.Chat{ID: 10, Type: "channel"}

	if chat.IsChannel() != true {
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
