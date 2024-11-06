package config

import (
	"encoding/json"
	"os"
)

const configFilePath = ".gatorconfig.json"

type Config struct {
	DbURL           string  `json:"db_url"`
	CurrentUserName *string `json:"current_user_name"`
}

func getConfigFilePAth() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return homeDir + "/" + configFilePath, nil
}

func (c *Config) SetUser(userName string) error {
	c.CurrentUserName = &userName

	return write(c)
}

func Read() (Config, error) {
	configFilePath, err := getConfigFilePAth()
	if err != nil {
		return Config{}, err
	}

	content, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func write(cfg *Config) error {
	configFilePath, err := getConfigFilePAth()
	if err != nil {
		return err
	}

	content, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(configFilePath, content, 0644)
	if err != nil {
		return err
	}

	return nil
}
