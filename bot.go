package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type Config struct {
	Token          string            `json:"token"`
	Plugins        map[string]string `json:"plugins"`
	EnabledPlugins map[string]bool   `json:"enabled"`
}

type Plugin interface {
	GetName() string
	GetCommands() []string
	GetHelpText() []string
	GotCommand(string, Message, []string)
	Setup()
}

var bot *BotApi
var plugins []Plugin
var config Config
var configPath *string

func main() {
	configPath = flag.String("config", "config.json", "path to config.json")

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

	plugins = []Plugin{&HelpPlugin{}, &FAPlugin{}, &ManagePlugin{}}

	for _, plugin := range plugins {
		val, ok := config.EnabledPlugins[plugin.GetName()]

		if !ok {
			fmt.Printf("Enable '%s'? [y/N] ", plugin.GetName())

			var enabled string
			fmt.Scanln(&enabled)

			if strings.ToLower(enabled) == "y" {
				plugin.Setup()
				log.Printf("Plugin '%s' started!\n", plugin.GetName())

				config.EnabledPlugins[plugin.GetName()] = true
			} else {
				config.EnabledPlugins[plugin.GetName()] = false
			}
		}

		if val {
			plugin.Setup()
			log.Printf("Plugin '%s' started!\n", plugin.GetName())
		}

		saveConfig()
	}

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
				val, _ := config.EnabledPlugins[plugin.GetName()]
				if !val {
					continue
				}

				parts := strings.Split(update.Message.Text, " ")
				command := parts[0]

				for _, cmd := range plugin.GetCommands() {
					if cmd == command {
						if bot.config.debug {
							log.Printf("'%s' matched plugin '%s'", update.Message.Text, plugin.GetName())
						}

						args := append(parts[:0], parts[1:]...)

						plugin.GotCommand(command, update.Message, args)
					}
				}
			}
		}
	}
}

func saveConfig() {
	data, _ := json.MarshalIndent(config, "", "  ")

	ioutil.WriteFile(*configPath, data, 0600)
}
