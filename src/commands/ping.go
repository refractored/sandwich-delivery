package commands

import (
	"github.com/bwmarrin/discordgo"
	"sandwich-delivery/src/models"
)

type PingCommand struct{}

func (c PingCommand) getName() string {
	return "ping"
}

func (c PingCommand) getCommandData() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{Name: c.getName(), Description: "Pong!"}
}

func (c PingCommand) DMsAllowed() bool {
	return true
}

func (c PingCommand) registerGuild() string {
	return ""
}

func (c PingCommand) permissionLevel() models.UserPermissionLevel {
	return models.PermissionLevelUser
}

func (c PingCommand) execute(session *discordgo.Session, event *discordgo.InteractionCreate) {
	session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Pong!",
					Description: "Bot Ping: " + session.HeartbeatLatency().String(),
					Color:       0x00ff00, // Green color
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Executed by " + DisplayName(event),
						IconURL: GetUser(event).AvatarURL("256"),
					},
					Author: &discordgo.MessageEmbedAuthor{
						Name:    "Sandwich Delivery",
						IconURL: session.State.User.AvatarURL("256"),
					},
				},
			},
		},
	})
}
