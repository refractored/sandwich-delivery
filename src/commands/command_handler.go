package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
)

var commands = map[string]Command{
	CoinflipCommand{}.getName(): CoinflipCommand{},
}

func RegisterCommands(session *discordgo.Session) {
	for n, d := range commands {
		_, err := session.ApplicationCommandCreate(session.State.User.ID, "", d.getCommandData())
		if err == nil {
			fmt.Printf("Unable to register command %s\n" + n)
			os.Exit(1)
		}
	}
}

func HandleCommand(session *discordgo.Session, event *discordgo.InteractionCreate) {
	command := commands[event.ApplicationCommandData().Name]
	if command == nil {
		session.ChannelMessageSendReply(event.ChannelID, "The command "+event.ApplicationCommandData().Name+" does not exist.", event.Message.MessageReference)
		session.ApplicationCommandDelete(session.State.User.ID, "", event.ApplicationCommandData().ID)
	}
	command.execute(session, event)
}
