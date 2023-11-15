package commands

import (
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
	"sandwich-delivery/src/models"
)

func DelOrderCommand(s *discordgo.Session, m *discordgo.MessageCreate, db *gorm.DB) {

	//args := strings.Split(m.Content, " ")
	var user models.Order
	var displayname = DisplayName(s, m)

	err := db.Table("orders").Select("user_id").Where("user_id = ?", m.Author.ID).First(&user)
	if err.Error != nil {
		noPendingOrders := &discordgo.MessageEmbed{
			Title:       "Error!",
			Description: "You do not have a pending order!",
			Color:       0xff2c2c,
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

}
