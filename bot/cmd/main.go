package main

import (
	"fmt"
	"mattermost-botpoll/config"
)

func main() {
	configSetting, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Ошибка при получении настроек:", err)
		return
	}

	fmt.Println(configSetting.MattermostServerUrl, configSetting.MattermostToken)
	fmt.Println(configSetting.TarantoolIP, configSetting.TarantoolPort)
}
