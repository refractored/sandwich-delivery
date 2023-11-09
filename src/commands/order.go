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

	switch args[2] {

	case "cancel":
		err := db.Table("orders").Select("user_id").Where("user_id = ?", m.Author.ID).First(&user)
		if err.Error != nil {
			s.ChannelMessageSend(m.ChannelID, "You do not have a pending order!")
			return
		}
		err2 := db.Delete(&models.Order{}, "user_id = ?", m.Author.ID).Error
		if err2 != nil {
			s.ChannelMessageSend(m.ChannelID, "Error deleting order! Please try again.")
			return
		}
		s.ChannelMessageSend(m.ChannelID, "Canceling Order.")

	case "purge":
		s.ChannelMessageSend(m.ChannelID, "LOL?? :3")

	default:
		var user models.Order
		result := db.Table("orders").Select("user_id").Where("user_id = ?", m.Author.ID).First(&user)
		if result.Error == nil {
			s.ChannelMessageSend(m.ChannelID, "You already have a pending order!")
			return
		}

		foodOrder := strings.Join(args[2:], " ")
		err := db.Table("orders").Create(&models.Order{
			UserID:           m.Author.ID,
			OrderDescription: foodOrder,
			DisplayName:      m.Author.Username,
		}).Error

		if err != nil {
			log.Fatalf("Error Creating Order: %v", err)
			s.ChannelMessageSend(m.ChannelID, "Error Creating Order!")
			return
		}

		s.ChannelMessageSend(m.ChannelID, "Successfully created order!")
	}
}
