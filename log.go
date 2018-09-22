package tgbotapi

import (
	"errors"
	stdlog "log"
	"os"
)

var log = stdlog.New(os.Stderr, "", stdlog.LstdFlags)

// SetLogger specifies the logger that the package should use.
func SetLogger(newLog *stdlog.Logger) error {
	if newLog == nil {
		return errors.New("logger is nil")
	}
	log = newLog
	return nil
}
