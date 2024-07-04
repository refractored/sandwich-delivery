package commands

import (
	"github.com/bwmarrin/discordgo"
	"sandwich-delivery/src/config"
	"sandwich-delivery/src/database"
	"sandwich-delivery/src/models"
	"strconv"
)

type UserManageCommand struct{}

func (c UserManageCommand) getName() string {
	return "usermanage"
}

func (c UserManageCommand) getCommandData() *discordgo.ApplicationCommand {
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
						Description: "Set the user's permission level.",
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
	return config.GetConfig().GuildID
}

func (c UserManageCommand) permissionLevel() models.UserPermissionLevel {
	return models.PermissionLevelAdmin
}

func (c UserManageCommand) execute(session *discordgo.Session, event *discordgo.InteractionCreate) {
	switch event.ApplicationCommandData().Options[0].Name {
	case "resetdaily":
		UserManageResetDaily(session, event)
		break
	case "modify":
		UserManageModify(session, event)
		break
	case "purge":
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "purge",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		break
	case "addcredits":
		UserManageAddCredits(session, event)
		break
	case "takecredits":
		UserManageTakeCredits(session, event)
		break
	}
}

func UserManageResetDaily(session *discordgo.Session, event *discordgo.InteractionCreate) {
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
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func UserManageModify(session *discordgo.Session, event *discordgo.InteractionCreate) {
	if GetPermissionLevel(GetUser(event).ID) < GetPermissionLevel(event.ApplicationCommandData().Options[0].Options[0].UserValue(session).ID) {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You cannot modify permissions of users with higher permissions than you!",
				Flags:   discordgo.MessageFlagsEphemeral,
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
					Flags:   discordgo.MessageFlagsEphemeral,
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
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}
	database.GetDB().Save(&user)
	session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Modified Data of " + event.ApplicationCommandData().Options[0].Options[0].UserValue(session).Username,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	return
}
func UserManageAddCredits(session *discordgo.Session, event *discordgo.InteractionCreate) {
	var user models.User

	options := event.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options[0].Options))
	for _, opt := range options[0].Options {
		optionMap[opt.Name] = opt
	}

	var userid string
	if option, ok := optionMap["user"]; ok {
		userid = option.UserValue(nil).ID
	}
	resp := database.GetDB().Find(&user, "user_id = ?", userid)
	if resp.RowsAffected == 0 {
		user.UserID = userid
		database.GetDB().Save(&user)
	}
	if option, ok := optionMap["credits"]; ok {
		user.Credits += uint32(option.IntValue())
		if user.Credits < 0 {
			session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "User cannot have negative credits!",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}
		database.GetDB().Save(&user)
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Added " + strconv.Itoa(int(option.IntValue())) + " credits to " + event.ApplicationCommandData().Options[0].Options[0].UserValue(session).Username + "\n" +
					"User now has " + strconv.Itoa(int(user.Credits)) + " credits.",
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
	}
}
func UserManageTakeCredits(session *discordgo.Session, event *discordgo.InteractionCreate) {
	var user models.User

	options := event.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options[0].Options))
	for _, opt := range options[0].Options {
		optionMap[opt.Name] = opt
	}

	var userid string
	if option, ok := optionMap["user"]; ok {
		userid = option.UserValue(nil).ID
	}
	resp := database.GetDB().Find(&user, "user_id = ?", userid)
	if resp.RowsAffected == 0 {
		user.UserID = userid
		database.GetDB().Save(&user)
	}
	if option, ok := optionMap["credits"]; ok {
		if user.Credits < uint32(option.IntValue()) {
			session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "User cannot have negative credits!",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}
		user.Credits -= uint32(option.IntValue())
		database.GetDB().Save(&user)
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Took " + strconv.Itoa(int(option.IntValue())) + " credits from " + event.ApplicationCommandData().Options[0].Options[0].UserValue(session).Username + "\n" +
					"User now has " + strconv.Itoa(int(user.Credits)) + " credits.",
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
	}
}
