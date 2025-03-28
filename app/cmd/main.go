package main

import (
	"log"
	"mattermost-botpoll/bot"
	"mattermost-botpoll/config"
	"mattermost-botpoll/database"
)

func main() {
	configSetting, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка при получении настроек: %v", err)
	}

	dbConn, err := database.InitConnectionDB(configSetting)
	if err != nil {
		log.Fatalf("Ошибка при подключении к Tarantool: %v", err)
	}
	dbConn.InitSpaces()

	wsConn, botUserID := bot.InitConnection(configSetting)
	bot.ListenEvent(wsConn, dbConn, botUserID, configSetting)
}
