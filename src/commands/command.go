package commands

import (
	"github.com/bwmarrin/discordgo"
	"sandwich-delivery/src/config"
	"sandwich-delivery/src/database"
	"sandwich-delivery/src/models"
	"slices"
)

type Command interface {
	// Returns the name of the command.
	getName() string

	// Returns the command data
	getCommandData() *discordgo.ApplicationCommand

	// Returns the guild ID to register the command in, empty string if the command is global.
	registerGuild() string

	// Returns the required permission level to execute the command.
	permissionLevel() models.UserPermissionLevel

	// Executes the command
	execute(session *discordgo.Session, event *discordgo.InteractionCreate)
}

func DisplayName(event *discordgo.InteractionCreate) string {
	user := GetUser(event)

	if user.Discriminator == "0" {
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

func DoesUserExist(userID string) bool {
	var user models.User

	resp := database.GetDB().Find(&user, "user_id = ?", userID)

	return resp.RowsAffected != 0
}

func IsUserBlacklisted(userID string) bool {
	if GetPermissionLevel(userID) == models.PermissionLevelOwner {
		return false
	}

	var user models.User

	resp := database.GetDB().Find(&user, "user_id = ?", userID)

	if resp.RowsAffected == 0 {
		return false
	}

	return user.IsBlacklisted
}

// Returns true if the user is the owner of the bot (defined in the config).
func IsHardcodedOwner(userID string) bool {
	return slices.Contains(config.GetConfig().Owners, userID)
}

func GetPermissionLevel(userID string) models.UserPermissionLevel {
	if IsHardcodedOwner(userID) {
		return models.PermissionLevelOwner
	}

	var user models.User

	db := database.GetDB()

	resp := db.First(&user, "user_id = ?", userID)

	if resp.RowsAffected == 0 {
		return models.PermissionLevelUser
	}

	return user.PermissionLevel
}
