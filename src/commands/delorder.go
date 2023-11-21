package commands

import (
	"github.com/bwmarrin/discordgo"
	"sandwich-delivery/src/database"
	"sandwich-delivery/src/models"
)

type DelOrderCommand struct{}

func (d *DelOrderCommand) getName() string {
	return "delorder"
}

func (d *DelOrderCommand) getCommandData() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{Name: d.getName(), Description: "Deletes your pending order."}
}

func (d *DelOrderCommand) execute(session *discordgo.Session, event *discordgo.InteractionCreate) {
	if InteractionIsDM(event) {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				// todo https://github.com/refractored/sandwich-delivery/issues/5
				Content: "This command can only be used in servers!",
			},
		})
		return
	}

	var user models.Order

	resp := database.GetDB().First(&user, "user_id = ?", GetUser(event).ID)
	if resp.Error != nil {
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
	resp = database.GetDB().Delete(&models.Order{}, "user_id = ?", GetUser(event).ID)
	if resp.Error != nil {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					// todo https://github.com/refractored/sandwich-delivery/issues/5
					{
						Title: "Error!",
						Description: "There was a problem deleting your order! Please try again.\n" +
							"If this issue persists, Please report it!", Color: 0xff2c2c, // Green color
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
}
