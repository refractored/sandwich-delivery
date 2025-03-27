package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"runtime"
	"sandwich-delivery/src/database"
	"sandwich-delivery/src/models"
	"strconv"
)

type InfoCommand struct{}

func (c InfoCommand) getName() string {
	return "info"
}

func (c InfoCommand) getCommandData() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{Name: c.getName(),
		Description: "Lookup info of the bot or a user.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "user",
				Description: "Lookup an user.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionUser,
						Name:        "user",
						Description: "The user to lookup.",
						Required:    true,
					},
				},
				Type: discordgo.ApplicationCommandOptionSubCommand,
			},
			{
				Name:        "user",
				Description: "Lookup an order.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "order",
						Description: "The ID of the order to lookup.",
						Required:    true,
					},
				},
				Type: discordgo.ApplicationCommandOptionSubCommand,
			},
			{
				Name:        "bot",
				Description: "Lookup Bot Data",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
		},
	}
}

func (c InfoCommand) registerGuild() string {
	return ""
}

func (c InfoCommand) permissionLevel() models.UserPermissionLevel {
	return models.PermissionLevelUser
}

func (c InfoCommand) execute(session *discordgo.Session, event *discordgo.InteractionCreate) {
	switch event.ApplicationCommandData().Options[0].Name {
	case "bot":
		InfoBot(session, event)
		break
	case "user":
		InfoUser(session, event)
		break
	case "order":
		InfoOrder(session, event)
		break
	}
}

func InfoBot(session *discordgo.Session, event *discordgo.InteractionCreate) {
	var pendingOrderCount int64
	var completedOrderCount int64
	var canceledOrderCount int64
	var userCount int64

	result := database.GetDB().Model(&models.Order{}).Where("status < ?", models.StatusDelivered).Count(&pendingOrderCount)
	if result.Error != nil {
		log.Println("Error counting orders:", result.Error)
	}
	result = database.GetDB().Model(&models.Order{}).Where("status = ?", models.StatusDelivered).Count(&completedOrderCount)
	if result.Error != nil {
		log.Println("Error counting orders:", result.Error)
	}
	result = database.GetDB().Model(&models.Order{}).Where("status > ?", models.StatusDelivered).Count(&canceledOrderCount)
	if result.Error != nil {
		log.Println("Error counting orders:", result.Error)
	}
	result = database.GetDB().Model(&models.User{}).Count(&userCount)
	if result.Error != nil {
		log.Println("Error counting orders:", result.Error)
	}
	session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "Bot Information",
					Description: "Bot Name: " + session.State.User.Username + "#" + session.State.User.Discriminator + "\n" +
						"Guilds: " + strconv.Itoa(len(session.State.Guilds)) + "\n" +
						"Pending Orders: " + strconv.Itoa(int(pendingOrderCount)) + "\n" +
						"Completed Orders: " + strconv.Itoa(int(completedOrderCount)) + "\n" +
						"Canceled Orders: " + strconv.Itoa(int(canceledOrderCount)) + "\n" +
						"Sandwich Accounts: " + strconv.Itoa(int(userCount)) + "\n" +
						fmt.Sprintf("Library: DiscordGo (%s)", discordgo.VERSION) + "\n" +
						"Runtime: " + runtime.Version() + " " + runtime.GOARCH + "\n",

					Color: 0x00ff00,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Executed by " + DisplayName(event),
						IconURL: GetUser(event).AvatarURL("256"),
					},
					Author: &discordgo.MessageEmbedAuthor{
						Name:    "Sandwich Delivery",
						IconURL: session.State.User.AvatarURL("256"),
					},
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: session.State.User.AvatarURL("256"),
					},
				},
			},
		},
	})
}
func InfoUser(session *discordgo.Session, event *discordgo.InteractionCreate) {
	if event.ApplicationCommandData().Options[0].Options[0].UserValue(session).ID == session.State.User.ID {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You can lookup the bots information with /info bot",
			},
		})
		return
	}
	var user models.User

	resp := database.GetDB().Find(&user, "user_id = ?", event.ApplicationCommandData().Options[0].Options[0].UserValue(session).ID)
	if resp.RowsAffected == 0 {
		user.UserID = event.ApplicationCommandData().Options[0].Options[0].UserValue(session).ID

		database.GetDB().Save(&user)
	}
	userarg, _ := session.User(user.UserID)
	var dailyclaimed string
	if user.DailyClaimedAt != nil {
		dailyclaimed = user.DailyClaimedAt.String()
	} else {
		dailyclaimed = "Never"
	}
	session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "User Information",
					Description: "Orders Created: " + strconv.Itoa(int(user.OrdersCreated)) + "\n" +
						"Orders Accepted: " + strconv.Itoa(int(user.OrdersAccepted)) + "\n" +
						"Credits: " + strconv.Itoa(int(user.Credits)) + "\n" +
						"DB ID: " + strconv.Itoa(int(user.ID)) + "\n" +
						"Blacklisted: " + strconv.FormatBool(user.IsBlacklisted) + "\n" +
						"Permission Level: " + strconv.Itoa(int(user.OrdersAccepted)) + "\n" +
						"Daily Claimed At: " + dailyclaimed + "\n" +
						"Sandwich Account Creation: " + user.CreatedAt.String() + "\n",
					Color: 0x00ff00,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Executed by " + DisplayName(event),
						IconURL: GetUser(event).AvatarURL("256"),
					},
					Author: &discordgo.MessageEmbedAuthor{
						Name:    "Sandwich Delivery",
						IconURL: session.State.User.AvatarURL("256"),
					},
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: userarg.AvatarURL("256"),
					},
				},
			},
		},
	})
}

func InfoOrder(session *discordgo.Session, event *discordgo.InteractionCreate) {
	var order models.Order

	resp := database.GetDB().Find(&order, "id = ?", event.ApplicationCommandData().Options[0].Options[0].IntValue())

	if resp.RowsAffected == 0 {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Order not found.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	var user, _ = session.User(order.UserID)
	var server, _ = session.Guild(order.SourceServer)
	var channel, _ = session.Channel(order.SourceChannel)
	var assignee, _ = session.User(order.Assignee)

	var description string = "Order ID: " + strconv.Itoa(int(order.ID)) + "\n" +
		"Description: " + order.OrderDescription + "\n" +
		"Status: " + order.Status.String() + "\n" +
		"Created At: " + order.CreatedAt.String() + "\n"

	if user != nil {
		description += "User: " + user.Mention() + "\n"
	}

	if server != nil {
		description += "Server: " + server.Name + "\n"
	}

	if channel != nil {
		description += "Channel: " + channel.Name + "\n"
	}

	if assignee != nil {
		description += "Assignee: " + assignee.Mention() + "\n"
	}

	if order.AcceptedAt != nil {
		description += "Accepted At: " + order.AcceptedAt.String() + "\n"
	}

	if order.InTransitAt != nil {
		description += "In Transit At: " + order.InTransitAt.String() + "\n"
	}

	if order.DeliveredAt != nil {
		description += "Delivered At: " + order.DeliveredAt.String() + "\n"
	}

	session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "User Information",
					Description: description,
					Color:       0x00ff00,
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
