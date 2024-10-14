package main

import (
	"github.com/halushko/kino-cat-core-go/logger_helper"
	"kino-cat-text-go/listeners"
	"log"
)

//goland:noinspection ALL
func main() {
	logFile := logger_helper.SoftPrepareLogFile()
	log.SetOutput(logFile)

	go listeners.StartUserMessageListener()
	go listeners.StartGetHelpCommandListener()

	defer logger_helper.SoftLogClose(logFile)

	select {}
}
