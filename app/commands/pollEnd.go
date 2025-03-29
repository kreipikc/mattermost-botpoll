package commands

import (
	"fmt"
	"log"
	"mattermost-botpoll/database"
	"mattermost-botpoll/models"
	"regexp"
	"strconv"
)

func EndPoll(dbConn *database.DB, baseURL string, token string, post *models.Post) error {
	idPoll, err := parseEndPollMessage(post.Message)
	if err != nil {
		return fmt.Errorf("ошибка при валидации данных: %s", err)
	}

	err = dbConn.UpdatePollEnd(idPoll, post.UserId)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении данных: %s", err)
	}

	log.Printf("Голосование ID: %d завершено успешно", idPoll)

	return nil
}

func parseEndPollMessage(message string) (int, error) {
	re := regexp.MustCompile(`!end_poll\s+(\d+)`)

	matches := re.FindStringSubmatch(message)
	if matches == nil {
		return 0, fmt.Errorf("неверный формат команды: ожидается '!end_poll <id_poll>'")
	}

	idPollStr := matches[1]
	idPoll, err := strconv.Atoi(idPollStr)
	if err != nil {
		return 0, fmt.Errorf("ошибка преобразования id_poll '%s' в int: %v", idPollStr, err)
	}

	return idPoll, nil
}
