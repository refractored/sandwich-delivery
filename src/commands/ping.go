package commands

import "github.com/bwmarrin/discordgo"

type PingCommand struct{}

func (p PingCommand) getName() string {
	return "ping"
}

func (p PingCommand) getCommandData() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{Name: p.getName(), Description: "Pong!"}
}

func (p PingCommand) execute(session *discordgo.Session, event *discordgo.InteractionCreate) {
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
