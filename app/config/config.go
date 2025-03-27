package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	MattermostServerIp     string
	MattermostServerPort   string
	MattermostToken        string
	MattermostSeverBaseUrl string
	MattermostServerWsUrl  string
}

func LoadConfig() (Config, error) {
	viper.SetConfigFile("envs/.env")
	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, fmt.Errorf("ошибка при чтении .env файла: %w", err)
	}

	config := Config{
		MattermostServerIp:     viper.GetString("MATTERMOST_SERVER_IP"),
		MattermostServerPort:   viper.GetString("MATTERMOST_SERVER_PORT"),
		MattermostToken:        viper.GetString("MATTERMOST_TOKEN"),
		MattermostSeverBaseUrl: fmt.Sprintf("http://%s:%s/api/v4", viper.GetString("MATTERMOST_SERVER_IP"), viper.GetString("MATTERMOST_SERVER_PORT")),
		MattermostServerWsUrl:  fmt.Sprintf("ws://%s:%s/api/v4/websocket", viper.GetString("MATTERMOST_SERVER_IP"), viper.GetString("MATTERMOST_SERVER_PORT")),
	}

	if config.MattermostServerIp == "" {
		return Config{}, fmt.Errorf("переменная MATTERMOST_SERVER_IP не задана в envs/.bot.env файле")
	}
	if config.MattermostServerPort == "" {
		return Config{}, fmt.Errorf("переменная MATTERMOST_SERVER_PORT не задана в envs/.bot.env файле")
	}
	if config.MattermostToken == "" {
		return Config{}, fmt.Errorf("переменная MATTERMOST_TOKEN не задана в envs/.bot.env файле")
	}
	return config, nil
}
