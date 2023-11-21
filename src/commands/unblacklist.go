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

func (c UnblacklistCommand) execute(session *discordgo.Session, event *discordgo.InteractionCreate) {
	if !IsOwner(GetUser(event).ID) {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				// todo https://github.com/refractored/sandwich-delivery/issues/5
				Content: "You are not the bot owner!",
			},
		})
		return
	}

	user := event.ApplicationCommandData().Options[0].UserValue(session)

	if !IsUserBlacklisted(user.ID) {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				// todo https://github.com/refractored/sandwich-delivery/issues/5
				Content: "User is not blacklisted.",
			},
		})
		return
	}

	resp := database.GetDB().Delete(&models.BlacklistUser{}, "user_id = ?", user.ID)

	if resp.Error != nil {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				// todo https://github.com/refractored/sandwich-delivery/issues/5
				Content: "Error unblacklisting the user.",
			},
		})
		return
	}

	session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "User unblacklisted successfully.",
		},
	})
}
