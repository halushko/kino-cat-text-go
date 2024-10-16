package main

import (
	"github.com/halushko/kino-cat-core-go/logger_helper"
	"kino-cat-text-go/listeners"
)

//goland:noinspection ALL
func main() {
	logFile := logger_helper.SoftPrepareLogFile()

	go listeners.StartUserMessageListener()
	go listeners.StartGetHelpCommandListener()

	defer logger_helper.SoftLogClose(logFile)

	select {}
}
