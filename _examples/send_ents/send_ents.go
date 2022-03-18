package main

import (
	"flag"
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	chatID   int64  = 0
	botToken string = ""
)

// Sends a message with message entities to the specified Telegram channel.
// Uses HtmlToEntities() to convert an HTML string to a markup free string and an array of message entities.
//
//   go run send_ents.go -id=<chat_id> -tok=<bot_token>
//
func main() {
	// Parse command line arguments.
	flag.Int64Var(&chatID, "id", 0, "Telegram channel chat ID (e.g. -1001234567890).")
	flag.StringVar(&botToken, "tok", "", "Telegram bot token.")
	flag.Parse()
	if chatID == 0 || botToken == "" {
		fmt.Printf("Chat ID and Bot Token need to be set.\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("Send Message with Entities\n")
	fmt.Printf("--------------------------\n")
	fmt.Printf("Chat ID:   %d\n", chatID)

	// Create bot.
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		fmt.Printf("Error creating bot: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Bot:       %s\n", bot.Self.UserName)
	bot.Debug = true

	// Convert HTML string to markup free text and an array of message entities.
	msgHTML := "<b>Hello</b> <i>World!</i>"
	msgStr, msgEnts, err := tgbotapi.HtmlToEntities(msgHTML, false)
	if err != nil {
		fmt.Printf("Error converting message: %s\n", err)
		os.Exit(1)
	}

	// Send message with message entities.
	msgCfg := tgbotapi.NewMessage(chatID, msgStr)
	msgCfg.ParseMode = ""
	msgCfg.Entities = msgEnts

	_, err = bot.Send(msgCfg)
	if err != nil {
		fmt.Printf("Error sending message: %s\n", err)
		os.Exit(1)
	}
}
