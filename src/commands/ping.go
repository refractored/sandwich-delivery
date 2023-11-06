package commands

import (
	"github.com/bwmarrin/discordgo"
)

func PingCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	ping := s.HeartbeatLatency()

	s.ChannelMessageSend(m.ChannelID, "Bot Ping: "+ping.String())
}
