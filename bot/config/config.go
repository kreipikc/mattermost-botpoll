package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	MattermostServerUrl string
	MattermostToken     string
	TarantoolIP         string
	TarantoolPort       string
}

func LoadConfig() (Config, error) {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, fmt.Errorf("ошибка при чтении .env файла: %w", err)
	}

	config := Config{
		MattermostServerUrl: viper.Get("MATTERMOST_SERVER_URL").(string),
		MattermostToken:     viper.Get("MATTERMOST_TOKEN").(string),
		TarantoolIP:         viper.Get("TARANTOOL_IP").(string),
		TarantoolPort:       viper.Get("TARANTOOL_PORT").(string),
	}

	if config.MattermostServerUrl == "" {
		return Config{}, fmt.Errorf("переменная MATTERMOST_SERVER_URL не задана в .env файле")
	}
	if config.MattermostToken == "" {
		return Config{}, fmt.Errorf("переменная MATTERMOST_TOKEN не задана в .env файле")
	}
	if config.TarantoolIP == "" {
		return Config{}, fmt.Errorf("переменная TARANTOOL_IP не задана в .env файле")
	}
	if config.TarantoolPort == "" {
		return Config{}, fmt.Errorf("переменная TARANTOOL_PORT не задана в .env файле")
	}
	return config, nil
}
