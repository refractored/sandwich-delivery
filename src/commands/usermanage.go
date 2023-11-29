package commands

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"sandwich-delivery/src/database"
	"sandwich-delivery/src/models"
)

type UserManageCommand struct{}

func (c UserManageCommand) getName() string {
	return "usermanage"
}

func (c UserManageCommand) getCommandData() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{Name: c.getName(),
		Description: "Flip a virtual coin.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "resetdaily",
				Description: "Reset the daily timer of a user.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionUser,
						Name:        "user",
						Description: "The user to reset the daily timer.",
						Required:    true,
					},
				},
				Type: discordgo.ApplicationCommandOptionSubCommand,
			},
			{
				Name:        "view",
				Description: "View the information of a user.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionUser,
						Name:        "user",
						Description: "The user to view information of.",
						Required:    true,
					},
				},
				Type: discordgo.ApplicationCommandOptionSubCommand,
			},
			{
				Name:        "purge",
				Description: "Purges an user's orders",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionUser,
						Name:        "user",
						Description: "The user to set the purge command data of.",
						Required:    true,
					},
				},
				Type: discordgo.ApplicationCommandOptionSubCommand,
			},
			{
				Name:        "modify",
				Description: "Modify data of an user.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionUser,
						Name:        "user",
						Description: "The user to set the purge command data of.",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "credits",
						Description: "Sets user's credits.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionBoolean,
						Name:        "blacklisted",
						Description: "Change a user's blacklist status.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "permissionlevel",
						Description: "The user to set the purge command data of.",
						Required:    false,
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{
								Name:  "Default",
								Value: models.PermissionLevelUser,
							},
							{
								Name:  "Mod",
								Value: models.PermissionLevelMod,
							},
							{
								Name:  "Artist",
								Value: models.PermissionLevelArtist,
							},
							{
								Name:  "Admin",
								Value: models.PermissionLevelAdmin,
							},
							{
								Name:  "Owner",
								Value: models.PermissionLevelOwner,
							},
						},
					},
				},
				Type: discordgo.ApplicationCommandOptionSubCommand,
			},
		},
	}
}

func (c UserManageCommand) registerGuild() string {
	return ""
}

func (c UserManageCommand) permissionLevel() models.UserPermissionLevel {
	return models.PermissionLevelAdmin
}

func (c UserManageCommand) execute(session *discordgo.Session, event *discordgo.InteractionCreate) {
	options := event.ApplicationCommandData().Options
	var user models.User

	log.Println(event.ApplicationCommandData().Options[0].Options[0].UserValue(session).ID)

	switch options[0].Name {
	case "resetdaily":
		resp := database.GetDB().First(&user, "user_id = ?", event.ApplicationCommandData().Options[0].Options[0].UserValue(session).ID)
		if resp.RowsAffected == 0 {
			user.UserID = event.ApplicationCommandData().Options[0].Options[0].UserValue(session).ID

			database.GetDB().Save(&user)
		}
		user.DailyClaimedAt = nil
		database.GetDB().Save(&user)
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Reset daily of " + event.ApplicationCommandData().Options[0].Options[0].UserValue(session).Username,
			},
		})
	case "modify":
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "modify",
			},
		})
	case "purge":
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "purge",
			},
		})
	case "view":
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "view",
			},
		})
		//switch options[0].Name {
		//case "nested-subcommand":
		//	content = "Nice, now you know how to execute nested commands too"
		//default:
		//	content = "Oops, something went wrong.\n" +
		//		"Hol' up, you aren't supposed to see this message."
		//}
	}
}
