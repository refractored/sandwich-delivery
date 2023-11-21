package commands

import (
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
	"sandwich-delivery/src/config"
	"sandwich-delivery/src/models"
)

type Command interface {
	getName() string
	getCommandData() *discordgo.ApplicationCommand

	execute(session *discordgo.Session, event *discordgo.InteractionCreate)
}

func DisplayName(event *discordgo.InteractionCreate) string {
	user := GetUser(event)

	if user.Discriminator != "0" {
		return user.Username
	} else {
		return user.Username + "#" + user.Discriminator
	}
}

func GetUser(event *discordgo.InteractionCreate) *discordgo.User {
	if InteractionIsDM(event) {
		return event.User
	} else {
		return event.Member.User
	}
}

func InteractionIsDM(event *discordgo.InteractionCreate) bool {
	return event.GuildID == ""
}

func IsUserBlacklisted(db *gorm.DB, userID string) bool {
	var user models.BlacklistUser

	result := db.First(&user, "user_id = ?", userID)

	return result.Error == nil
}

func IsOwner(event *discordgo.InteractionCreate) bool {
	return event.Member.User.ID == config.GetConfig().OwnerID
}
