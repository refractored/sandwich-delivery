package commands

import (
	"github.com/bwmarrin/discordgo"
	"sandwich-delivery/src/config"
	"sandwich-delivery/src/database"
	"sandwich-delivery/src/models"
)

type SetPermissionLevelCommand struct{}

func (c SetPermissionLevelCommand) getName() string {
	return "setpermissionlevel"
}

func (c SetPermissionLevelCommand) getCommandData() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        c.getName(),
		Description: "Sets the permission level of a user.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user to set the permission level of.",
				Required:    true,

			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "level",
				Description: "The permission level to set the user to.",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "User",
						Value: models.PermissionLevelUser,
					},
					{
						Name:  "Mod",
						Value: models.PermissionLevelMod,
					},
					{
						Name:  "Artist",
						Value: models.PermissionLevelArtist,
					},
					{
						Name:  "Admin",
						Value: models.PermissionLevelAdmin,
					},
					{
						Name:  "Owner",
						Value: models.PermissionLevelOwner,
					},
				},
			},
		},
	}
}

func (c SetPermissionLevelCommand) registerGuild() string {
	return config.GetConfig().GuildID
}

func (c SetPermissionLevelCommand) permissionLevel() models.UserPermissionLevel {
	return models.PermissionLevelOwner
}

func (c SetPermissionLevelCommand) execute(session *discordgo.Session, event *discordgo.InteractionCreate) {
	user := event.ApplicationCommandData().Options[0].UserValue(session)
	level := models.UserPermissionLevel(event.ApplicationCommandData().Options[1].IntValue())

	var userRecord models.User

	resp := database.GetDB().First(&userRecord, "user_id = ?", user.ID)

	if resp.RowsAffected == 0 {
		userRecord = models.User{
			UserID:          user.ID,
			PermissionLevel: level,
		}
	} else {
		userRecord.PermissionLevel = level
	}

	resp = database.GetDB().Save(&userRecord)

	session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Permission level set!",
		},
	})
}
