package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mattermost-botpoll/commands"
	"mattermost-botpoll/config"

	"github.com/mattermost/mattermost-server/v6/model"
)

func main() {
	configSetting, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Ошибка при получении настроек:", err)
		return
	}

	client := model.NewAPIv4Client(fmt.Sprintf("http://%s:%s", configSetting.MattermostServerIp, configSetting.MattermostServerPort))
	client.SetToken(configSetting.MattermostToken)

	user, _, err := client.GetMe("")
	if err != nil {
		log.Fatalf("Ошибка получения пользователя: %v", err)
	}
	botUserID := user.Id
	log.Printf("Бот запущен, ID: %s", botUserID)

	// Подключаемся к WebSocket
	wsClient, err := model.NewWebSocketClient4(fmt.Sprintf("ws://%s:%s/api/v4/websocket", configSetting.MattermostServerIp, configSetting.MattermostServerPort), configSetting.MattermostToken)
	// wsClient, err := model.NewWebSocketClient4("ws://localhost:8065/api/v4/websocket", configSetting.MattermostToken)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			log.Fatalf("Ошибка WebSocket: %s, Status: %d, Details: %s", appErr.Message, appErr.StatusCode, appErr.DetailedError)
		}
		log.Fatalf("Ошибка подключения к WebSocket: %v", err)
	}
	defer wsClient.Close()

	// Слушаем события
	wsClient.Listen()

	for event := range wsClient.EventChannel {
		if event.EventType() == model.WebsocketEventPosted {
			// Парсим пост из JSON
			postData := event.GetData()["post"]
			var post model.Post
			err := json.Unmarshal([]byte(postData.(string)), &post)
			if err != nil {
				log.Printf("Ошибка декодирования поста: %v", err)
				continue
			}

			// Проверяем, что сообщение не от бота
			if post.UserId != botUserID {
				handlePost(client, &post)
			}
		}
	}
}

func handlePost(client *model.Client4, post *model.Post) {
	if post.Message == "!hello" {
		commands.Hello(client, post)
	}
}
