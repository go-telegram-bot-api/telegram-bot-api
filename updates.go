package tgbotapi

import (
	"log"
	"time"
)

// UpdatesChan starts a channel for getting updates.
func (bot *BotAPI) UpdatesChan(config UpdateConfig) error {
	bot.Updates = make(chan Update, 100)

	go func() {
		for {
			updates, err := bot.GetUpdates(config)
			if err != nil {
				log.Println(err)
				log.Println("Failed to get updates, retrying in 3 seconds...")
				time.Sleep(time.Second * 3)

				continue
			}

			for _, update := range updates {
				if update.UpdateID >= config.Offset {
					config.Offset = update.UpdateID + 1
					bot.Updates <- update
				}
			}
		}
	}()

	return nil
}
