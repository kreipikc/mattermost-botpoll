package commands

import (
	"fmt"
	"mattermost-botpoll/database"
	"mattermost-botpoll/models"
	"regexp"
	"strconv"
)

func GetInfo(dbConn *database.DB, baseURL string, token string, post *models.Post) (*models.PollBody, error) {
	idPoll, err := parseInfoPoll(post.Message)
	if err != nil {
		return nil, fmt.Errorf("ошибка при валидации данных: %s", err)
	}

	poll, err := dbConn.GetPollByID(idPoll)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении данных из Tarantool: %s", err)
	}

	return poll, nil
}

func parseInfoPoll(message string) (int, error) {
	re := regexp.MustCompile(`!info_poll\s+(\d+)`)

	matches := re.FindStringSubmatch(message)
	if matches == nil {
		return 0, fmt.Errorf("неверный формат команды: ожидается '!info_poll <id_poll>'")
	}

	idPollStr := matches[1]
	idPoll, err := strconv.Atoi(idPollStr)
	if err != nil {
		return 0, fmt.Errorf("ошибка преобразования id_poll '%s' в int: %v", idPollStr, err)
	}

	return idPoll, nil
}
