package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	MattermostServerIp   string
	MattermostServerPort string
	MattermostToken      string
}

func LoadConfig() (Config, error) {
	viper.SetConfigFile("envs/.bot.env")
	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, fmt.Errorf("ошибка при чтении .env файла: %w", err)
	}

	config := Config{
		MattermostServerIp:   viper.Get("MATTERMOST_SERVER_IP").(string),
		MattermostServerPort: viper.Get("MATTERMOST_SERVER_PORT").(string),
		MattermostToken:      viper.Get("MATTERMOST_TOKEN").(string),
	}

	if config.MattermostServerIp == "" {
		return Config{}, fmt.Errorf("переменная MATTERMOST_SERVER_IP не задана в .env файле")
	}
	if config.MattermostServerPort == "" {
		return Config{}, fmt.Errorf("переменная MATTERMOST_SERVER_PORT не задана в .env файле")
	}
	if config.MattermostToken == "" {
		return Config{}, fmt.Errorf("переменная MATTERMOST_TOKEN не задана в .env файле")
	}
	return config, nil
}
