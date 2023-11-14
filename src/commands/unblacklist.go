package commands

import (
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
	"sandwich-delivery/src/models"
	"strings"
)

func UnblacklistCommand(s *discordgo.Session, m *discordgo.MessageCreate, db *gorm.DB) {
	args := strings.Split(m.Content, " ")

	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !unblacklist <user_id>")
		return
	}

	userID := args[2]

	if !IsUserBlacklisted(db, userID) {
		s.ChannelMessageSend(m.ChannelID, "User is not blacklisted.")
		return
	}

	err := db.Delete(&models.BlacklistUser{}, "user_id = ?", userID).Error
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error unblacklisting the user.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "User unblacklisted successfully.")
}
