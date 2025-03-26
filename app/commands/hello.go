package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"mattermost-botpoll/models"
	"mattermost-botpoll/utils"
)

func Hello(baseURL string, token string, post *models.Post) {
	reply := map[string]interface{}{
		"channel_id": post.ChannelId,
		"message":    fmt.Sprintf("Привет, я бот! Ты написал: %s", post.Message),
	}
	replyData, _ := json.Marshal(reply)

	err := utils.SendResponse(baseURL, token, replyData)
	if err != nil {
		log.Fatalf("Ошибка формирования или отправки ответа: %v", err)
	}
}
