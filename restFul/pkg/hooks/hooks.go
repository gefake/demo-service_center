package hooks

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const WEBHOOK_URL = "http://127.0.0.1:8081/webhook"

func PostMessageToBot(postingType string, postingStruct interface{}) error {
	prepMap := make(map[string]interface{})
	prepMap["type"] = postingType
	prepMap["payload"] = postingStruct

	json, err := json.Marshal(prepMap)
	if err != nil {
		return err
	}

	http.Post(WEBHOOK_URL, "application/json", bytes.NewBuffer(json))

	return nil
}
