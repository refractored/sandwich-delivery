package config

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

// Config represents the bot configuration.
type Config struct {
	Token   string   `json:"token"`
	Owners  []string `json:"owners"`
	GuildID string   `json:"guildID"`

	StartupChannelID string `json:"startupChannelID"`

	Database struct {
		Host         string            `json:"host"`
		Port         int               `json:"port"`
		User         string            `json:"user"`
		Password     string            `json:"password"`
		DBName       string            `json:"database"`
		ExtraOptions map[string]string `json:"extraOptions"`

		URL string `json:"url,omitempty"`
	} `json:"database"`

	TokensPerOrder *uint32 `json:"tokensPerOrder"`
	DailyTokens    *uint32 `json:"dailyTokens"`
}

// Error represents a configuration error.
type Error struct {
	message string
}

func (e *Error) Error() string {
	return e.message
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

	setDefaults(config)

	return config, nil
}

// VerifyConfig returns true if valid, otherwise false with an error.
func VerifyConfig(config *Config) (bool, error) {
	if config.Token == "" {
		return false, &Error{"token is not set"}
	}

	if len(config.Owners) == 0 {
		return false, &Error{"There are no entries in owners"}
	}

	if config.GuildID == "" {
		return false, &Error{"guildID is not set"}
	}

	if config.Database.Host == "" {
		return false, &Error{"Database host is not set"}
	}

	if config.Database.Port == 0 {
		return false, &Error{"Database port is not set"}
	}

	if config.Database.User == "" {
		return false, &Error{"Database user is not set"}
	}

	if config.Database.DBName == "" {
		return false, &Error{"Database name is not set"}
	}

	for key, value := range config.Database.ExtraOptions {
		if key == "" || value == "" {
			return false, &Error{"Database extra option is not set"}
		}

		if strings.ContainsAny(key, " :+?/()=@&") {
			return false, &Error{"Database extra option key contains invalid characters: " + key}
		}

		if strings.ContainsAny(value, " :+?/()=@&") {
			return false, &Error{"Database extra option key contains invalid characters: " + value}
		}
	}

	if config.Database.URL != "" {
		if config.Database.Host != "" || config.Database.Port != 0 || config.Database.User != "" || config.Database.Password != "" || config.Database.DBName != "" || len(config.Database.ExtraOptions) != 0 {
			return false, &Error{"Database URL and other database options are set"}
		}
	}

	if *config.TokensPerOrder < 0 {
		return false, &Error{"Tokens per order is less than 0"}
	}

	if *config.DailyTokens < 0 {
		return false, &Error{"Daily tokens is less than 0"}
	}

	return true, nil
}

func setDefaults(config *Config) {
	if config.TokensPerOrder == nil {
		var tokensPerOrder uint32 = 1
		config.TokensPerOrder = &tokensPerOrder
	}

	if config.DailyTokens == nil {
		var dailyTokens uint32 = 1
		config.DailyTokens = &dailyTokens
	}
}

func GetConfig() *Config {
	return config
}
