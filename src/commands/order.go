package commands

import (
	"github.com/bwmarrin/discordgo"
	"sandwich-delivery/src/database"
	"sandwich-delivery/src/models"
)

type OrderCommand struct{}

func (c OrderCommand) getName() string {
	return "order"
}

func (c OrderCommand) getCommandData() *discordgo.ApplicationCommand {
	DMPermission := new(bool)
	*DMPermission = false
	return &discordgo.ApplicationCommand{
		Name:        c.getName(),
		Description: "Order something!",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "order",
				Description: "What do you want to order?",
				Required:    true,
			},
		},
		DMPermission: DMPermission,
	}
}

func (c OrderCommand) registerGuild() string {
	return ""
}

func (c OrderCommand) permissionLevel() models.UserPermissionLevel {
	return models.PermissionLevelUser
}

func (c OrderCommand) execute(session *discordgo.Session, event *discordgo.InteractionCreate) {
	if InteractionIsDM(event) {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				// todo https://github.com/refractored/sandwich-delivery/issues/5
				Content: "This command can only be used in servers!",
			},
		})
	}

	orderOption := event.ApplicationCommandData().Options[0].StringValue()

	var order models.Order
	var user models.User

	resp := database.GetDB().First(&order, "user_id = ? AND delivered = ?", GetUser(event).ID, false)
	if resp.RowsAffected > 0 {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					// todo https://github.com/refractored/sandwich-delivery/issues/5
					{
						Title:       "Error!",
						Description: "You already have a pending order!",
						Color:       0xff2c2c, // Green color
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "Executed by " + DisplayName(event),
							IconURL: GetUser(event).AvatarURL("256"),
						},
						Author: &discordgo.MessageEmbedAuthor{
							Name:    "Sandwich Delivery",
							IconURL: session.State.User.AvatarURL("256"),
						},
					},
				},
			},
		})
		return
	}

	resp = database.GetDB().First(&user, "user_id = ?", GetUser(event).ID)

	if user.Credits < 1 {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					// todo https://github.com/refractored/sandwich-delivery/issues/5
					{
						Title:       "Error!",
						Description: "You do not have enough Credits to order!",
						Color:       0xff2c2c, // Green color
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "Executed by " + DisplayName(event),
							IconURL: GetUser(event).AvatarURL("256"),
						},
						Author: &discordgo.MessageEmbedAuthor{
							Name:    "Sandwich Delivery",
							IconURL: session.State.User.AvatarURL("256"),
						},
					},
				},
			},
		})
		return
	}

	user.Credits = user.Credits - 1
	user.OrdersCreated = user.OrdersCreated + 1
	database.GetDB().Save(&user)

	order = models.Order{
		UserID:           GetUser(event).ID,
		OrderDescription: orderOption,
		Delivered:        false,
		SourceServer:     event.GuildID,
		SourceChannel:    event.GuildID,
	}

	resp = database.GetDB().Save(&order)

	session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "Order Placed!",
					Description: "Thanks for placing your order!" +
						"\nPlease give our staff some time!",
					Color: 0x00ff00, // Green color
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Executed by " + DisplayName(event),
						IconURL: GetUser(event).AvatarURL("256"),
					},
					Author: &discordgo.MessageEmbedAuthor{
						Name:    "Sandwich Delivery",
						IconURL: session.State.User.AvatarURL("256"),
					},
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Your Order:",
							Value:  orderOption,
							Inline: false,
						},
					},
				},
			},
		},
	})
}
