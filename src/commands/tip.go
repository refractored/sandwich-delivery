package commands

import (
	"github.com/bwmarrin/discordgo"
	"sandwich-delivery/src/database"
	"sandwich-delivery/src/models"
)

type TipCommand struct{}

func (c TipCommand) getName() string {
	return "tip"
}

func (c TipCommand) getCommandData() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{Name: c.getName(), Description: "Tips your last order"}
}

func (c TipCommand) registerGuild() string {
	return ""
}

func (c TipCommand) permissionLevel() models.UserPermissionLevel {
	return models.PermissionLevelUser
}

func (c TipCommand) execute(session *discordgo.Session, event *discordgo.InteractionCreate) {
	var order models.Order
	var customer models.User
	var employee models.User
	resp := database.GetDB().Last(&order, "user_id = ?", GetUser(event).ID)
	if resp.RowsAffected == 0 {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You've never ordered anything before!",
			},
		})
		return
	}
	if order.Status != models.StatusDelivered {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You can only tip delivered orders!",
			},
		})
		return
	}
	if order.Tipped == true {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You've already tipped this order! You can only tip once!",
			},
		})
		return
	}
	resp = database.GetDB().First(&customer, "user_id = ?", GetUser(event).ID)
	resp = database.GetDB().First(&employee, "user_id = ?", order.Assignee)
	if customer.Credits < 1 {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You don't have enough credits to tip!",
			},
		})
		return
	}
	customer.Credits = customer.Credits - 1
	employee.Credits = employee.Credits + 1
	order.Tipped = true
	database.GetDB().Save(&order)
	database.GetDB().Save(&customer)
	database.GetDB().Save(&employee)
	session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "tipped!",
		},
	})
}
