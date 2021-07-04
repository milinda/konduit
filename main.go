package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func createLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()

	if err != nil {
		log.Panic("Cannot initialize logger.", err)
	}

	return logger
}


func main() {
	// Setup logging
	var logger = createLogger()
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	hub := newHub()
	go hub.start()

	go startWebApp(hub)

	startMqtt(hub, "", "", "")

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
}
