package listeners

import (
	"github.com/halushko/kino-cat-core-go/nats_helper"
	"kino-cat-text-go/queue_processor"
	"log"
	"regexp"
	"strings"
)

func StartUserMessageListener() {
	processor := func(data []byte) {
		log.Printf("[StartUserMessageListener] Отримано повідомлення з NATS: %v", string(data))
		userId, messageText, err := nats_helper.ParseNatsBotText(data)
		if err != nil {
			log.Printf("[StartUserMessageListener] Помилка при парсингу повідомлення: %v", err)
		}

		log.Printf("[StartUserMessageListener] Парсинг повідомлення: chatID = %d, message = %s", userId, messageText)

		if userId != 0 && messageText != "" {
			queue, arguments := findDataToAnotherProcessorRedirection(messageText)

			nats_helper.PublishCommandMessage(queue, userId, arguments)
			log.Printf("[StartUserMessageListener] Команда \"%s\" відправлена на обробку", queue)
		} else {
			log.Printf("[StartUserMessageListener] Помилка: ID користувача чи текст повідомлення порожні")
		}
	}

	listener := &nats_helper.NatsListenerHandler{
		Function: processor,
	}

	err := nats_helper.StartNatsListener("TELEGRAM_INPUT_TEXT_QUEUE", listener)
	if err != nil {
		log.Printf("[StartUserMessageListener] Помилка: %v", err)
	}
}

func findDataToAnotherProcessorRedirection(message string) (string, []string) {
	messageLength := len(message)
	for i := 0; i < messageLength+1; i++ {
		command := message[:messageLength-i]
		queue, flag := queue_processor.FindQueueByMessage(command)
		if flag {
			log.Printf("[StartUserMessageListener] Queue \"%s\" found for \"%s\"", queue, message)

			var args []string

			switch {
			case queue == "PROCESS_HTTP_QUEUE":
				re := regexp.MustCompile(`https?://[^\s"']+`)
				args = re.FindAllString(message, -1)
			default:
				args = prepareArguments(message, command)
			}
			log.Printf("[StartUserMessageListener] Arguments for request \"%s\": \"%s\"", command, args)
			return queue, args
		}
	}
	return "", nil
}

func prepareArguments(message string, command string) []string {
	args := strings.ReplaceAll(message[len(command):], "_", " ")
	return strings.Fields(args)
}
