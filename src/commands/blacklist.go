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
	if InteractionIsDM(event) {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				// todo https://github.com/refractored/sandwich-delivery/issues/5
				Content: "This command can only be used in servers!",
			},
		})
	}

	options := event.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options[0].Options))
	for _, opt := range options[0].Options {
		optionMap[opt.Name] = opt
	}
	switch options[0].Name {

	case "add":
		userOption := event.ApplicationCommandData().Options[0].Options[0].UserValue(session)

		var user models.User

		database.GetDB().Find(&user, "user_id = ?", userOption.ID)

		user.IsBlacklisted = true
		database.GetDB().Save(&user)

		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "User blacklisted successfully.",
			},
		})
		break
	case "remove":
		userOption := event.ApplicationCommandData().Options[0].Options[0].UserValue(session)

		var user models.User

		resp := database.GetDB().Find(&user, "user_id = ?", userOption.ID)
		if resp.RowsAffected == 0 || user.IsBlacklisted == false {
			session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					// todo https://github.com/refractored/sandwich-delivery/issues/5
					Content: "User is not blacklisted.",
				},
			})
			return
		}

		user.IsBlacklisted = false
		database.GetDB().Save(&user)

		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "User unblacklisted successfully.",
			},
		})
		break
	}
}
