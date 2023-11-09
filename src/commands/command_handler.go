package commands

import (
	"github.com/bwmarrin/discordgo"
	"go-discord-bot/src/config"
	"go-discord-bot/src/models"
	"gorm.io/gorm"
	"strings"
)

func IsUserBlacklisted(db *gorm.DB, userID string) bool {
	var user models.BlacklistUser

	result := db.Select("user_id").Where("user_id = ?", userID).First(&user)

	return result.Error == nil
}

func HandleCommand(s *discordgo.Session, m *discordgo.MessageCreate, cfg *config.Config, db *gorm.DB) {

	args := strings.Split(m.Content, " ")

	prefix := cfg.Prefix

	if args[0] != prefix {
		return
	}
	if IsUserBlacklisted(db, m.Author.ID) {
		return
	}
	commandName := args[1]

	var commandRegistry = map[string]func(*discordgo.Session, *discordgo.MessageCreate){

		// Owner Only Commands
		"shutdown": ShutdownCommand,
		"coinflip": CoinflipCommand,
		"blacklist": func(s *discordgo.Session, m *discordgo.MessageCreate) {
			BlacklistCommand(s, m, db)
		},
		"unblacklist": func(s *discordgo.Session, m *discordgo.MessageCreate) {
			UnblacklistCommand(s, m, db)
		},
		// Support Server Only Commands

		// Everyone Commands
		"ping": PingCommand,
		"order": func(s *discordgo.Session, m *discordgo.MessageCreate) {
			OrderCommand(s, m, db)
		},
	}

	if commandFunc, ok := commandRegistry[commandName]; ok {
		commandFunc(s, m)
	}
}
