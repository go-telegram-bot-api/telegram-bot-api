package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/ddliu/go-httpclient"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type FAPlugin struct {
}

func (plugin *FAPlugin) GetName() string {
	return "FA Mirrorer"
}

func (plugin *FAPlugin) GetCommands() []string {
	return []string{"/fa"}
}

func (plugin *FAPlugin) GetHelpText() []string {
	return []string{"/fa [link] - mirrors an image from FurAffinity"}
}

func (plugin *FAPlugin) Setup() {
	a, ok := config.Plugins["fa_a"]
	if !ok {
		fmt.Print("FurAffinity Cookie a: ")
		fmt.Scanln(&a)

		config.Plugins["fa_a"] = a
	}

	b, ok := config.Plugins["fa_b"]
	if !ok {
		fmt.Print("FurAffinity Cookie b: ")
		fmt.Scanln(&b)

		config.Plugins["fa_b"] = b
	}
}

func (plugin *FAPlugin) GotCommand(command string, message Message, args []string) {
	if len(args) == 0 {
		bot.sendMessage(NewMessage(message.Chat.Id, "You need to include a link!"))

		return
	}

	bot.sendChatAction(NewChatAction(message.Chat.Id, CHAT_UPLOAD_PHOTO))

	_, err := strconv.Atoi(args[0])
	if err == nil {
		args[0] = "http://www.furaffinity.net/view/" + args[0]
	}

	resp, err := httpclient.WithCookie(&http.Cookie{
		Name:  "b",
		Value: config.Plugins["fa_b"],
	}).WithCookie(&http.Cookie{
		Name:  "a",
		Value: config.Plugins["fa_a"],
	}).Get(args[0], nil)
	if err != nil {
		bot.sendMessage(NewMessage(message.Chat.Id, "ERR : "+err.Error()))
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		bot.sendMessage(NewMessage(message.Chat.Id, "ERR : "+err.Error()))
	}

	sel := doc.Find("#submissionImg")
	for i := range sel.Nodes {
		single := sel.Eq(i)

		val, _ := single.Attr("src")

		tokens := strings.Split(val, "/")
		fileName := tokens[len(tokens)-1]

		output, _ := os.Create(fileName)
		defer output.Close()
		defer os.Remove(output.Name())

		resp, _ := http.Get("http:" + val)
		defer resp.Body.Close()

		io.Copy(output, resp.Body)

		bot.sendPhoto(NewPhotoUpload(message.Chat.Id, output.Name()))
	}
}
