package config

import (
	"encoding/json"
	"os"
)

// Config represents the bot configuration.
type Config struct {
	Token        string `json:"token"`
	Prefix       string `json:"prefix"`
	SupportGuild string `json:"supportguild"`
	MySQLDSN     string `json:"mysqldsn"`
	ownerID      string `json:"ownerid"`
}

// LoadConfig reads the configuration from a JSON file and returns a Config struct.
func LoadConfig(filePath string) (Config, error) {
	var config Config

	file, err := os.Open(filePath)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}
