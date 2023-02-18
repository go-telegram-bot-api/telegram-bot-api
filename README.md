# Golang bindings for the Telegram Bot API

Experimental, temporary replacement for [go-telegram-bot-api](http://github.com/go-telegram-bot-api/telegram-bot-api).

## Installation

1. Add this to your `go.mod` file:
```
replace github.com/go-telegram-bot-api/telegram-bot-api/v5 => github.com/iamwavecut/telegram-bot-api latest
```
2. Run `go mod tidy`
3. Import and use the package as usual:
```
import (
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)
```