package main

import (
	"fmt"
	"mattermost-botpoll/bot"
	"mattermost-botpoll/config"
)

func main() {
	configSetting, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Ошибка при получении настроек:", err)
		return
	}

	wsConn, botUserID := bot.InitConnection(configSetting)
	bot.ListenEvent(wsConn, botUserID, configSetting)
}
