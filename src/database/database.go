package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"sandwich-delivery/src/config"
	"sandwich-delivery/src/models"
)

var database *gorm.DB

func Init() *gorm.DB {
	log.Println("Generating DSN...")

	dsn := generateDSN()

	log.Println("Opening Connection...")

	db, err := gorm.Open(mysql.New(mysql.Config{DSN: dsn}), &gorm.Config{})

	database = db

	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	log.Println("Migrating Tables...")
	migrateTables()

	return database
}

func migrateTables() {
	log.Println("Migrating Users Database...")
	err := GetDB().AutoMigrate(&models.User{})
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

func generateDSN() string {
	var dsn string

	if config.GetConfig().Database.URL != "" {
		dsn = config.GetConfig().Database.URL
	} else {
		if len(config.GetConfig().Database.ExtraOptions) == 0 {
			dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.GetConfig().Database.User, config.GetConfig().Database.Password, config.GetConfig().Database.Host, config.GetConfig().Database.Port, config.GetConfig().Database.DBName)
		} else {
			dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?", config.GetConfig().Database.User, config.GetConfig().Database.Password, config.GetConfig().Database.Host, config.GetConfig().Database.Port, config.GetConfig().Database.DBName)

			if config.GetConfig().Database.ExtraOptions["charset"] == "" {
				dsn = dsn + "charset=utf8mb4&"
			}

			if config.GetConfig().Database.ExtraOptions["parseTime"] == "" {
				dsn = dsn + "parseTime=True&"
			}

			if config.GetConfig().Database.ExtraOptions["loc"] == "" {
				dsn = dsn + "loc=Local&"
			}

			for key, value := range config.GetConfig().Database.ExtraOptions {
				dsn = dsn + "&" + key + "=" + value
			}
		}
	}

	return dsn
}

func GetDB() *gorm.DB {
	return database
}
