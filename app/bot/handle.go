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
		// log.Printf("Получено событие: %+v\n", event)

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
	if strings.TrimSpace(post.Message) == "!hello" {
		err := commands.Hello(baseURL, token, post)
		if err != nil {
			log.Printf("Ошибка при обработки команды !hello: %v", err)
			utils.SendResponse(baseURL, token, post, fmt.Sprintf("Ошибка при обработки команды !hello: %v", err))
		}
	}
	if strings.TrimSpace(post.Message) == "!help" {
		err := commands.HelpCommandPoll(baseURL, token, post)
		if err != nil {
			log.Printf("Ошибка при обработки команды !help: %v", err)
			utils.SendResponse(baseURL, token, post, fmt.Sprintf("Ошибка при обработки команды !help: %v", err))
		}
	}
	if strings.HasPrefix(post.Message, "!poll") {
		err := commands.CreatePoll(dbConn, baseURL, token, post)
		if err != nil {
			log.Printf("Ошибка при обработки команды !poll: %v", err)
			utils.SendResponse(baseURL, token, post, fmt.Sprintf("Ошибка при обработки команды !poll: %v", err))
		}
	}
	if strings.HasPrefix(post.Message, "!vote_poll") {
		err := commands.PollVote(dbConn, baseURL, token, post)
		if err != nil {
			log.Printf("Ошибка при обработки команды !vote_poll: %v", err)
			utils.SendResponse(baseURL, token, post, fmt.Sprintf("Ошибка при обработки команды !vote_poll: %v", err))
			return
		}
		utils.SendResponse(baseURL, token, post, "Vote success")
	}
	if strings.HasPrefix(post.Message, "!info_poll") {
		poll, err := commands.GetInfo(dbConn, baseURL, token, post)
		if err != nil {
			log.Printf("Ошибка при обработки команды !info_poll: %v", err)
			utils.SendResponse(baseURL, token, post, fmt.Sprintf("Ошибка при обработки команды !info_poll: %v", err))
			return
		}
		message := fmt.Sprintf("Опрос:\nId: %d\nTitle: %s\nDescription: %s\nDate end: %s\nVariants: %v\nAuthorID: %s", poll.Id, poll.Title, poll.Description, poll.DateEnd, poll.Variants, poll.AuthorID)
		utils.SendResponse(baseURL, token, post, message)
	}
	if strings.HasPrefix(post.Message, "!end_poll") {
		err := commands.EndPoll(dbConn, baseURL, token, post)
		if err != nil {
			log.Printf("Ошибка при обработки команды !end_poll: %v", err)
			utils.SendResponse(baseURL, token, post, fmt.Sprintf("Ошибка при обработки команды !end_poll: %v", err))
			return
		}
		utils.SendResponse(baseURL, token, post, "Голосование звершено!")
	}
	if strings.HasPrefix(post.Message, "!delete_poll") {
		err := commands.DeletePoll(dbConn, baseURL, token, post)
		if err != nil {
			log.Printf("Ошибка при обработки команды !delete_poll: %v", err)
			utils.SendResponse(baseURL, token, post, fmt.Sprintf("Ошибка при обработки команды !delete_poll: %v", err))
			return
		}
		utils.SendResponse(baseURL, token, post, "Голосование удалено!")
	}
}
