package main

import (
	"calendar/calendar/infrastructure"
	"fmt"
	"go.uber.org/zap"
	"log"
)

func main() {
	config := infrastructure.InitConfig()
	listenAddr := fmt.Sprintf("%s:%d", config.HttpListen.Ip, config.HttpListen.Port)
	server, err := infrastructure.NewServer(listenAddr, zap.InfoLevel)
	if err != nil {
		log.Fatalf("Creating server error, %s", err)
	}
	server.Start()
}
