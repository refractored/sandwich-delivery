package commands

import (
	"github.com/bwmarrin/discordgo"
	"log"
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

	order := event.ApplicationCommandData().Options[0].StringValue()

	var user models.Order

	resp := database.GetDB().First(&user, "user_id = ?", GetUser(event).ID)
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

	resp = database.GetDB().Create(&models.Order{
		UserID:           GetUser(event).ID,
		OrderDescription: order,
		Username:         GetUser(event).Username,
		Discriminator:    GetUser(event).Discriminator,
		ServerID:         event.GuildID,
		ChannelID:        event.ChannelID,
	})
	if resp.Error != nil {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					// todo https://github.com/refractored/sandwich-delivery/issues/5
					{
						Title: "Error!",
						Description: "There was a problem creating your order! Please try again.\n" +
							"If this issue persists, Please report it!",
						Color: 0xff2c2c, // Green color
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
		log.Printf("Error Creating Order: %v", resp.Error)
		return
	}

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
							Value:  order,
							Inline: false,
						},
					},
				},
			},
		},
	})
}
