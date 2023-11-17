package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"sandwich-delivery/src/config"
	"sandwich-delivery/src/models"
)

var Database *gorm.DB

func Init(config config.Config) *gorm.DB {
	log.Println("Opening Connection...")
	Database, err := gorm.Open(mysql.Open(config.MySQLDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	migrateTables()

	return Database
}

func migrateTables() {
	log.Println("Migrating Blacklist Database...")
	err := Database.AutoMigrate(&models.BlacklistUser{})
	if err != nil {
		log.Fatal("Error migrating database:", err)
		return
	}

	log.Println("Migrating Orders Database...")
	err = Database.AutoMigrate(&models.Order{})
	if err != nil {
		log.Fatal("Error migrating database:", err)
		return
	}
}
