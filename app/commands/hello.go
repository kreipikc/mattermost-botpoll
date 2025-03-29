package commands

import (
	"fmt"
	"log"
	"mattermost-botpoll/models"
	"mattermost-botpoll/utils"
)

func Hello(baseURL string, token string, post *models.Post) error {
	message := fmt.Sprintf("Привет, я бот! Ты написал: %s", post.Message)

	err := utils.SendResponse(baseURL, token, post, message)
	if err != nil {
		return fmt.Errorf("ошибка формирования или отправки ответа: %v", err)
	}

	log.Println("Сообщение на !hello отправлено успешно")

	return nil
}
