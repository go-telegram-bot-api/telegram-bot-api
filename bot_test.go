package tgbotapi

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	TestToken               = "153667468:AAHlSHlMqSt1f_uFmVRJbm5gntu2HI4WW8I"
	ChatID                  = 76918703
	Channel                 = "@tgbotapitest"
	SupergroupChatID        = -1001120141283
	ReplyToMessageID        = 35
	ExistingPhotoFileID     = "AgADAgADw6cxG4zHKAkr42N7RwEN3IFShCoABHQwXEtVks4EH2wBAAEC"
	ExistingDocumentFileID  = "BQADAgADOQADjMcoCcioX1GrDvp3Ag"
	ExistingAudioFileID     = "BQADAgADRgADjMcoCdXg3lSIN49lAg"
	ExistingVoiceFileID     = "AwADAgADWQADjMcoCeul6r_q52IyAg"
	ExistingVideoFileID     = "BAADAgADZgADjMcoCav432kYe0FRAg"
	ExistingVideoNoteFileID = "DQADAgADdQAD70cQSUK41dLsRMqfAg"
	ExistingStickerFileID   = "BQADAgADcwADjMcoCbdl-6eB--YPAg"
)

func getBot(t *testing.T) (*BotAPI, error) {
	bot, err := NewBotAPI(TestToken)
	bot.Debug = true

	if err != nil {
		t.Error(err)
	}

	return bot, err
}

func TestNewBotAPI_notoken(t *testing.T) {
	_, err := NewBotAPI("")

	if err == nil {
		t.Error(err)
	}
}

func TestGetUpdates(t *testing.T) {
	bot, _ := getBot(t)

	u := NewUpdate(0)

	_, err := bot.GetUpdates(u)

	if err != nil {
		t.Error(err)
	}
}

func TestSendWithMessage(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewMessage(ChatID, "A test message from the test library in telegram-bot-api")
	msg.ParseMode = "markdown"
	_, err := bot.Send(msg)

	if err != nil {
		t.Error(err)
	}
}

func TestSendWithMessageReply(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewMessage(ChatID, "A test message from the test library in telegram-bot-api")
	msg.ReplyToMessageID = ReplyToMessageID
	_, err := bot.Send(msg)

	if err != nil {
		t.Error(err)
	}
}

func TestSendWithMessageForward(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewForward(ChatID, ChatID, ReplyToMessageID)
	_, err := bot.Send(msg)

	if err != nil {
		t.Error(err)
	}
}

func TestSendWithNewPhoto(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewPhotoUpload(ChatID, "tests/image.jpg")
	msg.Caption = "Test"
	_, err := bot.Send(msg)

	if err != nil {
		t.Error(err)
	}
}

func TestSendWithNewPhotoWithFileBytes(t *testing.T) {
	bot, _ := getBot(t)

	data, _ := ioutil.ReadFile("tests/image.jpg")
	b := FileBytes{Name: "image.jpg", Bytes: data}

	msg := NewPhotoUpload(ChatID, b)
	msg.Caption = "Test"
	_, err := bot.Send(msg)

	if err != nil {
		t.Error(err)
	}
}

func TestSendWithNewPhotoWithFileReader(t *testing.T) {
	bot, _ := getBot(t)

	f, _ := os.Open("tests/image.jpg")
	reader := FileReader{Name: "image.jpg", Reader: f, Size: -1}

	msg := NewPhotoUpload(ChatID, reader)
	msg.Caption = "Test"
	_, err := bot.Send(msg)

	if err != nil {
		t.Error(err)
	}
}

func TestSendWithNewPhotoReply(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewPhotoUpload(ChatID, "tests/image.jpg")
	msg.ReplyToMessageID = ReplyToMessageID

	_, err := bot.Send(msg)

	if err != nil {
		t.Error(err)
	}
}

func TestSendNewPhotoToChannel(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewPhotoUploadToChannel(Channel, "tests/image.jpg")
	msg.Caption = "Test"
	_, err := bot.Send(msg)

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestSendNewPhotoToChannelFileBytes(t *testing.T) {
	bot, _ := getBot(t)

	data, _ := ioutil.ReadFile("tests/image.jpg")
	b := FileBytes{Name: "image.jpg", Bytes: data}

	msg := NewPhotoUploadToChannel(Channel, b)
	msg.Caption = "Test"
	_, err := bot.Send(msg)

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestSendNewPhotoToChannelFileReader(t *testing.T) {
	bot, _ := getBot(t)

	f, _ := os.Open("tests/image.jpg")
	reader := FileReader{Name: "image.jpg", Reader: f, Size: -1}

	msg := NewPhotoUploadToChannel(Channel, reader)
	msg.Caption = "Test"
	_, err := bot.Send(msg)

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestSendWithExistingPhoto(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewPhotoShare(ChatID, ExistingPhotoFileID)
	msg.Caption = "Test"
	_, err := bot.Send(msg)

	if err != nil {
		t.Error(err)
	}
}

func TestSendWithNewDocument(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewDocumentUpload(ChatID, "tests/image.jpg")
	_, err := bot.Send(msg)

	if err != nil {
		t.Error(err)
	}
}

func TestSendWithExistingDocument(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewDocumentShare(ChatID, ExistingDocumentFileID)
	_, err := bot.Send(msg)

	if err != nil {
		t.Error(err)
	}
}

func TestSendWithNewAudio(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewAudioUpload(ChatID, "tests/audio.mp3")
	msg.Title = "TEST"
	msg.Duration = 10
	msg.Performer = "TEST"
	msg.MimeType = "audio/mpeg"
	msg.FileSize = 688
	_, err := bot.Send(msg)

	if err != nil {
		t.Error(err)
	}
}

func TestSendWithExistingAudio(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewAudioShare(ChatID, ExistingAudioFileID)
	msg.Title = "TEST"
	msg.Duration = 10
	msg.Performer = "TEST"

	_, err := bot.Send(msg)

	if err != nil {
		t.Error(err)
	}
}

func TestSendWithNewVoice(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewVoiceUpload(ChatID, "tests/voice.ogg")
	msg.Duration = 10
	_, err := bot.Send(msg)

	if err != nil {
		t.Error(err)
	}
}

func TestSendWithExistingVoice(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewVoiceShare(ChatID, ExistingVoiceFileID)
	msg.Duration = 10
	_, err := bot.Send(msg)

	if err != nil {
		t.Error(err)
	}
}

func TestSendWithContact(t *testing.T) {
	bot, _ := getBot(t)

	contact := NewContact(ChatID, "5551234567", "Test")

	if _, err := bot.Send(contact); err != nil {
		t.Error(err)
	}
}

func TestSendWithLocation(t *testing.T) {
	bot, _ := getBot(t)

	_, err := bot.Send(NewLocation(ChatID, 40, 40))

	if err != nil {
		t.Error(err)
	}
}

func TestSendWithVenue(t *testing.T) {
	bot, _ := getBot(t)

	venue := NewVenue(ChatID, "A Test Location", "123 Test Street", 40, 40)

	if _, err := bot.Send(venue); err != nil {
		t.Error(err)
	}
}

func TestSendWithNewVideo(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewVideoUpload(ChatID, "tests/video.mp4")
	msg.Duration = 10
	msg.Caption = "TEST"

	_, err := bot.Send(msg)

	if err != nil {
		t.Error(err)
	}
}

func TestSendWithExistingVideo(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewVideoShare(ChatID, ExistingVideoFileID)
	msg.Duration = 10
	msg.Caption = "TEST"

	_, err := bot.Send(msg)

	if err != nil {
		t.Error(err)
	}
}

func TestSendWithNewVideoNote(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewVideoNoteUpload(ChatID, 240, "tests/videonote.mp4")
	msg.Duration = 10

	_, err := bot.Send(msg)

	if err != nil {
		t.Error(err)
	}
}

func TestSendWithExistingVideoNote(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewVideoNoteShare(ChatID, 240, ExistingVideoNoteFileID)
	msg.Duration = 10

	_, err := bot.Send(msg)

	if err != nil {
		t.Error(err)
	}
}

func TestSendWithNewSticker(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewStickerUpload(ChatID, "tests/image.jpg")

	_, err := bot.Send(msg)

	if err != nil {
		t.Error(err)
	}
}

func TestSendWithExistingSticker(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewStickerShare(ChatID, ExistingStickerFileID)

	_, err := bot.Send(msg)

	if err != nil {
		t.Error(err)
	}
}

func TestSendWithNewStickerAndKeyboardHide(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewStickerUpload(ChatID, "tests/image.jpg")
	msg.ReplyMarkup = ReplyKeyboardRemove{
		RemoveKeyboard: true,
		Selective:      false,
	}
	_, err := bot.Send(msg)

	if err != nil {
		t.Error(err)
	}
}

func TestSendWithExistingStickerAndKeyboardHide(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewStickerShare(ChatID, ExistingStickerFileID)
	msg.ReplyMarkup = ReplyKeyboardRemove{
		RemoveKeyboard: true,
		Selective:      false,
	}

	_, err := bot.Send(msg)

	if err != nil {
		t.Error(err)
	}
}

func TestGetFile(t *testing.T) {
	bot, _ := getBot(t)

	file := FileConfig{
		FileID: ExistingPhotoFileID,
	}

	_, err := bot.GetFile(file)

	if err != nil {
		t.Error(err)
	}
}

func TestSendChatConfig(t *testing.T) {
	bot, _ := getBot(t)

	_, err := bot.Request(NewChatAction(ChatID, ChatTyping))

	if err != nil {
		t.Error(err)
	}
}

func TestSendEditMessage(t *testing.T) {
	bot, _ := getBot(t)

	msg, err := bot.Send(NewMessage(ChatID, "Testing editing."))
	if err != nil {
		t.Error(err)
	}

	edit := EditMessageTextConfig{
		BaseEdit: BaseEdit{
			ChatID:    ChatID,
			MessageID: msg.MessageID,
		},
		Text: "Updated text.",
	}

	_, err = bot.Send(edit)
	if err != nil {
		t.Error(err)
	}
}

func TestGetUserProfilePhotos(t *testing.T) {
	bot, _ := getBot(t)

	_, err := bot.GetUserProfilePhotos(NewUserProfilePhotos(ChatID))
	if err != nil {
		t.Error(err)
	}
}

func TestSetWebhookWithCert(t *testing.T) {
	bot, _ := getBot(t)

	time.Sleep(time.Second * 2)

	bot.Request(DeleteWebhookConfig{})

	wh := NewWebhookWithCert("https://example.com/tgbotapi-test/"+bot.Token, "tests/cert.pem")
	_, err := bot.Request(wh)
	if err != nil {
		t.Error(err)
	}

	_, err = bot.GetWebhookInfo()

	if err != nil {
		t.Error(err)
	}

	bot.Request(DeleteWebhookConfig{})
}

func TestSetWebhookWithoutCert(t *testing.T) {
	bot, _ := getBot(t)

	time.Sleep(time.Second * 2)

	bot.Request(DeleteWebhookConfig{})

	wh := NewWebhook("https://example.com/tgbotapi-test/" + bot.Token)
	_, err := bot.Request(wh)
	if err != nil {
		t.Error(err)
	}

	info, err := bot.GetWebhookInfo()

	if err != nil {
		t.Error(err)
	}
	if info.MaxConnections == 0 {
		t.Errorf("Expected maximum connections to be greater than 0")
	}
	if info.LastErrorDate != 0 {
		t.Errorf("failed to set webhook: %s", info.LastErrorMessage)
	}

	bot.Request(DeleteWebhookConfig{})
}

func TestSendWithMediaGroup(t *testing.T) {
	bot, _ := getBot(t)

	cfg := NewMediaGroup(ChatID, []interface{}{
		NewInputMediaPhoto("https://github.com/go-telegram-bot-api/telegram-bot-api/raw/0a3a1c8716c4cd8d26a262af9f12dcbab7f3f28c/tests/image.jpg"),
		NewInputMediaVideo("https://github.com/go-telegram-bot-api/telegram-bot-api/raw/0a3a1c8716c4cd8d26a262af9f12dcbab7f3f28c/tests/video.mp4"),
	})

	messages, err := bot.SendMediaGroup(cfg)
	if err != nil {
		t.Error(err)
	}

	if messages == nil {
		t.Error()
	}

	if len(messages) != 2 {
		t.Error()
	}
}

func ExampleNewBotAPI() {
	bot, err := NewBotAPI("MyAwesomeBotToken")
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Optional: wait for updates and clear them if you don't want to handle
	// a large backlog of old messages
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}

func ExampleNewWebhook() {
	bot, err := NewBotAPI("MyAwesomeBotToken")
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.Request(NewWebhookWithCert("https://www.google.com:8443/"+bot.Token, "cert.pem"))
	if err != nil {
		panic(err)
	}

	info, err := bot.GetWebhookInfo()

	if err != nil {
		panic(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("failed to set webhook: %s", info.LastErrorMessage)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServeTLS("0.0.0.0:8443", "cert.pem", "key.pem", nil)

	for update := range updates {
		log.Printf("%+v\n", update)
	}
}

func ExampleWebhookHandler() {
	bot, err := NewBotAPI("MyAwesomeBotToken")
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.Request(NewWebhookWithCert("https://www.google.com:8443/"+bot.Token, "cert.pem"))
	if err != nil {
		panic(err)
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		panic(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("[Telegram callback failed]%s", info.LastErrorMessage)
	}

	http.HandleFunc("/"+bot.Token, func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%+v\n", bot.HandleUpdate(w, r))
	})

	go http.ListenAndServeTLS("0.0.0.0:8443", "cert.pem", "key.pem", nil)
}

func ExampleInlineConfig() {
	bot, err := NewBotAPI("MyAwesomeBotToken") // create new bot
	if err != nil {
		panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.InlineQuery == nil { // if no inline query, ignore it
			continue
		}

		article := NewInlineQueryResultArticle(update.InlineQuery.ID, "Echo", update.InlineQuery.Query)
		article.Description = update.InlineQuery.Query

		inlineConf := InlineConfig{
			InlineQueryID: update.InlineQuery.ID,
			IsPersonal:    true,
			CacheTime:     0,
			Results:       []interface{}{article},
		}

		if _, err := bot.Request(inlineConf); err != nil {
			log.Println(err)
		}
	}
}

func TestDeleteMessage(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewMessage(ChatID, "A test message from the test library in telegram-bot-api")
	msg.ParseMode = "markdown"
	message, _ := bot.Send(msg)

	deleteMessageConfig := DeleteMessageConfig{
		ChatID:    message.Chat.ID,
		MessageID: message.MessageID,
	}
	_, err := bot.Request(deleteMessageConfig)

	if err != nil {
		t.Error(err)
	}
}

func TestPinChatMessage(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewMessage(SupergroupChatID, "A test message from the test library in telegram-bot-api")
	msg.ParseMode = "markdown"
	message, _ := bot.Send(msg)

	pinChatMessageConfig := PinChatMessageConfig{
		ChatID:              message.Chat.ID,
		MessageID:           message.MessageID,
		DisableNotification: false,
	}
	_, err := bot.Request(pinChatMessageConfig)

	if err != nil {
		t.Error(err)
	}
}

func TestUnpinChatMessage(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewMessage(SupergroupChatID, "A test message from the test library in telegram-bot-api")
	msg.ParseMode = "markdown"
	message, _ := bot.Send(msg)

	// We need pin message to unpin something
	pinChatMessageConfig := PinChatMessageConfig{
		ChatID:              message.Chat.ID,
		MessageID:           message.MessageID,
		DisableNotification: false,
	}

	if _, err := bot.Request(pinChatMessageConfig); err != nil {
		t.Error(err)
	}

	unpinChatMessageConfig := UnpinChatMessageConfig{
		ChatID: message.Chat.ID,
	}

	if _, err := bot.Request(unpinChatMessageConfig); err != nil {
		t.Error(err)
	}
}

func TestPolls(t *testing.T) {
	bot, _ := getBot(t)

	poll := NewPoll(SupergroupChatID, "Are polls working?", "Yes", "No")

	msg, err := bot.Send(poll)
	if err != nil {
		t.Error(err)
	}

	result, err := bot.StopPoll(NewStopPoll(SupergroupChatID, msg.MessageID))
	if err != nil {
		t.Error(err)
	}

	if result.Question != "Are polls working?" {
		t.Error("Poll question did not match")
	}

	if !result.IsClosed {
		t.Error("Poll did not end")
	}

	if result.Options[0].Text != "Yes" || result.Options[0].VoterCount != 0 || result.Options[1].Text != "No" || result.Options[1].VoterCount != 0 {
		t.Error("Poll options were incorrect")
	}
}

func TestSendDice(t *testing.T) {
	bot, _ := getBot(t)

	dice := NewSendDice(ChatID)

	msg, err := bot.Send(dice)
	if err != nil {
		t.Error("Unable to send dice roll")
	}

	if msg.Dice == nil {
		t.Error("Dice roll was not received")
	}
}

func TestSetCommands(t *testing.T) {
	bot, _ := getBot(t)

	setCommands := NewSetMyCommands(BotCommand{
		Command:     "test",
		Description: "a test command",
	})

	if _, err := bot.Request(setCommands); err != nil {
		t.Error("Unable to set commands")
	}

	commands, err := bot.GetMyCommands()
	if err != nil {
		t.Error("Unable to get commands")
	}

	if len(commands) != 1 {
		t.Error("Incorrect number of commands returned")
	}

	if commands[0].Command != "test" || commands[0].Description != "a test command" {
		t.Error("Commands were incorrectly set")
	}
}
