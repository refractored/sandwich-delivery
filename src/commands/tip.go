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
	var minvalue float64 = 1
	return &discordgo.ApplicationCommand{
		Name:        c.getName(),
		Description: "Tips your last order",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "amount",
				Description: "Amount to tip",
				Required:    true,
				MinValue:    &minvalue,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "id",
				Description: "Order ID",
				Required:    false,
			},
		},
	}
}

func (c TipCommand) registerGuild() string {
	return ""
}

func (c TipCommand) permissionLevel() models.UserPermissionLevel {
	return models.PermissionLevelUser
}

func (c TipCommand) execute(session *discordgo.Session, event *discordgo.InteractionCreate) {
	amount := uint32(event.ApplicationCommandData().Options[0].IntValue())
	var id uint64
	var order models.Order
	var customer models.User
	var employee models.User

	if len(event.ApplicationCommandData().Options) == 1 {
		resp := database.GetDB().Find(&order, "user_id = ?", GetUser(event).ID)
		if resp.RowsAffected == 0 {
			session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You've never ordered anything before!",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}
	} else {
		id = uint64(event.ApplicationCommandData().Options[1].IntValue())
		resp := database.GetDB().Find(&order, "id = ?", id)
		if resp.RowsAffected == 0 {
			session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "This order does not exist!",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}
		if order.UserID != GetUser(event).ID {
			session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You can only tip your own orders!",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}

	}
	if order.Status > models.StatusDelivered {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You can only tip delivered orders!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}
	if order.Tipped == true {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You've already tipped this order! You can only tip once!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}
	database.GetDB().Find(&customer, "user_id = ?", GetUser(event).ID)
	database.GetDB().Find(&employee, "user_id = ?", order.Assignee)
	if customer.Credits < amount {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You don't have enough credits to tip!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}
	customer.Credits = customer.Credits - amount
	employee.Credits = employee.Credits + amount
	order.Tipped = true
	database.GetDB().Save(&order)
	database.GetDB().Save(&customer)
	database.GetDB().Save(&employee)
	session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Tipped!",
		},
	})
}
