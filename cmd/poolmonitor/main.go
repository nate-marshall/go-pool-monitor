package main

import (
	"pool-monitor/internal/config"
	"pool-monitor/internal/monitor"
	"pool-monitor/pkg/logger"
	"pool-monitor/pkg/mqtt"
)

func main() {
	config.LoadConfig()
	logger.Init()

	client := mqtt.Connect()
	defer client.Disconnect(250)

	client.Subscribe(config.ORPTopic, 1, monitor.MessageHandler)
	client.Subscribe(config.PHTopic, 1, monitor.MessageHandler)

	logger.Info("Starting to monitor levels...")
	monitor.MonitorLevels(client)
}
