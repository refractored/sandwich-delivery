package commands

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

var commands = map[string]Command{
	CoinflipCommand{}.getName():           CoinflipCommand{},
	OrderCommand{}.getName():              OrderCommand{},
	DelOrderCommand{}.getName():           DelOrderCommand{},
	PingCommand{}.getName():               PingCommand{},
	ShutdownCommand{}.getName():           ShutdownCommand{},
	BlacklistCommand{}.getName():          BlacklistCommand{},
	UnblacklistCommand{}.getName():        UnblacklistCommand{},
	SetPermissionLevelCommand{}.getName(): SetPermissionLevelCommand{},
}

func RegisterCommands(session *discordgo.Session) {
	for n, d := range commands {
		_, err := session.ApplicationCommandCreate(session.State.User.ID, d.registerGuild(), d.getCommandData())
		if err != nil {
			log.Fatalf("Unable to register command %s: %v\n", n, err)
		}

		log.Printf("Registered command %s\n", n)
	}
}

func HandleCommand(session *discordgo.Session, event *discordgo.InteractionCreate) {
	if event.Type != discordgo.InteractionApplicationCommand {
		return
	}

	if IsUserBlacklisted(GetUser(event).ID) {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				// todo https://github.com/refractored/sandwich-delivery/issues/5
				Content: "You are blacklisted from using this bot.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	command := commands[event.ApplicationCommandData().Name]
	if command == nil {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				// todo https://github.com/refractored/sandwich-delivery/issues/5
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
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})

		err := session.ApplicationCommandDelete(session.State.User.ID, "", event.ApplicationCommandData().ID)
		if err != nil {
			log.Printf("Unable to delete global command %s: %v\n", event.ApplicationCommandData().Name, err)
		}

		err = session.ApplicationCommandDelete(session.State.User.ID, event.GuildID, event.ApplicationCommandData().ID)
		if err != nil {
			log.Printf("Unable to delete guild-specific command %s: %v\n", event.ApplicationCommandData().Name, err)
		}

		log.Printf("Deleted command %s\n", event.ApplicationCommandData().Name)
		return
	}

	if command.permissionLevel() > GetPermissionLevel(GetUser(event).ID) {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				// todo https://github.com/refractored/sandwich-delivery/issues/5
				Content: "You do not have permission to use this command.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	command.execute(session, event)
}
