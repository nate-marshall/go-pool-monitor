package monitor

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"pool-monitor/internal/config"
	"pool-monitor/pkg/alert"
	"pool-monitor/pkg/logger"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	lastORPChange time.Time
	lastPHChange  time.Time
	lastORPValue  string
	lastPHValue   string
	lastRPMValue  int
)

// Struct to unmarshal the RPM payload if it's in JSON format
type RPMPayload struct {
	RPM int `json:"rpm"`
}

func MessageHandler(client MQTT.Client, msg MQTT.Message) {
	topic := msg.Topic()
	payload := string(msg.Payload())

	logger.Debug("Received message", "topic", topic, "payload", payload)

	switch topic {
	case config.ORPTopic:
		if payload != lastORPValue {
			lastORPChange = time.Now()
			lastORPValue = payload
		}
	case config.PHTopic:
		if payload != lastPHValue {
			lastPHChange = time.Now()
			lastPHValue = payload
		}
	case config.RPMTopic:
		logger.Debug("Received RPM message", "payload", payload) // Log the raw RPM payload
		rpmValue, err := strconv.Atoi(strings.TrimSpace(payload))
		if err != nil {
			// Try parsing as JSON if the plain int parsing fails
			var rpmPayload RPMPayload
			err = json.Unmarshal([]byte(payload), &rpmPayload)
			if err != nil {
				logger.Error("Error parsing RPM value", "error", err, "payload", payload)
				return
			}
			rpmValue = rpmPayload.RPM
		}
		lastRPMValue = rpmValue
		logger.Debug("Updated RPM value", "rpm", rpmValue)
	}
}

func CheckRPMStatus(client MQTT.Client) bool {
	logger.Debug("Subscribing to RPM topic", "topic", config.RPMTopic)
	token := client.Subscribe(config.RPMTopic, 1, nil)
	token.Wait()
	if token.Error() != nil {
		logger.Error("Error subscribing to topic", "topic", config.RPMTopic, "error", token.Error())
		return false
	}

	// Wait a moment to ensure the message handler has processed the latest value
	time.Sleep(2 * time.Second)

	// Log the current RPM value
	logger.Debug("RPM status checked", "rpm", lastRPMValue)
	return lastRPMValue > 0 // Assuming RPM is above zero if the pump is running
}

func MonitorLevels(client MQTT.Client) {
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		if !CheckRPMStatus(client) {
			logger.Warn("Pool pump is not running")
			continue
		}

		if time.Since(lastORPChange) > config.AlertInterval {
			logger.Warn("ORP level has not changed for the alert interval")
			alert.SendAlert("ORP level has not changed for the alert interval")
		}
		if time.Since(lastPHChange) > config.AlertInterval {
			logger.Warn("pH level has not changed for the alert interval")
			alert.SendAlert("pH level has not changed for the alert interval")
		}
	}
}
