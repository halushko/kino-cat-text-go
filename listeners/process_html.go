package listeners

import (
	"github.com/halushko/kino-cat-core-go/nats_helper"
	"log"
	"strings"
)

func ProcessHtml() {
	processor := func(data []byte) {
		id, args, err := nats_helper.ParseNatsBotCommand(data)
		if err != nil {
			log.Printf("[ProcessHtml] Проблема при парсингу повідомлення: %s", err)
			return
		}

		log.Printf("[ProcessHtml] Користувач: %d, текст: %s", id, args)

		for i, arg := range args {
			if strings.HasPrefix(arg, "https://utp.to/torrents/") {
				processUtopia(i, args, id)
				break
			}
		}

	}

	listener := &nats_helper.NatsListenerHandler{
		Function: processor,
	}

	err := nats_helper.StartNatsListener("PROCESS_HTTP_QUEUE", listener)
	if err != nil {
		log.Printf("[StartUserMessageListener] Помилка: %v", err)
	}
}

func processUtopia(i int, args []string, id int64) {
	parts := strings.Split(args[i], "/")
	if len(parts) > 0 {
		nats_helper.PublishCommandMessage("UTOPIA_GET_TORRENT_INFO", id, []string{parts[len(parts)-1]})
	} else {
		log.Printf("[StartUserMessageListener] Помилка - не можу знайти ID торента для Utopia: %s", args)
	}
}
