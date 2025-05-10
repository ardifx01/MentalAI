package database

import (
	"chatbot/database/models"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

func ConnectDB() {
	Database, _ = gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	fmt.Println("Database Connected")

	config, _ := Database.DB()
	config.SetMaxIdleConns(10)
	config.SetMaxOpenConns(100)
	config.SetConnMaxLifetime(time.Hour)

	Database.AutoMigrate(&models.Akun{}, &models.PersonalSurvey{}, &models.Percakapan{}, &models.Omongan{})
}
