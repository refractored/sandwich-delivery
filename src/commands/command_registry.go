package commands

import (
	"github.com/bwmarrin/discordgo"
)

var commandRegistry = map[string]func(*discordgo.Session, *discordgo.MessageCreate){

	// Owner Only Commands
	"coinflip": CoinflipCommand,

	// Support Server Only Commands

	// Everyone Commands
	"ping": PingCommand,
}
