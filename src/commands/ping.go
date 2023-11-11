package commands

import (
	"github.com/bwmarrin/discordgo"
)

func PingCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	ping := s.HeartbeatLatency()
	var displayName = DisplayName(s, m)

	//s.ChannelMessageSend(m.ChannelID, "Bot Ping: "+ping.String())
	pingEmbed := &discordgo.MessageEmbed{
		Title:       "Pong!",
		Description: "Bot Ping: " + ping.String(),
		Color:       0x00ff00, // Green color
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Executed by " + displayName,
			IconURL: m.Author.AvatarURL("256"),
		},
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "Sandwich Delivery",
			IconURL: s.State.User.AvatarURL("256"),
		},
	}
	s.ChannelMessageSendEmbed(m.ChannelID, pingEmbed)
}
