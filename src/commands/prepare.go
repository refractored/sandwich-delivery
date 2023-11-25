package commands

import (
	"github.com/bwmarrin/discordgo"
	"sandwich-delivery/src/config"
	"sandwich-delivery/src/database"
	"sandwich-delivery/src/models"
	"time"
)

type PrepareOrderCommand struct{}

func (c PrepareOrderCommand) getName() string {
	return "prepareorder"
}

func (c PrepareOrderCommand) getCommandData() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{Name: c.getName(), Description: "Prepares the invite for the order so you can send it to the customer."}
}

func (c PrepareOrderCommand) registerGuild() string {
	return config.GetConfig().GuildID

}

func (c PrepareOrderCommand) permissionLevel() models.UserPermissionLevel {
	return models.PermissionLevelArtist
}

func (c PrepareOrderCommand) execute(session *discordgo.Session, event *discordgo.InteractionCreate) {
	var order models.Order
	resp := database.GetDB().First(&order, "assignee = ?", GetUser(event).ID)

	if resp.RowsAffected == 0 {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You dont have an order to prepare!",
			},
		})
		return
	}

	if order.Status != models.StatusAccepted {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You can only prepare orders that you already accepted!",
			},
		})
		return
	}
	_, err := session.ChannelMessageSendComplex(order.SourceChannel, &discordgo.MessageSend{
		Content: "<@" + order.UserID + ">",
		Embed: &discordgo.MessageEmbed{
			Title: "Order Out for Delivery!",
			Description: "Your order is ready and should arrive shortly!" + "\n" +
				"Dont forget you can tip!",
			Color: 0x00ff00,
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Order being delivered by " + DisplayName(event),
				IconURL: GetUser(event).AvatarURL("256"),
			},
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Sandwich Delivery",
				IconURL: session.State.User.AvatarURL("256"),
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Order:",
					Value:  order.OrderDescription,
					Inline: false,
				},
			},
		},
	})
	if err != nil {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Order Location Invalid! Order Deleted.",
			},
		})
		dmMessage, _ := session.UserChannelCreate(order.UserID)
		session.ChannelMessageSendComplex(dmMessage.ID, &discordgo.MessageSend{
			Content: "<@" + order.UserID + ">",
			Embed: &discordgo.MessageEmbed{
				Title: "Order Error!",
				Description: "Your order could not be completed!" + "\n" +
					"The bot was unable to send a message into the channel you ordered from so it was deleted!",
				Color: 0xff2c2c,
				Footer: &discordgo.MessageEmbedFooter{
					Text:    "Prepare Command ran by " + DisplayName(event),
					IconURL: GetUser(event).AvatarURL("256"),
				},
				Author: &discordgo.MessageEmbedAuthor{
					Name:    "Sandwich Delivery",
					IconURL: session.State.User.AvatarURL("256"),
				},
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Order:",
						Value:  order.OrderDescription,
						Inline: false,
					},
				},
			},
		})
		order.Status = models.StatusError
		database.GetDB().Save(&order)
		database.GetDB().Delete(&order)
		return
	}
	// TODO: Check if the order was placed inside the one set in the config
	invite, err := session.ChannelInviteCreate(order.SourceChannel, discordgo.Invite{MaxUses: 1})
	if err != nil {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Could not create invite! Order Deleted.",
			},
		})
		dmMessage, _ := session.UserChannelCreate(order.UserID)
		session.ChannelMessageSendComplex(dmMessage.ID, &discordgo.MessageSend{
			Content: "<@" + order.UserID + ">",
			Embed: &discordgo.MessageEmbed{
				Title: "Order Error!",
				Description: "Your order could not be completed!" + "\n" +
					"The bot was unable to make a invite for the channel you ordered from so it was deleted!",
				Color: 0xff2c2c,
				Footer: &discordgo.MessageEmbedFooter{
					Text:    "Prepare Command ran by " + DisplayName(event),
					IconURL: GetUser(event).AvatarURL("256"),
				},
				Author: &discordgo.MessageEmbedAuthor{
					Name:    "Sandwich Delivery",
					IconURL: session.State.User.AvatarURL("256"),
				},
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Order:",
						Value:  order.OrderDescription,
						Inline: false,
					},
				},
			},
		})
		order.Status = models.StatusError
		database.GetDB().Save(&order)
		database.GetDB().Delete(&order)
		return
	}
	dmMessage, _ := session.UserChannelCreate(order.UserID)
	_, err = session.ChannelMessageSendComplex(dmMessage.ID, &discordgo.MessageSend{
		Content: "<@" + order.UserID + ">" + " https://discord.gg/" + invite.Code,
		Embed: &discordgo.MessageEmbed{
			Title: "Ding ding!",
			Description: "https://discord.gg/" + invite.Code + "\n" +
				"The customer awaits their order",
			Color: 0x00ff00,
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Prepare Command ran by " + DisplayName(event),
				IconURL: GetUser(event).AvatarURL("256"),
			},
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Sandwich Delivery",
				IconURL: session.State.User.AvatarURL("256"),
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Order:",
					Value:  order.OrderDescription,
					Inline: false,
				},
			},
		},
	})
	if err != nil {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please make sure DMs are enabled!",
			},
		})
		return
	}
	var time = time.Now()
	order.InTransitAt = &time
	order.Status = models.StatusInTransit
	resp = database.GetDB().Save(&order)

	session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Check your DMs!",
		},
	})
}
