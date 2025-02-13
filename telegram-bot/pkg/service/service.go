package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"telegram_bot/pkg/logger"
	"time"

	"github.com/joho/godotenv"
)

var JWTToken string

type TrustedTelegramUsers struct {
	ID         int    `json:"id"`
	TelegramID string `json:"telegramID"`
}

var TrustedUsers []*TrustedTelegramUsers

const AUTH_URL = "http://127.0.0.1:8080/auth/admin/sign-in/"

func AddTrustedUser(user *TrustedTelegramUsers) {
	TrustedUsers = append(TrustedUsers, user)
}

func removeByIndex(array []*TrustedTelegramUsers, index int) []*TrustedTelegramUsers {
	return append(array[:index], array[index+1:]...)
}

func RemoveTrustedUser(user *TrustedTelegramUsers) bool {
	for i, v := range TrustedUsers {
		fmt.Printf("%d == %d", user.ID, v.ID)
		if user.ID == v.ID {
			// Удаление элемента из среза
			TrustedUsers = removeByIndex(TrustedUsers, i)

			tets := TrustedUsers

			fmt.Println(tets)

			return true
		}
	}

	return false
}

func IsTrustedUsername(username string) bool {
	for _, user := range TrustedUsers {
		if !strings.Contains(user.TelegramID, username) {
			continue
		}

		return true
	}

	return false
}

// Получает JWT токен у основного сервиса
func getJWTToken() error {
	if err := godotenv.Load(); err != nil {
		logger.Log.Error(".env файл не существует")
	}

	client := http.Client{
		Timeout: time.Second * 10,
	}

	r, err := client.Get(AUTH_URL)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	fmt.Println("Сервер доступен на порту 8080")

	login := os.Getenv("JWT_LOGIN")
	password := os.Getenv("JWT_PASS")

	data := map[string]string{
		"name":     login,
		"password": password,
	}
	payload, err := json.Marshal(data)

	if err != nil {
		return err
	}

	resp, err := http.Post(AUTH_URL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var tokenResponse map[string]string
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return err
	}

	fmt.Println(data)

	token, ok := tokenResponse["token"]
	if !ok {
		return errors.New("token not found")
	}

	JWTToken = token

	return nil
}

func getTrustedUsers() {
	request, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/api/admin/telegram-trust", nil)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	//fmt.Println("header: " + JWTToken)
	request.Header.Add("Authorization", "Bearer "+JWTToken)
	client := &http.Client{}
	r, err := client.Do(request)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&TrustedUsers)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	// body, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	logger.Log.Error(err.Error())
	// }

	// fmt.Println(string(body))

	//fmt.Println(TrustedUsers)
}

func init() {
	if err := getJWTToken(); err != nil {
		logger.Log.Error(err.Error())
		return
	}

	getTrustedUsers()
}
