package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Dburl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name,omitempty"`
}

const configFileName = "gatorconfig.json"

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func (c *Config) SetUser(name string) error {
	c.CurrentUserName = name

	jsonConfig, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}

	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	return os.WriteFile(path, jsonConfig, 0o644)
}

func getConfigFilePath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return dir + "/." + configFileName, nil
}
