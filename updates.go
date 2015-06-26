package tgbotapi

func (bot *BotApi) UpdatesChan(config UpdateConfig) (chan Update, error) {
	bot.Updates = make(chan Update, 100)

	go func() {
		updates, err := bot.GetUpdates(config)
		if err != nil {
			panic(err)
		}

		for _, update := range updates {
			if update.UpdateId > config.Offset {
				config.Offset = update.UpdateId + 1
			}

			bot.Updates <- update
		}
	}()

	return bot.Updates, nil
}
