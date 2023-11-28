package commands

import (
	"github.com/bwmarrin/discordgo"
	"os"
	"sandwich-delivery/src/config"
	"sandwich-delivery/src/models"
)

type PurgeCommand struct{}

func (c PurgeCommand) getName() string {
	return "purgecmds"
}

func (c PurgeCommand) getCommandData() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{Name: c.getName(), Description: "Purges all commands and restarts."}
}

func (c PurgeCommand) registerGuild() string {
	return config.GetConfig().GuildID
}

func (c PurgeCommand) permissionLevel() models.UserPermissionLevel {
	return models.PermissionLevelOwner
}

func (c PurgeCommand) execute(session *discordgo.Session, event *discordgo.InteractionCreate) {

	session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Give me a minute...",
		},
	})

	applicationCommands, err := session.ApplicationCommands(session.State.User.ID, "")
	if err != nil {
		return
	}
	for _, applicationCommand := range applicationCommands {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})
		session.ApplicationCommandDelete(session.State.User.ID, "", applicationCommand.ID)
	}

	if InteractionIsDM(event) {
		applicationCommands, err := session.ApplicationCommands(session.State.User.ID, event.GuildID)
		if err != nil {
			return
		}
		for _, applicationCommand := range applicationCommands {
			session.ApplicationCommandDelete(session.State.User.ID, event.GuildID, applicationCommand.ID)
		}
	}

	session.Close()
	os.Exit(0)
}
