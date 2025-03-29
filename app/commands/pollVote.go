package commands

import (
	"fmt"
	"mattermost-botpoll/database"
	"mattermost-botpoll/models"
	"regexp"
	"strconv"
	"strings"
)

func PollVote(dbConn *database.DB, baseURL string, token string, post *models.Post) error {
	id_poll, variant, err := parseVotePoll(post.Message)
	if err != nil {
		return fmt.Errorf("ошибка при валидации данных: %s", err)
	}

	err = dbConn.UpdatePollVote(id_poll, variant)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении данных в Tarantool: %s", err)
	}
	return nil
}

func parseVotePoll(command string) (int, string, error) {
	re := regexp.MustCompile(`!vote_poll\s+(\d+)\s+(.+)`)

	matches := re.FindStringSubmatch(command)
	if matches == nil {
		return 0, "", fmt.Errorf("неверный формат команды: ожидается '!vote_poll <id_poll> <variant>'")
	}

	idPollStr := strings.TrimSpace(matches[1])
	idPoll, err := strconv.Atoi(idPollStr)
	if err != nil {
		return 0, "", fmt.Errorf("ошибка преобразования id_poll '%s' в int: %v", idPollStr, err)
	}

	variant := matches[2]

	return idPoll, variant, nil
}
