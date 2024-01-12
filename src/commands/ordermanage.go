package commands

import (
	"github.com/bwmarrin/discordgo"
	"sandwich-delivery/src/config"
	"sandwich-delivery/src/database"
	"sandwich-delivery/src/models"
)

type OrderManageCommand struct{}

func (c OrderManageCommand) getName() string {
	return "usermanage"
}

func (c OrderManageCommand) getCommandData() *discordgo.ApplicationCommand {
	var creditMin float64 = 0
	return &discordgo.ApplicationCommand{Name: c.getName(),
		Description: "Manage data of an user.",
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
				Name:        "addcredits",
				Description: "Remove credits from an user.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionUser,
						Name:        "user",
						Description: "The user to add credits to",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "credits",
						Description: "The amount to add.",
						MinValue:    &creditMin,
						MaxValue:    65535,
						Required:    true,
					},
				},
				Type: discordgo.ApplicationCommandOptionSubCommand,
			},
			{
				Name:        "takecredits",
				Description: "Remove Credits from an user.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionUser,
						Name:        "user",
						Description: "The user to take credits of",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "credits",
						Description: "The amount to add.",
						MinValue:    &creditMin,
						MaxValue:    65535,
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
				Description: "Modify data of an order. Use with caution.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "id",
						Description: "The ID of the order to accept.",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "UserID",
						Description: "Changes the order's creator. Must be a valid user ID.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "OrderDescription",
						Description: "Changes the order's description.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "SourceServer",
						Description: "Changes the Source of the order. Server must contain SourceChannel.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "SourceChannel",
						Description: "Changes the Source channel. Channel must be from the SourceServer.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionBoolean,
						Name:        "Assignee",
						Description: "Changes who is assigned to the order.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionBoolean,
						Name:        "Tipped",
						Description: "Changes the tipped status. Does not actually tip the user.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "Status",
						Description: "Changes the order status. Does not execute the code that related to the status.",
						Required:    false,
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{
								Name:  "Waiting",
								Value: models.StatusWaiting,
							},
							{
								Name:  "Accepted",
								Value: models.StatusAccepted,
							},
							{
								Name:  "In Transit",
								Value: models.StatusInTransit,
							},
							{
								Name:  "Delivered",
								Value: models.StatusDelivered,
							},
							{
								Name:  "Cancelled",
								Value: models.StatusCancelled,
							},
							{
								Name:  "Moderated",
								Value: models.StatusModerated,
							},
							{
								Name:  "Error",
								Value: models.StatusError,
							},
						},
					},
				},
				Type: discordgo.ApplicationCommandOptionSubCommand,
			},
		},
	}
}

func (c OrderManageCommand) registerGuild() string {
	return config.GetConfig().GuildID
}

func (c OrderManageCommand) permissionLevel() models.UserPermissionLevel {
	return models.PermissionLevelAdmin
}

func (c OrderManageCommand) execute(session *discordgo.Session, event *discordgo.InteractionCreate) {
	switch event.ApplicationCommandData().Options[0].Name {
	case "delete":
		OrderManageResetDaily(session, event)
		break
	case "modify":
		OrderManageModify(session, event)
		break
	}
}

func OrderManageResetDaily(session *discordgo.Session, event *discordgo.InteractionCreate) {
	var user models.User

	resp := database.GetDB().Find(&user, "user_id = ?", event.ApplicationCommandData().Options[0].Options[0].UserValue(session).ID)
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
}

func OrderManageModify(session *discordgo.Session, event *discordgo.InteractionCreate) {
	if GetPermissionLevel(GetUser(event).ID) < GetPermissionLevel(event.ApplicationCommandData().Options[0].Options[0].UserValue(session).ID) {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You cannot modify permissions of users with higher permissions than you!",
			},
		})
		return
	}

	var user models.User

	options := event.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options[0].Options))
	for _, opt := range options[0].Options {
		optionMap[opt.Name] = opt
	}

	resp := database.GetDB().Find(&user, "user_id = ?", event.ApplicationCommandData().Options[0].Options[0].UserValue(session).ID)

	if resp.RowsAffected == 0 {
		user.UserID = event.ApplicationCommandData().Options[0].Options[0].UserValue(session).ID
		database.GetDB().Save(&user)
	}
	if option, ok := optionMap["blacklisted"]; ok {
		user.IsBlacklisted = option.BoolValue()
	}
	if option, ok := optionMap["credits"]; ok {
		user.Credits = uint32(option.IntValue())
	}
	if option, ok := optionMap["permissionlevel"]; ok {
		user.PermissionLevel = models.UserPermissionLevel(option.IntValue())
		if GetPermissionLevel(GetUser(event).ID) < models.UserPermissionLevel(option.IntValue()) {
			session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You cannot assign permissions higher than your own!",
				},
			})
			return
		}
	}
	if len(optionMap) == 0 {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No options provided.",
			},
		})
		return
	}
	database.GetDB().Save(&user)
	session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Modified Data of " + event.ApplicationCommandData().Options[0].Options[0].UserValue(session).Username,
		},
	})
	return
}
