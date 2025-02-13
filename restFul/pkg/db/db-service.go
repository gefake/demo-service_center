package database

import (
	"os"
	"service_api/pkg/helpers"

	"gorm.io/gorm"
)

var DataSource *Connect

type Connect struct {
	Context *gorm.DB
}

// Получает пользователя из БД по паролю и юзернейму
func (s *Connect) GetUser(username, pass string) (User, error) {
	var user User

	err := s.Context.Where("name = ? AND password = ?", username, helpers.HashPasswordWithSalt(pass)).First(&user).Error

	return user, err
}

func (s *Connect) InitAdmin() {
	var existingUser User
	adminName, _ := os.LookupEnv("ADMIN_NAME")
	adminPass, _ := os.LookupEnv("ADMIN_PASS")
	result := s.Context.Where("name = ?", adminName).First(&existingUser)

	if result.RowsAffected == 0 {
		DataSource.Context.Create(&User{
			Name:     adminName,
			Password: adminPass,
		})
	} else {
		if !comparePasswords(existingUser.Password, adminPass) {
			hashedPassword := helpers.HashPasswordWithSalt(adminPass)
			s.Context.Model(&existingUser).Update("Password", hashedPassword)
		}
	}
}

func comparePasswords(hashedPassword, plainPassword string) bool {
	return hashedPassword == helpers.HashPasswordWithSalt(plainPassword)
}
