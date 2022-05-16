package main

import (
	"github.com/DATATRONiQ/go-sparkplug-primary/internal/server"
	"github.com/DATATRONiQ/go-sparkplug-primary/internal/sparkplug"
)

func main() {

	// TODO: Make configurable
	go sparkplug.StartMQTTClient("tcp://localhost:1883", "go-primary", "go-primary", "", "")

	server.Start()
}
