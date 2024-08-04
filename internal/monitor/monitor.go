package monitor

import (
	"context"
	"encoding/json"
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

// Struct to unmarshal the RPM payload
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
		var rpmPayload RPMPayload
		err := json.Unmarshal([]byte(payload), &rpmPayload)
		if err != nil {
			logger.Error("Error parsing RPM value", "error", err, "payload", payload)
			return
		}
		lastRPMValue = rpmPayload.RPM
		logger.Debug("Updated RPM value", "rpm", lastRPMValue)
	}
}

func MonitorLevels(ctx context.Context, client MQTT.Client) {
	// Subscribe to the RPM topic once and maintain the subscription
	logger.Debug("Subscribing to RPM topic", "topic", config.RPMTopic)
	token := client.Subscribe(config.RPMTopic, 1, MessageHandler)
	token.Wait()
	if token.Error() != nil {
		logger.Error("Error subscribing to topic", "topic", config.RPMTopic, "error", token.Error())
		return
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Shutting down monitor...")
			return
		case <-ticker.C:
			if lastRPMValue == 0 {
				logger.Warn("Pool pump is not running")
			} else {
				logger.Debug("RPM status checked", "rpm", lastRPMValue)
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
}
