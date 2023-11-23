package config

import (
	"encoding/json"
	"log"
	"os"
)

// Config represents the bot configuration.
type Config struct {
	Token     string   `json:"token"`
	MySQLDSN  string   `json:"mysqldsn"`
	Owners    []string `json:"owners"`
	GuildID   string   `json:"guildid"`
	KitchenID string   `json:"kitchenid"`
}

var config *Config

// LoadConfig reads the configuration from a JSON file and returns a Config struct.
func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)

	if err != nil {
		log.Fatal(err)
	}

	return config, nil
}

func GetConfig() *Config {
	return config
}
