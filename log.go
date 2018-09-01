package tgbotapi

import (
	"os"
	"errors"
	stdlog "log"
)

var log = stdlog.New(os.Stderr, "", stdlog.LstdFlags)

func SetLogger(newLog *stdlog.Logger) error {
	if newLog == nil {
		return errors.New("logger is nil")
	}
	log = newLog
	return nil
}
