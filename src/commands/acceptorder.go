package commands

import (
	"github.com/bwmarrin/discordgo"
	"sandwich-delivery/src/config"
	"sandwich-delivery/src/database"
	"sandwich-delivery/src/models"
	"strconv"
	"time"
)

type AcceptOrderCommand struct{}

func (c AcceptOrderCommand) getName() string {
	return "acceptorder"
}

func (c AcceptOrderCommand) getCommandData() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        c.getName(),
		Description: "Accepts the order and assigns it to you.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "id",
				Description: "The ID of the order to accept.",
				Required:    true,
			},
		},
	}
}

func (c AcceptOrderCommand) registerGuild() string {
	return config.GetConfig().GuildID

}

func (c AcceptOrderCommand) permissionLevel() models.UserPermissionLevel {
	return models.PermissionLevelArtist
}

func (c AcceptOrderCommand) execute(session *discordgo.Session, event *discordgo.InteractionCreate) {
	orderID := models.UserPermissionLevel(event.ApplicationCommandData().Options[0].IntValue())

	var order models.Order
	resp := database.GetDB().First(&order, "assignee = ? AND status < ?", GetUser(event).ID, models.StatusDelivered)

	if resp.RowsAffected != 0 {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You can only accept one order at a time!",
			},
		})
		return
	}

	resp = database.GetDB().First(&order, "id = ?", orderID)

	if resp.RowsAffected == 0 {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Order not found.",
			},
		})
		return
	}
	if order.Status != models.StatusWaiting {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This order can no longer be accepted!",
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
	var time = time.Now()
	order.AcceptedAt = &time
	order.Status = models.StatusAccepted
	resp = database.GetDB().Save(&order)

	session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Order accepted!",
		},
	})
}
