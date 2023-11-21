package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"sandwich-delivery/src/config"
	"sandwich-delivery/src/models"
)

var database *gorm.DB

func Init(config *config.Config) *gorm.DB {
	log.Println("Opening Connection...")

	database, err := gorm.Open(mysql.New(mysql.Config{DSN: config.MySQLDSN}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	migrateTables()

	return database
}

func migrateTables() {
	log.Println("Migrating Blacklist Database...")
	err := GetDB().AutoMigrate(&models.BlacklistUser{})
	if err != nil {
		log.Fatal("Error migrating database:", err)
		return
	}

	log.Println("Migrating Orders Database...")
	err = GetDB().AutoMigrate(&models.Order{})
	if err != nil {
		log.Fatal("Error migrating database:", err)
		return
	}
}

func GetDB() *gorm.DB {
	return database
}
