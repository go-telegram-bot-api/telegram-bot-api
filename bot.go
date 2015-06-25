package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type Config struct {
	Token   string            `json:"token"`
	Plugins map[string]string `json:"plugins"`
}

type Plugin interface {
	GetName() string
	GetCommands() []string
	GetHelpText() []string
	GotCommand(string, Message, []string)
}

var bot *BotApi
var plugins []Plugin
var config Config

func main() {
	configPath := flag.String("config", "config.json", "path to config.json")

	flag.Parse()

	data, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Panic(err)
	}

	json.Unmarshal(data, &config)

	bot = NewBotApi(BotConfig{
		token: config.Token,
		debug: true,
	})

	plugins = []Plugin{&HelpPlugin{}, &FAPlugin{}}

	ticker := time.NewTicker(time.Second)

	lastUpdate := 0

	for range ticker.C {
		update := NewUpdate(lastUpdate + 1)
		update.Timeout = 30

		updates, err := bot.getUpdates(update)

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

				for _, cmd := range plugin.GetCommands() {
					if cmd == parts[0] {
						args := append(parts[:0], parts[1:]...)

						plugin.GotCommand(parts[0], update.Message, args)
					}
				}
			}
		}
	}
}
