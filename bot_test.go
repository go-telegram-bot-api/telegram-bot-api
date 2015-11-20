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
		log.Println(err.Error())
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

func TestSendMessage(t *testing.T) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))

	if err != nil {
		t.Fail()
	}

	msg := tgbotapi.NewMessage(36529758, "A test message from the test library in telegram-bot-api")
	bot.SendMessage(msg)
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

		bot.SendMessage(msg)
	}
}
