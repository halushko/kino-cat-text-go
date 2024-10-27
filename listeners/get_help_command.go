package listeners

import (
	"fmt"
	"github.com/halushko/kino-cat-core-go/nats_helper"
	"kino-cat-text-go/queue_processor"
	"log"
	"strings"
)

func StartGetHelpCommandListener() {
	processor := func(data []byte) {
		userId, message, err := nats_helper.ParseNatsBotText(data)
		if err != nil {
			log.Printf("[StartGetHelpCommandListener] ERROR: %v", err)
			return
		}
		log.Printf("[StartGetHelpCommandListener] Отримано повідомлення: \"%s\" з NATS від користувача: %d", message, userId)

		if userId != 0 {
			commands, order := queue_processor.GetAllDescriptions()

			var sb strings.Builder
			for _, value := range order {
				sb.WriteString(fmt.Sprintf("%s - %s\n", value, commands[value]))
			}
			result := sb.String()

			if err = nats_helper.PublishTextMessage("TELEGRAM_OUTPUT_TEXT_QUEUE", userId, result); err != nil {
				log.Printf("[StartUserMessageListener] Не вдалося надіслати повідомлення \"%s\" через Телеграм бот", result)
				return
			}
		} else {
			log.Printf("[StartGetHelpCommandListener] Помилка: ID користувача чи текст повідомлення порожні")
		}
	}

	listener := &nats_helper.NatsListenerHandler{
		Function: processor,
	}

	if err := nats_helper.StartNatsListener("DISPLAY_ALL_COMMANDS", listener); err != nil {
		log.Printf("[StartGetHelpCommandListener] Не вдалося почати роботу над обробкою команди /help")
	}
}
