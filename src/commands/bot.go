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
	}
}

func (c BotCommand) registerGuild() string {
	return config.GetConfig().GuildID
}

func (c BotCommand) permissionLevel() models.UserPermissionLevel {
	return models.PermissionLevelAdmin
}

func (c BotCommand) execute(session *discordgo.Session, event *discordgo.InteractionCreate) {

	options := event.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options[0].Options))
	for _, opt := range options[0].Options {
		optionMap[opt.Name] = opt
	}
	switch options[0].Name {
	case "shutdown":
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
		break
	case "purgecmds":
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Give me a minute...",
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
		break
	}
}
