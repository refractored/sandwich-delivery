package commands

import (
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
	"log"
	"sandwich-delivery/src/models"
	"strings"
)

func OrderCommand(s *discordgo.Session, m *discordgo.MessageCreate, db *gorm.DB) {

	var user models.Order
	args := strings.Split(m.Content, " ")
	foodOrder := strings.Join(args[1:], " ")
	displayname := DisplayName(s, m)

	if len(args[1]) < 3 {
		minCharacters := &discordgo.MessageEmbed{
			Title:       "Error!",
			Description: "Your order needs to be 3 characters or more!",
			Color:       0xff2c2c, // Green color
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Executed by " + displayname,
				IconURL: m.Author.AvatarURL("256"),
			},
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Sandwich Delivery",
				IconURL: s.State.User.AvatarURL("256"),
			},
		}
		s.ChannelMessageSendEmbed(m.ChannelID, minCharacters)
		return
	}

	result := db.Table("orders").Select("user_id").Where("user_id = ?", m.Author.ID).First(&user)
	if result.Error == nil {
		pendingOrder := &discordgo.MessageEmbed{
			Title:       "Error!",
			Description: "You already have a pending order!",
			Color:       0xff2c2c, // Green color
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Executed by " + displayname,
				IconURL: m.Author.AvatarURL("256"),
			},
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Sandwich Delivery",
				IconURL: s.State.User.AvatarURL("256"),
			},
		}
		s.ChannelMessageSendEmbed(m.ChannelID, pendingOrder)
		return
	}

	orderConformation := &discordgo.MessageEmbed{
		Title: "Order Placed!",
		Description: "Thanks for placing your order!" +
			"\nPlease give our staff some time!",
		Color: 0x00ff00, // Green color
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Executed by " + displayname,
			IconURL: m.Author.AvatarURL("256"),
		},
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "Sandwich Delivery",
			IconURL: s.State.User.AvatarURL("256"),
		},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Your Order:",
				Value:  foodOrder,
				Inline: false,
			},
		},
	}
	err := db.Table("orders").Create(&models.Order{
		UserID:           m.Author.ID,
		OrderDescription: foodOrder,
		Username:         m.Author.Username,
		Discriminator:    m.Author.Discriminator,
		ServerID:         m.GuildID,
		ChannelID:        m.ChannelID,
	}).Error
	if err != nil {
		errorCreatingOrder := &discordgo.MessageEmbed{
			Title: "Error!",
			Description: "There was a problem creating your order! Please try again.\n" +
				"If this issue persists, Please report it!",
			Color: 0xff2c2c, // Green color
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Executed by " + displayname,
				IconURL: m.Author.AvatarURL("256"),
			},
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Sandwich Delivery",
				IconURL: s.State.User.AvatarURL("256"),
			},
		}
		s.ChannelMessageSendEmbed(m.ChannelID, errorCreatingOrder)
		log.Printf("Error Creating Order: %v", err)
		return
	}
	s.ChannelMessageSendEmbed(m.ChannelID, orderConformation)
}
