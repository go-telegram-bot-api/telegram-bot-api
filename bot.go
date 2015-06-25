package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Token string `json:"token"`
}

type Plugin interface {
	GetCommand() string
	GotCommand(Message, []string)
}

type ColonThree struct {
}

func (plugin *ColonThree) GetCommand() string {
	return "/three"
}

func (plugin *ColonThree) GotCommand(message Message, args []string) {
	if len(args) > 0 {
		n, err := strconv.Atoi(args[0])
		if err != nil {
			msg := NewMessage(message.Chat.Id, "Bad number!")
			msg.ReplyToMessageId = message.MessageId

			bot.sendMessage(msg)

			return
		}

		if n > 5 {
			msg := NewMessage(message.Chat.Id, "That's a bit much, no?")
			msg.ReplyToMessageId = message.MessageId

			bot.sendMessage(msg)

			return
		}

		for i := 0; i < n; i++ {
			bot.sendMessage(NewMessage(message.Chat.Id, ":3"))
		}
	} else {
		bot.sendMessage(NewMessage(message.Chat.Id, ":3"))

		bot.sendPhoto(NewPhotoUpload(message.Chat.Id, "fox.png"))
	}
}

var bot *BotApi

func main() {
	configPath := flag.String("config", "config.json", "path to config.json")

	flag.Parse()

	data, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Panic(err)
	}

	var cfg Config
	json.Unmarshal(data, &cfg)

	bot = NewBotApi(BotConfig{
		token: cfg.Token,
		debug: true,
	})

	plugins := []Plugin{&ColonThree{}}

	ticker := time.NewTicker(5 * time.Second)

	lastUpdate := 0

	for range ticker.C {
		updates, err := bot.getUpdates(NewUpdate(lastUpdate + 1))

		if err != nil {
			log.Panic(err)
		}

		for _, update := range updates {
			lastUpdate = update.UpdateId

			if update.Message.Text == "" {
				continue
			}

			for _, plugin := range plugins {
				parts := strings.Split(update.Message.Text, " ")

				if plugin.GetCommand() == parts[0] {
					parts = append(parts[:0], parts[1:]...)

					plugin.GotCommand(update.Message, parts)
				}
			}
		}
	}
}
