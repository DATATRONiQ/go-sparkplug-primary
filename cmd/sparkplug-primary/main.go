package main

import (
	"github.com/DATATRONiQ/go-sparkplug-primary/internal/server"
	"github.com/DATATRONiQ/go-sparkplug-primary/internal/sparkplug"
	"github.com/DATATRONiQ/go-sparkplug-primary/internal/util"
)

var (
	logFormat       = util.LookupEnv("LOG_FORMAT", "text")
	logFile         = util.LookupEnv("LOG_FILE", "")
	logLevel        = util.LookupEnv("LOG_LEVEL", "info")
	mqttEndpoint    = util.LookupEnv("MQTT_ENDPOINT", "tcp://localhost:1883")
	mqttClientID    = util.LookupEnv("MQTT_CLIENT_ID", "go-primary")
	mqttUsername    = util.LookupEnv("MQTT_USERNAME", "")
	mqttPassword    = util.LookupEnv("MQTT_PASSWORD", "")
	sparkplugHostID = util.LookupEnv("SPARKPLUG_HOST_ID", "go-primary")
)

func main() {

	util.InitLogger(logFormat, logFile, logLevel)

	// TODO: Make configurable
	go sparkplug.StartMQTTClient(mqttEndpoint, mqttClientID, sparkplugHostID, mqttUsername, mqttPassword)

	server.Start()
}
