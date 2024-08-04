package mqtt

import (
	"pool-monitor/internal/config"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func Connect() MQTT.Client {
	opts := MQTT.NewClientOptions().
		AddBroker(config.MQTTBroker).
		SetClientID("pool-monitor").
		SetUsername(config.MQTTUsername).
		SetPassword(config.MQTTPassword).
		SetDefaultPublishHandler(nil)

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return client
}
