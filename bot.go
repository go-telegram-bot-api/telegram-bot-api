package tgbotapi

type BotApi struct {
	Token   string      `json:"token"`
	Debug   bool        `json:"debug"`
	Updates chan Update `json:"-"`
}

func NewBotApi(token string) *BotApi {
	return &BotApi{
		Token: token,
	}
}
