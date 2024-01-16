package storage

import (
	"bookingHouses-server/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func connectToDB() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		panic("error loading .env file")
	}

	dsn := os.Getenv("DB_CONNECTION_STRING")
	db, dbError := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if dbError != nil {
		log.Panic("error connection to database")
	}

	DB = db
	return db
}

func performMigration(db *gorm.DB) {
	db.AutoMigrate(
		&models.User{},
	)
}

func InitializeDB() *gorm.DB {
	db := connectToDB()
	performMigration(db)
	return db
}
