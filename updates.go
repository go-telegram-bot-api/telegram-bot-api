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
				if bot.Debug == true {
					panic(err)
				} else {
					log.Println(err)
					log.Println("Fail to GetUpdates,Retry in 3 Seconds...")
					time.Sleep(time.Second * 3)
				}
			} else {
				for _, update := range updates {
					if update.UpdateID >= config.Offset {
						config.Offset = update.UpdateID + 1
						bot.Updates <- update
					}
				}
			}

		}
	}()

	return bot.Updates, nil
}
