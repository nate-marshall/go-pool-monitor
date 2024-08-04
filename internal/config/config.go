package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

var (
	MQTTBroker           string
	MQTTUsername         string
	MQTTPassword         string
	ORPTopic             string
	PHTopic              string
	RPMTopic             string
	MattermostWebhookURL string
	AlertInterval        time.Duration
	LogLevel             string
)

func LoadConfig() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	MQTTBroker = viper.GetString("MQTT_BROKER")
	MQTTUsername = viper.GetString("MQTT_USERNAME")
	MQTTPassword = viper.GetString("MQTT_PASSWORD")
	ORPTopic = viper.GetString("ORP_TOPIC")
	PHTopic = viper.GetString("PH_TOPIC")
	RPMTopic = viper.GetString("RPM_TOPIC")
	MattermostWebhookURL = viper.GetString("MATTERMOST_WEBHOOK_URL")
	LogLevel = viper.GetString("LOG_LEVEL")

	alertIntervalMinutes := viper.GetInt("ALERT_INTERVAL_MINUTES")
	AlertInterval = time.Duration(alertIntervalMinutes) * time.Minute
}
