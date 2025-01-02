package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configURL = ".gatorconfig.json"

func ReadConfig() (Config, error) {
	jsonFile, err := os.Open(configURL)
	if err != nil {
		return Config{}, err
	}
	defer jsonFile.Close()

	var config Config
	decoder := json.NewDecoder(jsonFile)
	if err := decoder.Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func SetUser(cfg *Config, newUserName string) error {
	cfg.CurrentUserName = newUserName

	if err := write(*cfg); err != nil {
		return err
	}

	return nil
}

func write(cfg Config) error {
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(configURL, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}
