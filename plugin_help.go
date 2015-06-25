package main

import (
	"bytes"
	"log"
)

type HelpPlugin struct {
}

func (plugin *HelpPlugin) GetName() string {
	return "Plugins help"
}

func (plugin *HelpPlugin) GetCommands() []string {
	return []string{"/help"}
}

func (plugin *HelpPlugin) GetHelpText() []string {
	return []string{"/help (/command) - returns help about a command"}
}

func (plugin *HelpPlugin) GotCommand(command string, message Message, args []string) {
	msg := NewMessage(message.Chat.Id, "")
	msg.ReplyToMessageId = message.MessageId
	msg.DisableWebPagePreview = true

	var buffer bytes.Buffer

	if len(args) > 0 {
		for _, plug := range plugins {
			for _, cmd := range plug.GetCommands() {
				log.Println(cmd)
				log.Println(args[0])
				log.Println(args[0][1:])
				if cmd == args[0] || cmd[1:] == args[0] {
					buffer.WriteString(plug.GetName())
					buffer.WriteString("\n")

					for _, help := range plug.GetHelpText() {
						buffer.WriteString("  ")
						buffer.WriteString(help)
						buffer.WriteString("\n")
					}
				}
			}
		}
	} else {
		buffer.WriteString(config.Plugins["about_text"])
		buffer.WriteString("\n\n")

		for _, plug := range plugins {
			buffer.WriteString(plug.GetName())
			buffer.WriteString("\n")

			for _, cmd := range plug.GetHelpText() {
				buffer.WriteString("  ")
				buffer.WriteString(cmd)
				buffer.WriteString("\n")
			}

			buffer.WriteString("\n")
		}
	}

	msg.Text = buffer.String()
	bot.sendMessage(msg)
}
