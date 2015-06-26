package tgbotapi

type BotApi struct {
	Token   string      `json:"token"`
	Debug   bool        `json:"debug"`
	Self    User        `json:"-"`
	Updates chan Update `json:"-"`
}

func NewBotApi(token string) (*BotApi, error) {
	bot := &BotApi{
		Token: token,
	}

	self, err := bot.GetMe()
	if err != nil {
		return &BotApi{}, err
	}

	bot.Self = self

	return bot, nil
}
