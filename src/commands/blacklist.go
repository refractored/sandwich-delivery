package commands

import (
	"github.com/bwmarrin/discordgo"
	"sandwich-delivery/src/config"
	"sandwich-delivery/src/database"
	"sandwich-delivery/src/models"
)

type BlacklistCommand struct{}

func (c BlacklistCommand) getName() string {
	return "blacklist"
}

func (c BlacklistCommand) getCommandData() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        c.getName(),
		Description: "Blacklist a user from using the bot.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "add",
				Description: "Add an user to the blacklisted users.",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionUser,
						Name:        "user",
						Description: "The user to blacklist.",
						Required:    true,
					},
				},
			},
			{
				Name:        "remove",
				Description: "remove an user from the blacklisted users.",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionUser,
						Name:        "user",
						Description: "The user to remove.",
						Required:    true,
					},
				},
			},
		},
	}
}

func (c BlacklistCommand) registerGuild() string {
	return config.GetConfig().GuildID
}

func (c BlacklistCommand) permissionLevel() models.UserPermissionLevel {
	return models.PermissionLevelAdmin
}

func (c BlacklistCommand) execute(session *discordgo.Session, event *discordgo.InteractionCreate) {
	switch event.ApplicationCommandData().Options[0].Name {
	case "add":
		BlacklistAdd(session, event)
		break
	case "remove":
		BlacklistRemove(session, event)
		break
	}
}

func BlacklistAdd(session *discordgo.Session, event *discordgo.InteractionCreate) {
	userOption := event.ApplicationCommandData().Options[0].Options[0].UserValue(session)

	var user models.User

	resp := database.GetDB().Find(&user, "user_id = ?", userOption.ID)
	if resp.RowsAffected == 0 {
		user.UserID = userOption.ID
	}
	user.IsBlacklisted = true
	database.GetDB().Save(&user)

	session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "User blacklisted successfully.",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func BlacklistRemove(session *discordgo.Session, event *discordgo.InteractionCreate) {
	userOption := event.ApplicationCommandData().Options[0].Options[0].UserValue(session)

	var user models.User

	resp := database.GetDB().Find(&user, "user_id = ?", userOption.ID)
	if resp.RowsAffected == 0 {
		user.UserID = userOption.ID
	}
	user.IsBlacklisted = true
	database.GetDB().Save(&user)
	session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "User blacklisted successfully.",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
