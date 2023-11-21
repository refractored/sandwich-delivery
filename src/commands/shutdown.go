package commands

import (
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"os"
	"sandwich-delivery/src/config"
)

type ShutdownCommand struct{}

func (c ShutdownCommand) getName() string {
	return "shutdown"
}

func (c ShutdownCommand) getCommandData() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{Name: c.getName(), Description: "Shuts down the bot."}
}

func (c ShutdownCommand) registerGuild() string {
	return config.GetConfig().GuildID
}

func (c ShutdownCommand) execute(session *discordgo.Session, event *discordgo.InteractionCreate) {
	var shutdownMessages = []string{
		"Was it something I did? :( *(Shutting Down)*",
		"Whatever you say... *(Shutting Down)*",
		"Whatever. *(Shutting Down)*",
		"Rude. *(Shutting Down)*",
		"Fine... I guess... :( *(Shutting Down)*",
	}

	selection := rand.Intn(len(shutdownMessages))

	if !IsOwner(GetUser(event).ID) {
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
