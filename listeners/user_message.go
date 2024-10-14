package listeners

import (
	"encoding/json"
	"github.com/halushko/kino-cat-core-go/nats_helper"
	"github.com/nats-io/nats.go"
	"kino-cat-text-go/queue_processor"
	"log"
	"strings"
)

type TelegramUserNatsMessage struct {
	ChatId int64  `json:"chat_id"`
	Text   string `json:"text"`
}

type commandNatsMessage struct {
	ChatID    int64    `json:"chat_id"`
	Arguments []string `json:"arguments"`
}

func StartUserMessageListener() {
	processor := func(msg *nats.Msg) {
		log.Printf("[StartNatsListener] Отримано повідомлення з NATS: %s", string(msg.Data))
		chatId, messageText := parseNatsMessage(msg.Data)

		log.Printf("[StartNatsListener] Парсинг повідомлення: chatID = %d, message = %s", chatId, messageText) // Новый лог для проверки данных

		if chatId != 0 && messageText != "" {
			queue, arguments := findDataToAnotherProcessorRedirection(messageText)

			if request, errMarshal := json.Marshal(commandNatsMessage{
				ChatID:    chatId,
				Arguments: arguments,
			}); errMarshal == nil {
				if errPublish := nats_helper.PublishToNATS(queue, request); errPublish != nil {
					log.Printf("[StartUserMessageListener] ERROR in publish to %s:%s", queue, errPublish)
				}
			} else {
				log.Printf("[StartUserMessageListener] ERROR in publish to %s:%s", queue, errMarshal)
			}

		} else {
			log.Println("[StartNatsListener] Помилка: ID користувача чи текст повідомлення порожні")
		}
	}

	listener := &nats_helper.NatsListener{
		Handler: processor,
	}

	nats_helper.StartNatsListener("TELEGRAM_INPUT_TEXT_QUEUE", listener)
}

func parseNatsMessage(data []byte) (int64, string) {
	var msg TelegramUserNatsMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Printf("[StartNatsListener] Помилка при розборі повідомлення з NATS: %v", err)
		return 0, ""
	}

	return msg.ChatId, msg.Text
}

func findDataToAnotherProcessorRedirection(message string) (string, []string) {
	messageLength := len(message)
	for i := 0; i < messageLength; i++ {
		command := message[:messageLength-i]
		queue, flag := queue_processor.FindQueueByMessage(command)
		if flag {
			log.Printf("[StartUserMessageListener] Queue \"%s\" found for \"%s\"", queue, message)
			args := prepareArguments(message, command)
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
