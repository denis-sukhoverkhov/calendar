package main

import (
	"calendar/calendar/infrastructure"
	"log"
)

func main() {
	config := infrastructure.InitConfig()
	server, err := infrastructure.NewServer(config)
	if err != nil {
		log.Fatalf("Creating server error, %s", err)
	}
	server.Start()
}
