package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"sandwich-delivery/src/config"
	"sandwich-delivery/src/database"
	"sandwich-delivery/src/models"
	"time"
)

type DailyCommand struct{}

func (c DailyCommand) getName() string {
	return "daily"
}

func (c DailyCommand) getCommandData() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{Name: c.getName(), Description: "Collect your daily rewards!"}
}

func (c DailyCommand) registerGuild() string {
	return ""
}

func (c DailyCommand) permissionLevel() models.UserPermissionLevel {
	return models.PermissionLevelUser
}

func (c DailyCommand) execute(session *discordgo.Session, event *discordgo.InteractionCreate) {
	if *config.GetConfig().DailyTokens == 0 {
		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Error!",
						Description: "Daily rewards are disabled!",
						Color:       0xff2c2c,
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "Executed by " + DisplayName(event),
							IconURL: GetUser(event).AvatarURL("256"),
						},
					},
				},
			},
		})
	}

	var user models.User

	database.GetDB().First(&user, "user_id = ?", GetUser(event).ID)
	claimedAt := *user.DailyClaimedAt
	elapsed := time.Since(claimedAt)

	if elapsed.Hours() >= 24 {
		user.Credits = user.Credits + *config.GetConfig().DailyTokens
		TimeNow := *user.DailyClaimedAt

		user.DailyClaimedAt = &TimeNow
		database.GetDB().Save(&user)

		var tokenMessage string

		if *config.GetConfig().DailyTokens == 1 {
			tokenMessage = "+1 Credit"
		} else {
			tokenMessage = fmt.Sprintf("+%d Credits", *config.GetConfig().DailyTokens)
		}

		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title: "Reward Claimed!",
						Description: tokenMessage +
							"\nCome back in 24 hours for another reward!",
						Color: 0x00ff00,
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
	} else {
		remainingTime := 24*time.Hour - elapsed
		remainingTimeString := fmt.Sprintf("%d hours, %d minutes, %d seconds",
			int(remainingTime.Hours()), int(remainingTime.Minutes())%60, int(remainingTime.Seconds())%60)

		session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title: "Error!",
						Description: "You need to wait a bit longer!" +
							"\nYou can claim your daily reward in " + remainingTimeString,
						Color: 0xff2c2c,
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
}
