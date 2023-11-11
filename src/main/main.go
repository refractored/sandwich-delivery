package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"go-discord-bot/src/commands"
	"go-discord-bot/src/config"
	"go-discord-bot/src/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	startTime := time.Now()
	configPath := "config.json"
	log.Println("Loading Config...")
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}
	log.Println("Migrating Opening Connection...")
	db, err := gorm.Open(mysql.Open(cfg.MySQLDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	log.Println("Migrating Blacklist Database...")
	migrateBlacklistErr := db.AutoMigrate(&models.BlacklistUser{})
	if migrateBlacklistErr != nil {
		log.Fatal("Error migrating database:", err)
		return
	}
	log.Println("Migrating Orders Database...")
	migrateOrderErr := db.AutoMigrate(&models.Order{})
	if migrateOrderErr != nil {
		log.Fatal("Error migrating database:", err)
		return
	}

	log.Println("Opening Session on Discord...")
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
	sess.UpdateGameStatus(0, "Bot started!")
	sess.ChannelMessageSend("1171665367454716016", startupMessage)

	go func() {
		updateStatusPeriodically(sess, db)
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
