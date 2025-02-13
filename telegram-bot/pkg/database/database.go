package database

import (
	"strings"
	"telegram_bot/pkg/logger"

	db "github.com/sonyarouje/simdb"
)

type TelegramAuthUser struct {
	TelegramID   int64  `json:"telegram_id"`
	TelegramName string `json:"telegram_name"`
}

func (c TelegramAuthUser) ID() (jsonField string, value interface{}) {
	value = c.TelegramName
	jsonField = "telegram_name"
	return
}

var driver *db.Driver

func IsAuthorizedUser(TelegramID int64) bool {
	err := driver.Open(TelegramAuthUser{}).Where("telegram_id", "=", TelegramID).First().AsEntity(TelegramAuthUser{})

	if err != nil && strings.Contains(err.Error(), "not found") {
		return false
	}

	return true
}

func IsAuthorizedUserByName(TelegramName string) bool {
	err := driver.Open(TelegramAuthUser{}).Where("telegram_name", "=", TelegramName).First().AsEntity(TelegramAuthUser{})

	if err != nil && strings.Contains(err.Error(), "not found") {
		return false
	}

	return true
}

// TODO: Создать метод наличия пользователя в бд
func GetUser(TelegramID int64) (*TelegramAuthUser, error) {
	var newUser *TelegramAuthUser

	err := driver.Open(newUser).Where("telegram_id", "=", TelegramID).First().AsEntity(&newUser)

	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

func GetUserByName(TelegramName string) (*TelegramAuthUser, error) {
	var newUser *TelegramAuthUser

	err := driver.Open(newUser).Where("telegram_name", "=", TelegramName).First().AsEntity(&newUser)

	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

func GetUsers() ([]TelegramAuthUser, error) {
	var users []TelegramAuthUser

	err := driver.Open(TelegramAuthUser{}).Get().AsEntity(&users)

	return users, err
}

func InsertUser(user *TelegramAuthUser) error {
	err := driver.Insert(user)

	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

func DeleteUser(user *TelegramAuthUser) error {
	err := driver.Delete(user)

	return err
}

func init() {
	d, err := db.New("data")

	if err != nil {
		logger.Log.Error(err.Error())
	}

	driver = d

	// newUser := TelegramAuthUser{
	// 	TelegramID:   1244,
	// 	TelegramName: "userMsg.From.UserName",
	// }

	// fmt.Println("User was posting into db")
	// InsertUser(&newUser)
}
