package commands

import (
	"github.com/bwmarrin/discordgo"
	"go-discord-bot/src/config"
	"strings"
)

//var commandRegistry = map[string]func(*discordgo.Session, *discordgo.MessageCreate){
//	"coinflip": CoinflipCommand,
//	"ping":     PingCommand,
//}

func HandleCommand(s *discordgo.Session, m *discordgo.MessageCreate, cfg *config.Config) {
	args := strings.Split(m.Content, " ")

	prefix := cfg.Prefix

	if args[0] != prefix {
		return
	}

	// Get the command name
	commandName := args[1]

	// Check if the command exists in the registry
	if commandFunc, ok := commandRegistry[commandName]; ok {
		commandFunc(s, m)
	}
}
