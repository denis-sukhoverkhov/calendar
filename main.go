package main

import (
	"github.com/denis-sukhoverkhov/calendar/internal/infrastructure"
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
