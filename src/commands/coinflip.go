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
	return &discordgo.ApplicationCommand{Name: c.getName(), Description: "Flip a virtual coin."}
}

func (c CoinflipCommand) registerGuild() string {
	return ""
}

func (c CoinflipCommand) execute(session *discordgo.Session, event *discordgo.InteractionCreate) {
	coin := []string{"Heads", "Tails"}
	selection := rand.Intn(len(coin))

	session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Flipping a coin...",
					Description: coin[selection],
					Color:       0x00ff00, // Green color
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Executed by " + DisplayName(event),
						IconURL: GetUser(event).AvatarURL("256"),
					},
					Author: &discordgo.MessageEmbedAuthor{
						Name:    "Sandwich Delivery",
						IconURL: GetUser(event).AvatarURL("256"),
					},
				},
			},
		},
	})
}
