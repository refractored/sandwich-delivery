package commands

import (
	"github.com/bwmarrin/discordgo"
	"go-discord-bot/src/config"
	"log"
	"strings"
)

func HandleCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Split(m.Content, " ")

	configPath := "config.json"

	// I actually don't know if this is a bad approach to load the config twice, but it works for now.
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	prefix := cfg.Prefix
	supportserver := cfg.SupportGuild

	if args[0] != prefix {
		return
	}

	switch args[1] {
	case "coinflip":

		if m.GuildID != supportserver {
			return
		}
		CoinflipCommand(s, m)

	case "ping":
		PingCommand(s, m)
	}
}
