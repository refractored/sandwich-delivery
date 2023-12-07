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
	"sandwich-delivery/src/models"
	"strconv"
	"syscall"
	"time"
)

func main() {
	startTime := time.Now()
	configPath := "config.json"

	log.Println("Loading Config...")
	cfg, err := config.LoadConfig(configPath)

	log.Println("Verifying Config...")
	success, err := config.VerifyConfig(cfg)

	if err != nil {
		log.Fatal(err)
	}
	if !success {
		log.Fatal("Config verification failed!")
	}

	log.Println("Initializing Database...")
	database.Init()

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

	sess.UpdateGameStatus(0, "Bot started!")
	_, err = sess.ChannelMessageSend(config.GetConfig().StartupChannelID, fmt.Sprintf("Bot started! (%[1]s)", time.Since(startTime)))
	if err != nil && config.GetConfig().StartupChannelID != "" {
		log.Println("Error sending startup message:", err)
	}

	go func() {
		updateStatusPeriodically(sess, database.GetDB())
	}()

	log.Println("Bot is running!")

	defer func() {
		log.Println("Bot is shutting down...")
		_, err := sess.ChannelMessageSend(config.GetConfig().StartupChannelID, "Shutting down...")
		if err != nil && config.GetConfig().StartupChannelID != "" {
			log.Println("Error sending shutdown message:", err)
		}
		sess.Close()
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func updateStatusPeriodically(s *discordgo.Session, db *gorm.DB) {
	var updateInterval = 2 * time.Minute

	for {
		var orderCount int64

		result := db.Model(&models.Order{}).Where("status < ?", models.StatusDelivered).Count(&orderCount)
		if result.Error != nil {
			log.Println("Error counting orders:", result.Error)
			continue
		}
		orderCountString := strconv.Itoa(int(orderCount))

		s.UpdateGameStatus(0, "Orders: "+orderCountString)

		log.Println("Bot status updated. Orders:", orderCount)
		time.Sleep(updateInterval)
	}

}
