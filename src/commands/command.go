package commands

import (
	"github.com/bwmarrin/discordgo"
	"sandwich-delivery/src/config"
	"sandwich-delivery/src/database"
	"sandwich-delivery/src/models"
	"slices"
)

type Command interface {
	/**
	 * Returns the name of the command.
	 */
	getName() string

	/**
	 * Returns the command data for this command.
	 */
	getCommandData() *discordgo.ApplicationCommand

	/**
	 * Returns the guild ID that this command should be registered to, or an empty string if it should be registered globally.
	 */
	registerGuild() string

	/**
	 * Executes the command.
	 */
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

func IsUserBlacklisted(userID string) bool {
	if IsOwner(userID) {
		return false
	}

	var user models.BlacklistUser

	db := database.GetDB()

	result := db.First(&user, "user_id = ?", userID)

	return result.RowsAffected > 0
}

func IsOwner(userID string) bool {
	return slices.Contains(config.GetConfig().Owners, userID)
}
