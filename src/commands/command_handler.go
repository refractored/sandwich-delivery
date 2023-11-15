package commands

import (
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
	"sandwich-delivery/src/config"
	"sandwich-delivery/src/models"
	"strings"
)

func IsUserBlacklisted(db *gorm.DB, userID string) bool {
	var user models.BlacklistUser

	result := db.Select("user_id").Where("user_id = ?", userID).First(&user)

	return result.Error == nil
}

func DisplayName(s *discordgo.Session, m *discordgo.MessageCreate) string {
	var displayname string
	if m.Author.Discriminator != "0" {
		displayname = m.Author.Username + "#" + m.Author.Discriminator
	} else {
		displayname = m.Author.Username
	}
	return displayname
}
func HandleCommand(s *discordgo.Session, m *discordgo.MessageCreate, cfg *config.Config, db *gorm.DB) {

	if !strings.HasPrefix(m.Content, cfg.Prefix) {
		return
	}
	test := strings.Replace(m.Content, cfg.Prefix, "", -1)
	args := strings.Fields(test)

	if IsUserBlacklisted(db, m.Author.ID) {
		s.ChannelMessageSend(m.ChannelID, "You are blacklisted from using this bot!")
		return
	}
	commandName := args[0]

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
		"delorder": func(s *discordgo.Session, m *discordgo.MessageCreate) {
			DelOrderCommand(s, m, db)
		},
	}

	if commandFunc, ok := commandRegistry[commandName]; ok {
		commandFunc(s, m)
	}
}
