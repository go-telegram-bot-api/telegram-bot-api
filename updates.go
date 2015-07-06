package tgbotapi

import (
	"log"
	"time"
)

// UpdatesChan returns a chan that is called whenever a new message is gotten.
func (bot *BotAPI) UpdatesChan(config UpdateConfig) (chan Update, error) {
	bot.Updates = make(chan Update, 100)

	go func() {
		for {
			updates, err := bot.GetUpdates(config)
			if err != nil {
				if bot.Debug {
					panic(err)
				} else {
					log.Println(err)
					log.Println("Failed to get updates, retrying in 3 seconds...")
					time.Sleep(time.Second * 3)
				}

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

	return bot.Updates, nil
}
