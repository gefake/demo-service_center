package webhookutil

import (
	"encoding/json"
	"errors"
	"telegram_bot/pkg/application"
)

func unpackTask(webhookData map[string]interface{}) (*application.ApplicationForCall, error) {
	var application *application.ApplicationForCall
	payloadData, ok := webhookData["payload"].(map[string]interface{})
	if !ok {
		return application, errors.New("invalid payload data")
	}

	appJson, err := json.Marshal(payloadData)
	if err != nil {
		return application, err
	}

	err = json.Unmarshal(appJson, &application)
	if err != nil {
		return application, err
	}

	return application, nil
}
