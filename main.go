package kino_cat_text_go

import (
	"kino-cat-text-go/listeners"
	"log"
	"os"
)

//goland:noinspection ALL
func main() {
	logFile := prepareLogFile()
	log.SetOutput(logFile)
	defer func() {
		if err := logFile.Close(); err != nil {
			log.Println("Помилка при спробі закрити лог файл")
		}
	}()

	listeners.StartUserMessageListener()
	listeners.StartGetHelpCommandListener()
}

func prepareLogFile() *os.File {
	log.Print("Старт бота")

	logFile, err := os.OpenFile("bot.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Fatal(err)
	}

	return logFile
}
