package commands

import (
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"os"
	"sandwich-delivery/src/config"
	"sandwich-delivery/src/models"
)

type BotCommand struct{}

func (c BotCommand) getName() string {
	return "bot"
}

func (c BotCommand) getCommandData() *discordgo.ApplicationCommand {
	DMPermission := new(bool)
	*DMPermission = false
	return &discordgo.ApplicationCommand{Name: c.getName(),
		Description: "Manage the bot.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "shutdown",
				Description: "Shutdown the bot.",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
			{
				Name:        "purgecmds",
				Description: "Purges all commands and restarts.",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
		},
		DMPermission: DMPermission,
	}
}

func (c BotCommand) registerGuild() string {
	return config.GetConfig().GuildID
}

func (c BotCommand) permissionLevel() models.UserPermissionLevel {
	return models.PermissionLevelAdmin
}

func (c BotCommand) execute(session *discordgo.Session, event *discordgo.InteractionCreate) {
	switch event.ApplicationCommandData().Options[0].Name {
	case "shutdown":
		BotShutdown(session, event)
		break
	case "purgecmds":
		BotPurgeCMDS(session, event)
		break
	}
}
func BotShutdown(session *discordgo.Session, event *discordgo.InteractionCreate) {
	var shutdownMessages = []string{
		"Was it something I did? :( *(Shutting Down)*",
		"Whatever you say... *(Shutting Down)*",
		"Whatever. *(Shutting Down)*",
		"Rude. *(Shutting Down)*",
		"Fine... I guess... :( *(Shutting Down)*",
	}

	selection := rand.Intn(len(shutdownMessages))

	session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: shutdownMessages[selection],
		},
	})

	session.Close()
	os.Exit(0)
}

func BotPurgeCMDS(session *discordgo.Session, event *discordgo.InteractionCreate) {
	session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Give me a minute...",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	applicationCommands, err := session.ApplicationCommands(session.State.User.ID, "")
	if err != nil {
		return
	}
	for _, applicationCommand := range applicationCommands {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})
		session.ApplicationCommandDelete(session.State.User.ID, "", applicationCommand.ID)
	}

	if InteractionIsDM(event) {
		applicationCommands, err := session.ApplicationCommands(session.State.User.ID, event.GuildID)
		if err != nil {
			return
		}
		for _, applicationCommand := range applicationCommands {
			session.ApplicationCommandDelete(session.State.User.ID, event.GuildID, applicationCommand.ID)
		}
	}

	session.Close()
	os.Exit(0)
}
