package commands

import (
	"github.com/bwmarrin/discordgo"
	"math/rand"
)

type CoinflipCommand struct{}

func (c CoinflipCommand) getName() string {
	return "coinflip"
}

func (c CoinflipCommand) getCommandData() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{Name: c.getName(), Description: "idk"}
}

func (c CoinflipCommand) execute(session *discordgo.Session, event *discordgo.InteractionCreate) {
	displayName := DisplayName(session, event)

	coin := []string{"Heads", "Tails"}
	selection := rand.Intn(len(coin))

	pingEmbed := &discordgo.MessageEmbed{
		Title:       "Flipping a coin...",
		Description: coin[selection],
		Color:       0x00ff00, // Green color
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Executed by " + displayName,
			IconURL: event.User.AvatarURL("256"),
		},
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "Sandwich Delivery",
			IconURL: session.State.User.AvatarURL("256"),
		},
	}

	session.ChannelMessageSendEmbed(event.ChannelID, pingEmbed)
}
