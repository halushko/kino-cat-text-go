package queue_processor

import (
	"sort"
	"unicode/utf8"
)

var queues = map[string][]string{
	"/space":                 {"FILE_MOVE_TO_FOLDER"},
	"/start_":                {"FILE_SHOW_FREE_SPACE"},
	"/list":                  {"EXECUTE_TORRENT_COMMAND_LIST", "відобразити всі торенти, можна запускати для конкретного сховища"},
	"/list_":                 {},
	"/more_":                 {"EXECUTE_TORRENT_COMMAND_SHOW_COMMANDS", "перелічити можливі команди для вказаного торента"},
	"/resume_":               {"EXECUTE_TORRENT_COMMAND_RESUME_TORRENT", "продовжити закачувати цей торент"},
	"/pause_":                {"EXECUTE_TORRENT_COMMAND_PAUSE_TORRENT", "тимчасово припинити закачувати цей торент"},
	"/resume_all":            {"EXECUTE_TORRENT_COMMAND_RESUME_ALL_TORRENTS", "продовжити закачувати всі торенти, можна запускати для конкретного сховища"},
	"/resume_all_":           {},
	"/pause_all":             {"EXECUTE_TORRENT_COMMAND_PAUSE_ALL_TORRENTS", "тимчасово припинити закачувати всі торенти, можна запускати для конкретного сховища"},
	"/pause_all_":            {},
	"/info_":                 {"EXECUTE_TORRENT_COMMAND_INFO", "відобразити інформацію по торенту"},
	"/approve_with_files_":   {"EXECUTE_TORRENT_COMMAND_DELETE_WITH_FILES", "видалити вказаний торент разом з файлами"},
	"/approve_just_torrent_": {"EXECUTE_TORRENT_COMMAND_DELETE_ONLY_TORRENT", "видалити вказаний торент, але залишити файли"},
	"/files_":                {"EXECUTE_TORRENT_COMMAND_LIST_FILES", "відобразити всі файли, що будуть скачані в цьому торенті"},
	"/remove_":               {"EXECUTE_TORRENT_COMMAND_DELETE_ONLY_TORRENT", "видалення вказаного торента"},
	"/downloads":             {"EXECUTE_LIST_TORRENTS_IN_DOWNLOAD_STATUS", "відобразити всі торенти що знаходяться в стані \"завантаження\""},
	"/help":                  {"DISPLAY_ALL_COMMANDS", "вивести інформацію по всім командам"},
	"":                       {"EXECUTE_TORRENT_COMMAND_SEARCH_BY_NAME"},
	"<якийсь текст>":         {"", "пошук торентів по частині назві"},
}

func FindQueueByMessage(message string) (string, bool) {
	for key := range queues {
		if key == message {
			if len(queues[key]) == 0 {
				if len(message) > 0 {
					_, size := utf8.DecodeLastRuneInString(message)
					correctedMessage := message[:len(message)-size]
					return FindQueueByMessage(correctedMessage)
				} else {
					return "", false
				}
			}
			return queues[key][0], true
		}
	}
	return "", false
}

func GetAllDescriptions() (map[string]string, []string) {
	result := make(map[string]string)
	sortedKeys := make([]string, 0)

	for key, value := range queues {
		if len(value) > 1 && len(key) > 0 {
			result[key] = value[1]
			sortedKeys = append(sortedKeys, key)
		}
	}

	sort.Slice(sortedKeys, func(i, j int) bool { return sortedKeys[i] < sortedKeys[j] })
	return result, sortedKeys
}
