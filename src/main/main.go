package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
	"sandwich-delivery/src/commands"
	"sandwich-delivery/src/config"
	"sandwich-delivery/src/database"
	"strconv"
	"syscall"
	"time"
)

func main() {
	startTime := time.Now()
	configPath := "config.json"

	log.Println("Loading Config...")
	cfg, err := config.LoadConfig(configPath)

	log.Println("Initializing Database...")
	database.Init(cfg)

	log.Println("Opening Session on Discord...")
	sess, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		log.Fatal(err)
	}

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Registering Commands...")
	commands.RegisterCommands(sess)

	sess.AddHandler(func(session *discordgo.Session, event *discordgo.InteractionCreate) {
		commands.HandleCommand(session, event)
	})

	startupTime := time.Since(startTime)
	startupMessage := fmt.Sprintf("Bot started! (%[1]s)", startupTime)
	sess.UpdateGameStatus(0, "Bot started!")
	sess.ChannelMessageSend("1171665367454716016", startupMessage)

	go func() {
		updateStatusPeriodically(sess, database.GetDB())
	}()

	log.Println("Bot is running!")

	defer func() {
		log.Println("Bot is shutting down...")
		sess.ChannelMessageSend("1171665367454716016", "Shutting down...")
		sess.Close()
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func updateStatusPeriodically(s *discordgo.Session, db *gorm.DB) {
	updateInterval := 2 * time.Minute
	ticker := time.NewTicker(updateInterval)

	for {
		select {
		case <-ticker.C:
			var orderCount int64
			result := db.Table("orders").Count(&orderCount)
			if result.Error != nil {
				log.Println("Error counting orders:", result.Error)
				continue
			}

			orderCountString := strconv.Itoa(int(orderCount))

			s.UpdateGameStatus(0, "Orders: "+orderCountString)

			log.Println("Bot status updated. Orders:", orderCount)
		}
	}
}
