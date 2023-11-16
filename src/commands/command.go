package commands

import (
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
	"sandwich-delivery/src/models"
)

type Command interface {
	getName() string
	getCommandData() *discordgo.ApplicationCommand

	execute(session *discordgo.Session, event *discordgo.InteractionCreate)
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

func IsUserBlacklisted(db *gorm.DB, userID string) bool {
	var user models.BlacklistUser

	result := db.Select("user_id").Where("user_id = ?", userID).First(&user)

	return result.Error == nil
}
