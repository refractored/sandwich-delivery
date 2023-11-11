package commands

import (
	"github.com/bwmarrin/discordgo"
	"math/rand"
)

func CoinflipCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	var displayname = DisplayName(s, m)

	coin := []string{"Heads", "Tails"}
	selection := rand.Intn(len(coin))

	pingEmbed := &discordgo.MessageEmbed{
		Title:       "Flipping a coin...",
		Description: coin[selection],
		Color:       0x00ff00, // Green color
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Executed by " + displayname,
			IconURL: m.Author.AvatarURL("256"),
		},
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "Sandwich Delivery",
			IconURL: s.State.User.AvatarURL("256"),
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, pingEmbed)

}
