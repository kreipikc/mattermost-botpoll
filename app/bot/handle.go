package bot

import (
	"encoding/json"
	"fmt"
	"log"
	"mattermost-botpoll/commands"
	"mattermost-botpoll/config"
	"mattermost-botpoll/database"
	"mattermost-botpoll/models"
	"mattermost-botpoll/utils"
	"strings"

	"github.com/gorilla/websocket"
)

func ListenEvent(wsConn *websocket.Conn, dbConn *database.DB, botUserID string, config *config.Config) {
	defer wsConn.Close()

	for {
		var event map[string]interface{}
		err := wsConn.ReadJSON(&event)
		if err != nil {
			log.Printf("Ошибка чтения сообщения: %v", err)
			break
		}
		log.Printf("Получено событие: %+v\n", event)

		if event["event"] == "posted" {
			postData := event["data"].(map[string]interface{})["post"]
			var post models.Post
			err := json.Unmarshal([]byte(postData.(string)), &post)
			if err != nil {
				log.Printf("Ошибка декодирования поста: %v", err)
				continue
			}
			if post.UserId != botUserID {
				handlePost(dbConn, config.MattermostConf.MattermostSeverBaseUrl, config.MattermostConf.MattermostToken, &post)
			}
		}
	}
}

func handlePost(dbConn *database.DB, baseURL, token string, post *models.Post) {
	if post.Message == "!hello" {
		err := commands.Hello(baseURL, token, post)
		if err != nil {
			utils.SendResponse(baseURL, token, post, fmt.Sprintf("Ошибка при обработки команды !hello: %v", err))
		}
	}
	if strings.HasPrefix(post.Message, "!poll") {
		err := commands.CreatePoll(dbConn, baseURL, token, post)
		if err != nil {
			utils.SendResponse(baseURL, token, post, fmt.Sprintf("Ошибка при обработки команды !poll: %v", err))
		}
	}
}
