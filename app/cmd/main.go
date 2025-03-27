package main

import (
	"log"
	"mattermost-botpoll/bot"
	"mattermost-botpoll/config"
)

func main() {
	configSetting, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка при получении настроек: %v", err)
	}

	wsConn, botUserID := bot.InitConnection(configSetting)
	bot.ListenEvent(wsConn, botUserID, configSetting)
}
