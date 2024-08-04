package main

import (
	"context"
	"os"
	"os/signal"
	"pool-monitor/internal/config"
	"pool-monitor/internal/monitor"
	"pool-monitor/pkg/logger"
	"pool-monitor/pkg/mqtt"
	"syscall"
)

func main() {
	config.LoadConfig()
	logger.Init()

	client := mqtt.Connect()
	defer client.Disconnect(250)

	ctx, cancel := context.WithCancel(context.Background())
	go monitor.MonitorLevels(ctx, client)

	// Handle system signals for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	logger.Info("Received shutdown signal")
	cancel()
}
