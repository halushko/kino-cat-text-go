package listeners

import (
	"encoding/json"
	"github.com/halushko/kino-cat-core-go/nats_helper"
	"github.com/nats-io/nats.go"
	"kino-cat-text-go/queue_processor"
	"log"
)

func StartGetHelpCommandListener() {
	processor := func(msg *nats.Msg) {
		log.Printf("[StartGetHelpCommandListener] Отримано повідомлення з NATS: %s", string(msg.Data))
		chatId, messageText := parseNatsMessage(msg.Data)

		log.Printf("[StartGetHelpCommandListener] Парсинг повідомлення: chatID = %d, message = %s", chatId, messageText) // Новый лог для проверки данных

		if chatId != 0 {

			jsonData, err := json.Marshal(TelegramUserNatsMessage{
				ChatId: chatId,
				Text:   messageText,
			})
			if err != nil {
				log.Printf("[StartGetHelpCommandListener] ERROR:%s", err)
				return
			}
			if err = nats_helper.PublishToNATS("TELEGRAM_OUTPUT_TEXT_QUEUE", jsonData); err != nil {
				log.Printf("[StartGetHelpCommandListener] ERROR:%s", err)
				return
			}

			commands, order := queue_processor.GetAllDescriptions()
			result := ""
			for _, value := range order {
				result = result + value + " - " + commands[value] + "\n"
			}
			queue := "TELEGRAM_OUTPUT_TEXT_QUEUE"
			if request, errMarshal := json.Marshal(TelegramUserNatsMessage{
				ChatId: chatId,
				Text:   result,
			}); errMarshal == nil {
				if errPublish := nats_helper.PublishToNATS(queue, request); errPublish != nil {
					log.Printf("[StartUserMessageListener] ERROR in publish to %s:%s", queue, errPublish)
				}
			} else {
				log.Printf("[StartGetHelpCommandListener] ERROR in publish to %s:%s", queue, errMarshal)
			}

		} else {
			log.Println("[StartGetHelpCommandListener] Помилка: ID користувача чи текст повідомлення порожні")
		}
	}

	listener := &nats_helper.NatsListener{
		Handler: processor,
	}

	nats_helper.StartNatsListener("DISPLAY_ALL_COMMANDS", listener)
}
