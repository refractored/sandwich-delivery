package commands

import "github.com/bwmarrin/discordgo"

func owneronly(s *discordgo.Session, m *discordgo.MessageCreate, silent bool) {
	if !silent {
		s.ChannelMessageSend(m.ChannelID, "This command can only be ran by the bot owner!")
	}
	return
}
func supportonly(s *discordgo.Session, m *discordgo.MessageCreate, silent bool) {
	if silent == false {
		s.ChannelMessageSend(m.ChannelID, "This command can only be executed in the support server!")
	}
	return
}
func lackingpermissions(s *discordgo.Session, m *discordgo.MessageCreate, silent bool) {
	if !silent {
		s.ChannelMessageSend(m.ChannelID, "You do not have the required permissions to use this command.")
	}
	return
}
