package commands

import (
	"github.com/bwmarrin/discordgo"
	"sandwich-delivery/src/config"
	"sandwich-delivery/src/database"
	"sandwich-delivery/src/models"
)

type UnblacklistCommand struct{}

func (c UnblacklistCommand) getName() string {
	return "unblacklist"
}

func (c UnblacklistCommand) getCommandData() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        c.getName(),
		Description: "Unblacklist a user from using the bot.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user to unblacklist.",
				Required:    true,
			},
		},
	}
}

func (c UnblacklistCommand) registerGuild() string {
	return config.GetConfig().GuildID
}

func (c UnblacklistCommand) permissionLevel() models.UserPermissionLevel {
	return models.PermissionLevelAdmin
}

func (c UnblacklistCommand) execute(session *discordgo.Session, event *discordgo.InteractionCreate) {
	userOption := event.ApplicationCommandData().Options[0].UserValue(session)

	var user models.User

	resp := database.GetDB().First(&user, "user_id = ?", userOption.ID)
	if resp.RowsAffected == 0 || user.IsBlacklisted == false {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				// todo https://github.com/refractored/sandwich-delivery/issues/5
				Content: "User is not blacklisted.",
			},
		})
		return
	}

	user.IsBlacklisted = false
	database.GetDB().Save(&user)

	session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "User unblacklisted successfully.",
		},
	})
}
