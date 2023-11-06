package commands

import (
	"github.com/bwmarrin/discordgo"
	"math/rand"
)

func CoinflipCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	coin := []string{"Heads", "Tails"}
	selection := rand.Intn(len(coin))
	s.ChannelMessageSend(m.ChannelID, coin[selection])
}
dsadad