package commands

import (
	"github.com/bwmarrin/discordgo"
	"go-discord-bot/src/models"
	"gorm.io/gorm"
	"log"
	"strings"
)

func OrderCommand(s *discordgo.Session, m *discordgo.MessageCreate, db *gorm.DB) {

	args := strings.Split(m.Content, " ")
	var user models.Order
	var displayname string
	if m.Author.Discriminator != "0" {
		displayname = m.Author.Username + "#" + m.Author.Discriminator
	} else {
		displayname = m.Author.Username
	}

	switch args[2] {

	case "cancel":
		err := db.Table("orders").Select("user_id").Where("user_id = ?", m.Author.ID).First(&user)
		if err.Error != nil {
			noPendingOrders := &discordgo.MessageEmbed{
				Title:       "Error!",
				Description: "You do not have a pending order!",
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
			s.ChannelMessageSendEmbed(m.ChannelID, noPendingOrders)
			return
		}
		err2 := db.Delete(&models.Order{}, "user_id = ?", m.Author.ID).Error
		if err2 != nil {
			errorDeletingOrder := &discordgo.MessageEmbed{
				Title: "Error!",
				Description: "There was a problem deleting your order! Please try again.\n" +
					"If this issue persists, Please report it!", Color: 0xff2c2c, // Green color
				Footer: &discordgo.MessageEmbedFooter{
					Text:    "Executed by " + displayname,
					IconURL: m.Author.AvatarURL("256"),
				},
				Author: &discordgo.MessageEmbedAuthor{
					Name:    "Sandwich Delivery",
					IconURL: s.State.User.AvatarURL("256"),
				},
			}
			s.ChannelMessageSendEmbed(m.ChannelID, errorDeletingOrder)
			return
		}
		cancelOrder := &discordgo.MessageEmbed{
			Title:       "Order Canceled!",
			Description: "You now can create a new order!",
			Color:       0x00ff00, // Green color
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Executed by " + displayname,
				IconURL: m.Author.AvatarURL("256"),
			},
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Sandwich Delivery",
				IconURL: s.State.User.AvatarURL("256"),
			},
		}
		s.ChannelMessageSendEmbed(m.ChannelID, cancelOrder)

	case "purge":
		s.ChannelMessageSend(m.ChannelID, "LOL?? :3")

	default:
		var user models.Order

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

		foodOrder := strings.Join(args[2:], " ")
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
			log.Println("Error Creating Order: %v", err)
			return
		}

		s.ChannelMessageSendEmbed(m.ChannelID, orderConformation)
	}
}
