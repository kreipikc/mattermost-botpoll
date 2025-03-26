package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func SendResponse(baseURL string, token string, replyData []byte) error {
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
