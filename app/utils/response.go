package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mattermost-botpoll/models"
	"net/http"
)

func SendResponse(baseURL string, token string, post *models.Post, message string) error {
	reply := map[string]interface{}{
		"channel_id": post.ChannelId,
		"message":    message,
	}
	replyData, _ := json.Marshal(reply)

	req, err := http.NewRequest("POST", baseURL+"/posts", bytes.NewBuffer(replyData))
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка отправки ответа: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ошибка отправки ответа: %d, тело: %s", resp.StatusCode, string(body))
	}
	return nil
}
