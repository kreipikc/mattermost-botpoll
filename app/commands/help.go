package commands

import (
	"fmt"
	"log"
	"mattermost-botpoll/models"
	"mattermost-botpoll/utils"
)

func HelpCommandPoll(baseURL string, token string, post *models.Post) error {
	message := "Команды для работы с poll:\n1. !poll <title> <description> <date_end> <variants: obj1, obj2...> - создание голосования;\n2. !info_poll <id_poll> - вывод данных о голосовании;\n3. !vote_poll <id_poll> <variant> - проголосовать за вариант;\n4. !end_poll <id_poll> - завершение голосования (только автор);\n5. !delete_poll <id_poll> - удаление голосования (только автор)."

	err := utils.SendResponse(baseURL, token, post, message)
	if err != nil {
		return fmt.Errorf("ошибка формирования или отправки ответа: %v", err)
	}

	log.Println("Информация о командах отправлена успешно")

	return nil
}
