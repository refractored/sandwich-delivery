package commands

import (
	"github.com/bwmarrin/discordgo"
	"sandwich-delivery/src/database"
	"sandwich-delivery/src/models"
	"strconv"
	"time"
)

type StatusCommand struct{}

func (c StatusCommand) getName() string {
	return "status"
}

func (c StatusCommand) getCommandData() *discordgo.ApplicationCommand {
	DMPermission := new(bool)
	*DMPermission = false
	return &discordgo.ApplicationCommand{
		Name:        c.getName(),
		Description: "Update the status of a customer's order.",
		Options: []*discordgo.ApplicationCommandOption{

			{
				Name:        "accept",
				Description: "Assign an order to ONLY you.",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "id",
						Description: "The ID of the order to accept.",
						Required:    true,
					},
				},
			},
			{
				Name:        "transit",
				Description: "Get the invite of the order location and mark it as in transit.",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
			{
				Name:        "delivered",
				Description: "Mark an order as delivered.",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
		},
		DMPermission: DMPermission,
	}
}

func (c StatusCommand) registerGuild() string {
	return ""
}

func (c StatusCommand) permissionLevel() models.UserPermissionLevel {
	return models.PermissionLevelArtist
}

func (c StatusCommand) execute(session *discordgo.Session, event *discordgo.InteractionCreate) {
	switch event.ApplicationCommandData().Options[0].Name {
	case "accept":
		OrderStatusAccept(session, event)
		break
	case "transit":
		OrderStatusInTransit(session, event)
		break
	case "delivered":
		OrderStatusDelivered(session, event)
		break
	}
}
func OrderStatusAccept(session *discordgo.Session, event *discordgo.InteractionCreate) {
	var order models.Order

	orderID := models.UserPermissionLevel(event.ApplicationCommandData().Options[0].Options[0].IntValue())

	// Check if artist already has an order
	resp := database.GetDB().Find(&order, "assignee = ? AND status < ?", GetUser(event).ID, models.StatusDelivered)

	if resp.RowsAffected != 0 {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You can only accept one order at a time!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}
	resp = database.GetDB().Find(&order, "id = ?", orderID)
	if resp.RowsAffected == 0 {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Order not found.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	if order.Status > models.StatusDelivered {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This order was cancelled!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	if order.Status != models.StatusWaiting {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This order can no longer be accepted! It's either already accepted or delivered.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}
	_, err := session.ChannelMessageSendComplex(order.SourceChannel, &discordgo.MessageSend{
		Content: "<@" + order.UserID + ">",
		Embed: &discordgo.MessageEmbed{
			Title: "Order Accepted!",
			Description: "Your order has been accepted!" + "\n" +
				"It's currently being prepared and will out for delivery soon!",
			Color: 0x00ff00,
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Order Accepted by " + DisplayName(event),
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
				{
					Name:   "ID:",
					Value:  strconv.Itoa(int(order.ID)),
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
				Flags:   discordgo.MessageFlagsEphemeral,
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
					Text:    "Accept Command ran by " + DisplayName(event),
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
		return
	}
	order.Assignee = GetUser(event).ID
	var timeNow = time.Now()
	order.AcceptedAt = &timeNow
	order.Status = models.StatusAccepted
	resp = database.GetDB().Save(&order)
	session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Order accepted!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
func OrderStatusInTransit(session *discordgo.Session, event *discordgo.InteractionCreate) {
	var order models.Order

	resp := database.GetDB().Find(&order, "assignee = ? AND status < ?", GetUser(event).ID, models.StatusDelivered)

	if resp.RowsAffected == 0 {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You don't have an order mark as in transit!" + "\n" +
					"Use /status accept to accept an order first!",
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	if order.Status != models.StatusAccepted {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You can only prepare orders that you accepted!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}
	_, err := session.ChannelMessageSendComplex(order.SourceChannel, &discordgo.MessageSend{
		Content: "<@" + order.UserID + ">",
		Embed: &discordgo.MessageEmbed{
			Title: "Order Out for Delivery!",
			Description: "Your order is ready and should arrive shortly!" + "\n" +
				"Don't forget you can tip!",
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
				{
					Name:   "ID:",
					Value:  strconv.Itoa(int(order.ID)),
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
				Flags:   discordgo.MessageFlagsEphemeral,
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
		return
	}
	invite, err := session.ChannelInviteCreate(order.SourceChannel, discordgo.Invite{MaxUses: 1, Temporary: true, Unique: true, MaxAge: 300})
	if err != nil {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Could not create invite! Order Deleted.",
				Flags:   discordgo.MessageFlagsEphemeral,
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
					{
						Name:   "ID:",
						Value:  strconv.Itoa(int(order.ID)),
						Inline: false,
					},
				},
			},
		})
		order.Status = models.StatusError
		database.GetDB().Save(&order)
		return
	}
	dmMessage, _ := session.UserChannelCreate(GetUser(event).ID)
	_, err = session.ChannelMessageSendComplex(dmMessage.ID, &discordgo.MessageSend{
		Content: "<@" + order.UserID + ">" + " https://discord.gg/" + invite.Code,
		Embed: &discordgo.MessageEmbed{
			Title: "Ding ding!",
			Description: "https://discord.gg/" + invite.Code + "\n" +
				"The customer awaits their order",
			Color: 0x00ff00,
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Transit Command ran by " + DisplayName(event),
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
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}
	var timeNow = time.Now()
	order.InTransitAt = &timeNow
	order.Status = models.StatusInTransit
	resp = database.GetDB().Save(&order)

	session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Check your DMs!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func OrderStatusDelivered(session *discordgo.Session, event *discordgo.InteractionCreate) {
	var order models.Order
	resp := database.GetDB().Find(&order, "assignee = ? AND status < ?", GetUser(event).ID, models.StatusDelivered)

	if resp.RowsAffected == 0 {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You don't have an order to mark as delivered!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	if order.Status != models.StatusInTransit {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You can only mark orders as delivered if they are in transit!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}
	session.ChannelMessageSendComplex(order.SourceChannel, &discordgo.MessageSend{
		Content: "<@" + order.UserID + ">",
		Embed: &discordgo.MessageEmbed{
			Title: "Order Delivered!",
			Description: "Your order has been marked as delivered!" + "\n" +
				"You may tip your artist with /tip!",
			Color: 0x00ff00,
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Order marked as delivered by " + DisplayName(event),
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
				{
					Name:   "ID:",
					Value:  strconv.Itoa(int(order.ID)),
					Inline: false,
				},
			},
		},
	})
	var time = time.Now()
	order.DeliveredAt = &time
	order.Status = models.StatusDelivered
	resp = database.GetDB().Save(&order)

	session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Marked order as finished!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
