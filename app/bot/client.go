package bot

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mattermost-botpoll/config"
	"mattermost-botpoll/models"
	"net/http"

	"github.com/gorilla/websocket"
)

func InitConnection(config config.Config) (*websocket.Conn, string) {
	user, err := getMe(config.MattermostSeverBaseUrl, config.MattermostToken)
	if err != nil {
		log.Fatalf("Ошибка получения пользователя: %v", err)
	}
	botUserID := user.Id
	log.Printf("Бот запущен, ID: %s", botUserID)

	log.Printf("Попытка подключения к WebSocket: %s", config.MattermostServerWsUrl)

	dialer := websocket.DefaultDialer
	header := map[string][]string{
		"Authorization": {"Bearer " + config.MattermostToken},
	}

	wsConn, _, err := dialer.Dial(config.MattermostServerWsUrl, header)
	if err != nil {
		log.Fatalf("Ошибка подключения к WebSocket: %v", err)
	}

	log.Println("WebSocket подключён успешно")

	authMessage := map[string]interface{}{
		"seq":    1,
		"action": "authentication_challenge",
		"data": map[string]string{
			"token": config.MattermostToken,
		},
	}
	err = wsConn.WriteJSON(authMessage)
	if err != nil {
		log.Fatalf("Ошибка отправки сообщения аутентификации: %v", err)
	}

	return wsConn, botUserID
}

func getMe(baseURL string, token string) (*models.User, error) {
	req, err := http.NewRequest("GET", baseURL+"/users/me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ошибка: %d, тело: %s", resp.StatusCode, string(body))
	}

	var user models.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	return &user, err
}
