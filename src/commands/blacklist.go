package commands

import (
	"github.com/bwmarrin/discordgo"
	"sandwich-delivery/src/config"
	"sandwich-delivery/src/database"
	"sandwich-delivery/src/models"
)

type BlacklistCommand struct{}

func (c BlacklistCommand) getName() string {
	return "blacklist"
}

func (c BlacklistCommand) getCommandData() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        c.getName(),
		Description: "Blacklist a user from using the bot.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user to blacklist.",
				Required:    true,
			},
		},
	}
}

func (c BlacklistCommand) registerGuild() string {
	return config.GetConfig().GuildID
}

func (c BlacklistCommand) execute(session *discordgo.Session, event *discordgo.InteractionCreate) {
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

	if IsUserBlacklisted(user.ID) {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				// todo https://github.com/refractored/sandwich-delivery/issues/5
				Content: "User is already blacklisted.",
			},
		})
		return
	}

	resp := database.GetDB().Create(&models.BlacklistUser{UserID: user.ID})
	if resp.Error != nil {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				// todo https://github.com/refractored/sandwich-delivery/issues/5
				Content: "Error blacklisting the user.",
			},
		})
		return
	}

	session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "User blacklisted successfully.",
		},
	})
}
