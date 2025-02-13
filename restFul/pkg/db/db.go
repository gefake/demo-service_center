package database

import (
	"fmt"
	"os"
	"service_api/pkg/helpers"
	"service_api/pkg/logger"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID       int    `gorm:"primarykey;autoIncrement" json:"id"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	u.Password = helpers.HashPasswordWithSalt(u.Password)
	return
}

type TrustedTelegramUsers struct {
	ID         int    `gorm:"primarykey;autoIncrement" json:"id"`
	TelegramID string `json:"telegramID" bindinging:"required"`
}

type Service struct {
	ID    int    `gorm:"primarykey;autoIncrement" json:"id"`
	Name  string `json:"name" binding:"required"`
	Price int    `json:"price" binding:"required"`
}

type ApplicationForCall struct {
	ID          int    `gorm:"primarykey;autoIncrement" json:"id"`
	Name        string `json:"name" binding:"required" `
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Date        int    `json:"date"`
}

// type Employee struct {
// 	ID             int    `gorm:"primarykey;autoIncrement" json:"id"`
// 	Name           string `json:"name" binding:"required"`
// 	JobTitle       string `json:"jobTitle" binding:"required"`
// 	Characteristic string `json:"characteristic" binding:"required"`
// 	AvatarPath     string `json:"avatarPath" binding:"required"`
// }

// type Review struct {
// 	ID        int    `gorm:"primarykey;autoIncrement" json:"id"`
// 	ImagePath string `json:"imagePath" binding:"required" `
// }

func init() {
	logger.Log.Info("Connecting to database")

	if err := godotenv.Load(); err != nil {
		logger.Log.Fatal(".env файл не существует. Дальнешее подключение к БД невозможно осуществить")
	}

	dbHost, _ := os.LookupEnv("DB_HOST")
	dbPort, _ := os.LookupEnv("DB_PORT")
	dbName, _ := os.LookupEnv("DB_NAME")
	dbUser, _ := os.LookupEnv("DB_USER")
	dbPass, _ := os.LookupEnv("DB_PASSWORD")

	dataSetName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)

	//logger.Log.Info("Connect with using dataset" + dataSetName)

	if db, err := gorm.Open(mysql.Open(dataSetName), &gorm.Config{}); err != nil {
		logger.Log.Fatal("Failed to connect to database " + err.Error())
	} else {
		logger.Log.Info("Connected to database")

		db.AutoMigrate(&Service{})
		db.AutoMigrate(&User{})
		db.AutoMigrate(&ApplicationForCall{})
		db.AutoMigrate(&TrustedTelegramUsers{})
		// db.AutoMigrate(&Employee{})
		// db.AutoMigrate(&Review{})

		DataSource = &Connect{
			Context: db,
		}
	}
}
