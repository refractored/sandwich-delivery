package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"go-discord-bot/src/commands"
	"go-discord-bot/src/config"
	"go-discord-bot/src/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	startTime := time.Now()
	configPath := "config.json"
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}
	db, err := gorm.Open(mysql.Open(cfg.MySQLDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	db.AutoMigrate(&models.BlacklistUser{})
	db.AutoMigrate(&models.Order{})

	sess, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		log.Fatal(err)
	}

	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		commands.HandleCommand(sess, m, &cfg, db)
	})

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}

	startupTime := time.Since(startTime)
	startupMessage := fmt.Sprintf("Bot started! (%[1]s)", startupTime)
	sess.UpdateGameStatus(0, "HAI GUYS :3")
	sess.ChannelMessageSend("1171665367454716016", startupMessage)
	fmt.Println("Bot is running!")

	defer func() {
		fmt.Println("Bot is shutting down...")
		sess.ChannelMessageSend("1171665367454716016", "Shutting down...")
		sess.Close()
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
