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
		MattermostServerIp:     viper.Get("MATTERMOST_SERVER_IP").(string),
		MattermostServerPort:   viper.Get("MATTERMOST_SERVER_PORT").(string),
		MattermostToken:        viper.Get("MATTERMOST_TOKEN").(string),
		MattermostSeverBaseUrl: fmt.Sprintf("http://%s:%s/api/v4", viper.Get("MATTERMOST_SERVER_IP").(string), viper.Get("MATTERMOST_SERVER_PORT").(string)),
		MattermostServerWsUrl:  fmt.Sprintf("ws://%s:%s/api/v4/websocket", viper.Get("MATTERMOST_SERVER_IP").(string), viper.Get("MATTERMOST_SERVER_PORT").(string)),
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
