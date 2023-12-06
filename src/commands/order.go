package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"sandwich-delivery/src/config"
	"sandwich-delivery/src/database"
	"sandwich-delivery/src/models"
	"time"
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
		Description: "Manage your own order.",
		Options: []*discordgo.ApplicationCommandOption{

			{
				Name:        "create",
				Description: "Place an order",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "order",
						Description: "What do you want to order?",
						Required:    true,
					},
				},
			},
			{
				Name:        "cancel",
				Description: "Cancel your order",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
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

	options := event.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options[0].Options))
	for _, opt := range options[0].Options {
		optionMap[opt.Name] = opt
	}
	switch options[0].Name {

	case "create":
		orderOption := event.ApplicationCommandData().Options[0].Options[0].StringValue()
		var order models.Order
		var user models.User

		resp := database.GetDB().Find(&order, "user_id = ? AND status < ?", GetUser(event).ID, models.StatusDelivered)

		if resp.RowsAffected != 0 {
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

		if user.Credits < *config.GetConfig().TokensPerOrder {
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

		user.Credits = user.Credits - *config.GetConfig().TokensPerOrder
		user.OrdersCreated = user.OrdersCreated + 1
		database.GetDB().Save(&user)

		order = models.Order{
			UserID:           GetUser(event).ID,
			OrderDescription: orderOption,
			Status:           models.StatusWaiting,
			SourceServer:     event.GuildID,
			SourceChannel:    event.ChannelID,
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

		unixTimeString := fmt.Sprintf("%d", time.Now().Unix())

		session.ChannelMessageSendEmbed(config.GetConfig().KitchenChannelID, &discordgo.MessageEmbed{
			Title: "Order Created!",
			Description: fmt.Sprintf("Order ID: %d", order.ID) +
				"\nOrdered at: <t:" + unixTimeString + ":f>" +
				"\nPlaced: <t:" + unixTimeString + ":R>",
			Color: 0x00ff00,
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Order Created by " + DisplayName(event),
				IconURL: GetUser(event).AvatarURL("256"),
			},
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Sandwich Delivery",
				IconURL: session.State.User.AvatarURL("256"),
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Order:",
					Value:  orderOption,
					Inline: false,
				},
			},
		})
		break
	case "cancel":
		var order models.Order

		resp := database.GetDB().Find(&order, "user_id = ?", GetUser(event).ID)
		if resp.RowsAffected == 0 {
			session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						// todo https://github.com/refractored/sandwich-delivery/issues/5
						{
							Title:       "Error!",
							Description: "You do not have a pending order!",
							Color:       0xff2c2c,
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
		order.Status = models.StatusCancelled
		database.GetDB().Save(&order)

		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					// todo https://github.com/refractored/sandwich-delivery/issues/5
					{
						Title:       "Order Canceled!",
						Description: "You now can create a new order!",
						Color:       0x00ff00, // Green color
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
		break
	}
}
