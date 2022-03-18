package tgbotapi

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf16"

	"golang.org/x/net/html"
)

// Converts an HTML string to a markup free string and an array of Telegram message entities.
// Set strict to true if you want the function to error out on bad HTML or false to ignore.
// NOTE: The offset and length fields in the message entities are the number of UTF-16 code units (uint16) and not the number of characters.
// For example the 😂 emoji is one 4 byte UTF-16 character made up of 2 16-bit code units (0xd83d and 0xde02).
// Its code unit length is 2.
// https://core.telegram.org/api/entities
// https://core.telegram.org/bots/api#html-style
// https://core.telegram.org/bots/api#messageentity
// https://pkg.go.dev/golang.org/x/net/html
// https://www.ibm.com/docs/en/db2-for-zos/11?topic=unicode-utfs
func HtmlToEntities(str string, strict bool) (text string, entities []MessageEntity, err error) {
	// Keep track of a LIFO stack of entities. When the start tag is scanned, push a new entity on the stack.
	// When end tag is scanned, pop the last entity off the stack and add to entities array.
	// Mismatched tags are either ignored or return an error depending on if strict is true.
	stack := []MessageEntity{}
	t := html.NewTokenizer(strings.NewReader(str))

	uStr := utf16.Encode([]rune(text))
loop:
	for {
		tt := t.Next()
		switch tt {
		case html.ErrorToken:
			err = t.Err()
			if err == io.EOF {
				err = nil
				break loop
			}
			return
		case html.TextToken:
			uStr = append(uStr, utf16.Encode([]rune(string(t.Text())))...)
		case html.StartTagToken:
			// push on stack
			me := getEntity(t)
			if me.Type == "" {
				// ignore tags we don't know about
				continue
			}
			me.Offset = len(uStr)
			stack = append(stack, me)
		case html.EndTagToken:
			// pop off statck
			me := getEntity(t)
			if me.Type == "" {
				// ignore tags we don't know about
				continue
			}
			if len(stack) == 0 {
				if strict {
					err = fmt.Errorf("unexpected end tag: %s", me.Tag)
					return
				}
				continue
			}
			last := stack[len(stack)-1]
			if last.Tag != me.Tag {
				if strict {
					err = fmt.Errorf("unexpected end tag: %s", me.Tag)
					return
				}
				continue
			}

			stack = stack[:len(stack)-1] // pop
			last.Length = len(uStr) - last.Offset
			if last.Length == 0 {
				// skip tags that have no content
				continue
			}
			entities = append(entities, last)
		}
	}
	// convert UTF-16 to UTF-8
	text += string(utf16.Decode(uStr))
	return
}

// Gets attribute value for the current tag. Returns empty string if not found.
func getAttr(t *html.Tokenizer, findKey string) string {
	hasMore := true
	var key, value []byte
	for hasMore {
		key, value, hasMore = t.TagAttr()
		if string(key) == findKey {
			return string(value)
		}
	}
	return ""
}

// Creates a message entity for the current token.
// Figures out the entity type equivalent to the token's HTML tag (e.g "b" -> "bold").
// Sets the type to empty string if no mapping found.
//
// https://core.telegram.org/api/entities
// https://core.telegram.org/bots/api#formatting-options
//
//   Entity Type    Tags
//   -----------    ----
//   bold           <b>, <strong>
//   code           <code>
//   italic         <em>, <i>
//   pre            <pre language="{language}">
//   spoiler        <span class="tg-spoiler">, <tg-spoiler>
//   strikethrough  <del>, <s>, <strike>
//   text_link      <a href="https://...">
//   text_mention   <a href="tg://user?id={user}">
//   underline      <ins>, <u>
//
func getEntity(t *html.Tokenizer) (me MessageEntity) {
	name, hasAttr := t.TagName()
	me.Tag = string(name)
	switch strings.ToLower(string(name)) {
	case "a":
		me.Type = "text_link"
		if hasAttr {
			href := getAttr(t, "href")
			if strings.HasPrefix(href, "tg://user?id=") {
				me.Type = "text_mention"
				me.User.ID, _ = strconv.ParseInt(href[len("tg://user?id="):], 10, 64)
			} else {
				me.URL = href
			}
		}
	case "b":
		me.Type = "bold"
	case "code":
		me.Type = "code"
	case "del":
		me.Type = "strikethrough"
	case "em":
		me.Type = "italic"
	case "i":
		me.Type = "italic"
	case "ins":
		me.Type = "underline"
	case "pre":
		me.Type = "pre"
		if hasAttr {
			me.Language = getAttr(t, "language")
		}
	case "s":
		me.Type = "strikethrough"
	case "span":
		if hasAttr {
			if getAttr(t, "class") == "tg-spoiler" {
				me.Type = "spoiler"
			}
		}
	case "strike":
		me.Type = "strikethrough"
	case "strong":
		me.Type = "bold"
	case "tg-spoiler":
		me.Type = "spoiler"
	case "u":
		me.Type = "underline"
	}
	return
}
