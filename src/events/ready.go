package events

import (
	"github.com/bwmarrin/discordgo"
)

func HandleReady(session *discordgo.Session, event *discordgo.Ready) {
	println("Logged in as: " + session.State.User.Username + "#" + session.State.User.Discriminator)
}
