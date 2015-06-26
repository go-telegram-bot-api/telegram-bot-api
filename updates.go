package tgbotapi

// UpdatesChan returns a chan that is called whenever a new message is gotten.
func (bot *BotAPI) UpdatesChan(config UpdateConfig) (chan Update, error) {
	bot.Updates = make(chan Update, 100)

	go func() {
		updates, err := bot.GetUpdates(config)
		if err != nil {
			panic(err)
		}

		for _, update := range updates {
			if update.UpdateID > config.Offset {
				config.Offset = update.UpdateID + 1
			}

			bot.Updates <- update
		}
	}()

	return bot.Updates, nil
}
