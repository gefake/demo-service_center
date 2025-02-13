package webhookutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"telegram_bot/pkg/database"
	"telegram_bot/pkg/service"
)

func unpackAddTrustUser(webhookData map[string]interface{}) (*service.TrustedTelegramUsers, error) {
	var user *service.TrustedTelegramUsers
	payloadData, ok := webhookData["payload"].(map[string]interface{})
	if !ok {
		return user, errors.New("invalid payload data")
	}

	appJson, err := json.Marshal(payloadData)
	if err != nil {
		return user, err
	}

	err = json.Unmarshal(appJson, &user)
	if err != nil {
		return user, err
	}

	name := strings.ReplaceAll(user.TelegramID, "@", "")

	user.TelegramID = name

	_ = append(service.TrustedUsers, user)

	return user, nil
}

func unpackRemoveTrustUser(webhookData map[string]interface{}) error {
	var user *service.TrustedTelegramUsers
	payloadData, ok := webhookData["payload"].(map[string]interface{})
	if !ok {
		return errors.New("invalid payload data")
	}

	appJson, err := json.Marshal(payloadData)
	if err != nil {
		return err
	}

	err = json.Unmarshal(appJson, &user)
	if err != nil {
		return err
	}

	name := strings.ReplaceAll(user.TelegramID, "@", "")

	user.TelegramID = name

	hasDeleted := service.RemoveTrustedUser(user)
	isAuth := database.IsAuthorizedUserByName(user.TelegramID)

	fmt.Println(service.TrustedUsers)
	if hasDeleted && isAuth {
		user, err := database.GetUserByName(user.TelegramID)

		if err != nil {
			return err
		}

		database.DeleteUser(user)
	}

	trusrz := service.TrustedUsers
	fmt.Println("trusted uzrs")
	fmt.Println(trusrz)

	return nil
}
