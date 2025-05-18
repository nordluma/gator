package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() Config {
	configFile, err := getConfigFilePath()
	if err != nil {
		log.Fatal(err)
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Could read %s. Error: %s", configFile, err)
	}

	var config Config
	if err = json.Unmarshal(data, &config); err != nil {
		log.Fatal("Could not deserialize '.gatorconfig.json'")
	}

	return config
}

func (c *Config) SetUser(user string) {
	c.CurrentUserName = user
	if err := write(*c); err != nil {
		log.Fatal(err)
	}
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Could not get home dir: %s\n", err)
	}

	return filepath.Join(homeDir, configFileName), nil
}

func write(cfg Config) error {
	configFile, err := getConfigFilePath()
	if err != nil {
		return err
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("Could not serialize config. Error: %s", err)
	}

	if err = os.WriteFile(configFile, data, os.FileMode(os.O_TRUNC|os.O_CREATE)); err != nil {
		return fmt.Errorf("Could not write config to file. Error: %s", err)
	}

	return nil
}
