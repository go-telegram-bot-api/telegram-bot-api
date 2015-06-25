package main

import (
	"fmt"
	"log"
	"strings"
)

type ManagePlugin struct {
}

func (plugin *ManagePlugin) GetName() string {
	return "Plugin manager"
}

func (plugin *ManagePlugin) GetCommands() []string {
	return []string{
		"/enable",
		"Enable",
		"/disable",
		"Disable",
		"/reload",
	}
}

func (plugin *ManagePlugin) GetHelpText() []string {
	return []string{
		"/enable [name] - enables a plugin",
		"/disable [name] - disables a plugin",
		"/reload - reloads bot configuration",
	}
}

func (plugin *ManagePlugin) Setup() {
}

func (plugin *ManagePlugin) GotCommand(command string, message Message, args []string) {
	log.Println(command)

	if command == "/enable" {
		keyboard := [][]string{}

		hasDisabled := false
		for _, plug := range plugins {
			enabled, _ := config.EnabledPlugins[plug.GetName()]
			if enabled {
				continue
			}

			hasDisabled = true
			keyboard = append(keyboard, []string{"Enable " + plug.GetName()})
		}

		if !hasDisabled {
			msg := NewMessage(message.Chat.Id, "All plugins are enabled!")
			msg.ReplyToMessageId = message.MessageId

			bot.sendMessage(msg)

			return
		}

		msg := NewMessage(message.Chat.Id, "Please specify which plugin to enable")
		msg.ReplyToMessageId = message.MessageId
		msg.ReplyMarkup = ReplyKeyboardMarkup{
			Keyboard:        keyboard,
			OneTimeKeyboard: true,
			Selective:       true,
			ResizeKeyboard:  true,
		}

		bot.sendMessage(msg)
	} else if command == "Enable" {
		pluginName := strings.SplitN(message.Text, " ", 2)

		msg := NewMessage(message.Chat.Id, "")
		msg.ReplyToMessageId = message.MessageId
		msg.ReplyMarkup = ReplyKeyboardHide{
			HideKeyboard: true,
			Selective:    true,
		}

		_, ok := config.EnabledPlugins[pluginName[1]]
		if !ok {
			msg.Text = "Unknown plugin!"
			msg.ReplyToMessageId = message.MessageId
			bot.sendMessage(msg)

			return
		}

		config.EnabledPlugins[pluginName[1]] = true
		msg.Text = fmt.Sprintf("Enabled '%s'!", pluginName[1])
		bot.sendMessage(msg)
	} else if command == "/disable" {
		keyboard := [][]string{}

		hasEnabled := false
		for _, plug := range plugins {
			enabled, _ := config.EnabledPlugins[plug.GetName()]
			if !enabled {
				continue
			}

			hasEnabled = true
			keyboard = append(keyboard, []string{"Disable " + plug.GetName()})
		}

		if !hasEnabled {
			msg := NewMessage(message.Chat.Id, "All plugins are disabled!")
			msg.ReplyToMessageId = message.MessageId

			bot.sendMessage(msg)

			return
		}

		msg := NewMessage(message.Chat.Id, "Please specify which plugin to disable")
		msg.ReplyToMessageId = message.MessageId
		msg.ReplyMarkup = ReplyKeyboardMarkup{
			Keyboard:        keyboard,
			OneTimeKeyboard: true,
			Selective:       true,
			ResizeKeyboard:  true,
		}

		bot.sendMessage(msg)
	} else if command == "Disable" {
		pluginName := strings.SplitN(message.Text, " ", 2)

		msg := NewMessage(message.Chat.Id, "")
		msg.ReplyToMessageId = message.MessageId
		msg.ReplyMarkup = ReplyKeyboardHide{
			HideKeyboard: true,
			Selective:    true,
		}

		_, ok := config.EnabledPlugins[pluginName[1]]
		if !ok {
			msg.Text = "Unknown plugin!"
			msg.ReplyToMessageId = message.MessageId
			bot.sendMessage(msg)

			return
		}

		config.EnabledPlugins[pluginName[1]] = false
		msg.Text = fmt.Sprintf("Disabled '%s'!", pluginName[1])
		bot.sendMessage(msg)
	}

	saveConfig()
}
