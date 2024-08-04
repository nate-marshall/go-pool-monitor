package alert

import (
	"bytes"
	"encoding/json"
	"net/http"
	"pool-monitor/internal/config"
	"pool-monitor/pkg/logger"
)

type MattermostPayload struct {
	Text string `json:"text"`
}

func SendAlert(message string) {
	payload := MattermostPayload{Text: message}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Error marshalling JSON payload", "error", err)
		return
	}

	resp, err := http.Post(config.MattermostWebhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		logger.Error("Error sending alert to Mattermost", "error", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("Received non-OK response from Mattermost", "status", resp.Status)
	}
}
