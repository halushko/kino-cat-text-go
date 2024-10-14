package kino_cat_text_go

import (
	"github.com/halushko/kino-cat-core-go/logger_helper"
	"kino-cat-text-go/listeners"
	"log"
)

//goland:noinspection ALL
func main() {
	logFile := logger_helper.SoftPrepareLogFile()
	log.SetOutput(logFile)

	listeners.StartUserMessageListener()
	listeners.StartGetHelpCommandListener()

	logger_helper.SoftLogClose(logFile)
}
