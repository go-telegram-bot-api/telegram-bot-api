package tgbotapi_test

import (
	"github.com/zhulik/telegram-bot-api"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	botToken := os.Getenv("TELEGRAM_API_TOKEN")

	if botToken == "" {
		log.Panic("You must provide a TELEGRAM_API_TOKEN env variable to test!")
	}

	os.Exit(m.Run())
}

func TestNewBotAPI_notoken(t *testing.T) {
	_, err := tgbotapi.NewBotAPI("")

	if err == nil {
		t.Fail()
	}
}

func TestNewBotAPI_token(t *testing.T) {
	_, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))

	if err != nil {
		t.Fail()
	}
}

func TestGetUpdates(t *testing.T) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))

	if err != nil {
		t.Fail()
	}

	u := tgbotapi.NewUpdate(0)

	_, err = bot.GetUpdates(u)

	if err != nil {
		t.Fail()
	}
}

func TestSendWithMessage(t *testing.T) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))

	if err != nil {
		t.Fail()
	}

	msg := tgbotapi.NewMessage(76918703, "A test message from the test library in telegram-bot-api")
	_, err = bot.Send(msg)

	if err != nil {
		t.Fail()
	}
}

func TestSendWithNewPhoto(t *testing.T) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))

	if err != nil {
		t.Fail()
	}

	msg := tgbotapi.NewPhotoUpload(76918703, "tests/image.jpg")
	_, err = bot.Send(msg)

	if err != nil {
		t.Fail()
	}
}



func TestSendWithExistingPhoto(t *testing.T) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))

	if err != nil {
		t.Fail()
	}

	msg := tgbotapi.NewPhotoShare(76918703, "AgADAgADxKcxG4cBswqt13DnHOgbmBxDhCoABC0h01_AL4SKe20BAAEC")
	_, err = bot.Send(msg)

	if err != nil {
		t.Fail()
	}
}

func TestSendWithNewDocument(t *testing.T) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))

	if err != nil {
		t.Fail()
	}

	msg := tgbotapi.NewDocumentUpload(76918703, "tests/image.jpg")
	_, err = bot.Send(msg)

	if err != nil {
		t.Fail()
	}
}

func TestSendWithExistingDocument(t *testing.T) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))

	if err != nil {
		t.Fail()
	}

	msg := tgbotapi.NewDocumentShare(76918703, "BQADAgADBwADhwGzCjWgiUU4T8VNAg")
	_, err = bot.Send(msg)

	if err != nil {
		t.Fail()
	}
}

func TestGetFile(t *testing.T) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))

	if err != nil {
		t.Fail()
	}

	file := tgbotapi.FileConfig{"BQADAgADBwADhwGzCjWgiUU4T8VNAg"}

	_, err = bot.GetFile(file)

	if err != nil {
		t.Fail()
	}
}

func TestSendChatConfig(t *testing.T) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))

	if err != nil {
		t.Fail()
	}

	err = bot.SendChatAction(tgbotapi.NewChatAction(76918703, tgbotapi.ChatTyping))

	if err != nil {
		t.Fail()
	}
}

func TestGetUserProfilePhotos(t *testing.T) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))

	if err != nil {
		t.Fail()
	}

	_, err = bot.GetUserProfilePhotos(tgbotapi.NewUserProfilePhotos(76918703))
	if err != nil {
		t.Fail()
	}
}

func TestUpdatesChan(t *testing.T) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))

	if err != nil {
		t.Fail()
	}

	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	err = bot.UpdatesChan(ucfg)
	
	if err != nil {
		t.Fail()
	}
}

func ExampleNewBotAPI() {
	bot, err := tgbotapi.NewBotAPI("MyAwesomeBotToken")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	err = bot.UpdatesChan(u)

	for update := range bot.Updates {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
