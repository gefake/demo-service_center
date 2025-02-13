package helpers

import (
	"crypto/sha1"
	"fmt"
	"regexp"
)

const Salt = "askjlmn135poopaj"

func HashPasswordWithSalt(pass string) string {
	hash := sha1.New()
	hash.Write([]byte(pass))

	return fmt.Sprintf("%x", hash.Sum([]byte(Salt)))
}

func ValidatePhoneNumber(phoneNumber string) bool {
	re := regexp.MustCompile(`^\+7\d{10}$`)

	//print(re.MatchString(phoneNumber))
	return re.MatchString(phoneNumber)
}
