package commands

import (
	"github.com/bwmarrin/discordgo"
	"sandwich-delivery/src/database"
	"sandwich-delivery/src/models"
	"strconv"
	"time"
)

type DeliverCommand struct{}

func (c DeliverCommand) getName() string {
	return "deliver"
}

func (c DeliverCommand) getCommandData() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{Name: c.getName(), Description: "Mark an order as delivered."}
}

func (c DeliverCommand) registerGuild() string {
	return ""
}

func (c DeliverCommand) permissionLevel() models.UserPermissionLevel {
	return models.PermissionLevelArtist
}

func (c DeliverCommand) execute(session *discordgo.Session, event *discordgo.InteractionCreate) {
	var order models.Order
	resp := database.GetDB().First(&order, "assignee = ? AND status < ?", GetUser(event).ID, models.StatusDelivered)

	if resp.RowsAffected == 0 {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You dont have an order to mark as delivered!",
			},
		})
		return
	}

	if order.Status != models.StatusInTransit {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You can only mark orders as delivered that are in transit!",
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
		},
	})
}
