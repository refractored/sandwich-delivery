package commands

import (
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"os"
)

type ShutdownCommand struct{}

func (s ShutdownCommand) getName() string {
	return "shutdown"
}

func (s ShutdownCommand) getCommandData() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{Name: s.getName(), Description: "Shuts down the bot."}
}

func (s ShutdownCommand) execute(session *discordgo.Session, event *discordgo.InteractionCreate) {
	var shutdownMessages = []string{
		"Was it something I did? :( *(Shutting Down)*",
		"Whatever you say... *(Shutting Down)*",
		"Whatever. *(Shutting Down)*",
		"Rude. *(Shutting Down)*",
		"Fine... I guess... :( *(Shutting Down)*",
	}

	selection := rand.Intn(len(shutdownMessages))

	if !IsOwner(event) {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				// todo https://github.com/refractored/sandwich-delivery/issues/5
				Content: "You are not the bot owner!",
			},
		})
		return
	}
	session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: shutdownMessages[selection],
		},
	})
	session.Close()
	os.Exit(0)
}
