package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type MattermostConfig struct {
	MattermostServerIp     string
	MattermostServerPort   string
	MattermostToken        string
	MattermostSeverBaseUrl string
	MattermostServerWsUrl  string
}

type TarantoolConfig struct {
	TarantoolServerIp   string
	TarantoolServerPort string
	TarantoolUser       string
	TarantoolPassword   string
}

type Config struct {
	MattermostConf MattermostConfig
	TarantoolConf  TarantoolConfig
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile("envs/.env")
	err := viper.ReadInConfig()
	if err != nil {
		return &Config{}, fmt.Errorf("ошибка при чтении .env файла: %w", err)
	}

	config := &Config{
		MattermostConf: MattermostConfig{
			MattermostServerIp:     viper.GetString("MATTERMOST_SERVER_IP"),
			MattermostServerPort:   viper.GetString("MATTERMOST_SERVER_PORT"),
			MattermostToken:        viper.GetString("MATTERMOST_TOKEN"),
			MattermostSeverBaseUrl: fmt.Sprintf("http://%s:%s/api/v4", viper.GetString("MATTERMOST_SERVER_IP"), viper.GetString("MATTERMOST_SERVER_PORT")),
			MattermostServerWsUrl:  fmt.Sprintf("ws://%s:%s/api/v4/websocket", viper.GetString("MATTERMOST_SERVER_IP"), viper.GetString("MATTERMOST_SERVER_PORT")),
		},
		TarantoolConf: TarantoolConfig{
			TarantoolServerIp:   viper.GetString("TARANTOOL_SERVER_IP"),
			TarantoolServerPort: viper.GetString("TARANTOOL_SERVER_PORT"),
			TarantoolUser:       viper.GetString("TARANTOOL_USER"),
			TarantoolPassword:   viper.GetString("TARANTOOL_PASSWORD"),
		},
	}

	err = checkDataConfig(config)
	if err != nil {
		return &Config{}, err
	}

	return config, nil
}

func checkDataConfig(config *Config) error {
	// Mattermost
	if config.MattermostConf.MattermostServerIp == "" {
		return fmt.Errorf("переменная MATTERMOST_SERVER_IP не задана в envs/.env файле")
	}
	if config.MattermostConf.MattermostServerPort == "" {
		return fmt.Errorf("переменная MATTERMOST_SERVER_PORT не задана в envs/.env файле")
	}
	if config.MattermostConf.MattermostToken == "" {
		return fmt.Errorf("переменная MATTERMOST_TOKEN не задана в envs/.env файле")
	}

	// Tarantool
	if config.TarantoolConf.TarantoolServerIp == "" {
		return fmt.Errorf("переменная TARANTOOL_SERVER_IP не задана в envs/.env")
	}
	if config.TarantoolConf.TarantoolServerPort == "" {
		return fmt.Errorf("переменная TARANTOOL_SERVER_PORT не задана в envs/.env")
	}
	if config.TarantoolConf.TarantoolUser == "" {
		return fmt.Errorf("переменная TARANTOOL_USER не задана в envs/.env")
	}

	return nil
}
