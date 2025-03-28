package commands

import (
	"fmt"
	"mattermost-botpoll/database"
	"mattermost-botpoll/models"
	"mattermost-botpoll/utils"
	"regexp"
	"strings"
	"time"
)

func CreatePoll(dbConn *database.DB, baseURL string, token string, post *models.Post) error {
	poll, err := parsePollString(post.Message)
	if err != nil {
		return fmt.Errorf("ошибка при валидации команды: %v", err)
	}
	poll.AuthorID = post.UserId

	idPoll, err := dbConn.CreatePoll(poll)
	if err != nil {
		return fmt.Errorf("ошибка при создании опроса в Tarantool: %v", err)
	}

	message := fmt.Sprintf("Опрос:\nId: %d\nTitle: %s\nDescription: %s\nDate end: %s\nVariants: %v\nAuthorID: %s", idPoll, poll.Title, poll.Description, poll.DateEnd, poll.Variants, poll.AuthorID)

	err = utils.SendResponse(baseURL, token, post, message)
	if err != nil {
		return fmt.Errorf("ошибка формирования или отправки ответа: %s", err)
	}
	return nil
}

func parsePollString(pollString string) (*models.PollBody, error) {
	pattern := `!poll\s+(?P<title>.+?)\s+(?P<description>.+?)\s+(?P<date_end>.+?)\s+(?P<variants>.+)`
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(pollString)

	if match == nil {
		return nil, fmt.Errorf("строка не соответствует формату команды")
	}

	title := strings.TrimSpace(match[1])
	description := strings.TrimSpace(match[2])
	dateEndStr := strings.TrimSpace(match[3])
	variants := strings.TrimSpace(match[4])

	dateEnd, err := time.Parse("02.01.2006", dateEndStr)
	if err != nil {
		return nil, fmt.Errorf("дата конца не соответствует формату dd.mm.yyyy")
	}

	variantsList := strings.Split(variants, ",")
	variantsMap := make(map[string]int)
	for _, zn := range variantsList {
		variantsMap[strings.TrimSpace(zn)] = 0
	}

	return &models.PollBody{
		Title:       title,
		Description: description,
		Variants:    variantsMap,
		DateEnd:     dateEnd,
	}, nil
}
