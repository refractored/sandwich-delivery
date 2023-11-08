package commands

import (
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"os"
)

func ShutdownCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	shutdownMessage := []string{
		"Was it something I did? :( *(Shutting Down)*",
		"Whatever you say... *(Shutting Down)*",
		"Whatever. *(Shutting Down)*",
		"Rude. *(Shutting Down)*",
		"Fine... I guess... :( *(Shutting Down)*"}
	selection := rand.Intn(len(shutdownMessage))
	s.ChannelMessageSend(m.ChannelID, shutdownMessage[selection])
	s.Close()
	os.Exit(0)
}
