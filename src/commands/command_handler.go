package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
)

var commands = map[string]Command{
	CoinflipCommand{}.getName(): CoinflipCommand{},
	OrderCommand{}.getName():    OrderCommand{},
}

func RegisterCommands(session *discordgo.Session) {
	for n, d := range commands {
		_, err := session.ApplicationCommandCreate(session.State.User.ID, "823334764307415122", d.getCommandData())
		if err != nil {
			fmt.Printf("Unable to register command %s: %v\n", n, err)
			os.Exit(1)
		}

		fmt.Printf("Registered command %s\n", n)
	}
}

func HandleCommand(session *discordgo.Session, event *discordgo.InteractionCreate) {
	if event.Type != discordgo.InteractionApplicationCommand {
		return
	}

	command := commands[event.ApplicationCommandData().Name]
	if command != nil {
		command.execute(session, event)
		return
	}

	session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Error",
					Description: "Unable to find the slash command: " + event.ApplicationCommandData().Name,
					Color:       0xff0000,
					Author: &discordgo.MessageEmbedAuthor{
						Name:    "Sandwich Delivery",
						IconURL: session.State.User.AvatarURL("256"),
					},
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Executed by " + DisplayName(event),
						IconURL: GetUser(event).AvatarURL("256"),
					},
				},
			},
		},
	})

	err := session.ApplicationCommandDelete(session.State.User.ID, "", event.ApplicationCommandData().Name)
	if err != nil {
		return
	}
}
