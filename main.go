package main

import (
	"calendar/calendar"
	"fmt"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
)

func hello(writer http.ResponseWriter, request *http.Request) {
	_, err := io.WriteString(writer, fmt.Sprintf("Hello world!"))
	if err != nil {
		log.Fatalf("/hello, %s", err)
	}
}

func handleNotFound(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.NotFound(writer, request)
		return
	}

	_, err := io.WriteString(writer, fmt.Sprintf("Welcome to the main page!"))
	if err != nil {
		log.Fatalf("/, %s", err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/hello", calendar.LoggerMiddleware(http.HandlerFunc(hello)))
	mux.Handle("/", calendar.LoggerMiddleware(http.HandlerFunc(handleNotFound)))

	config := calendar.InitConfig()
	listenAddr := fmt.Sprintf("%s:%d", config.HttpListen.Ip, config.HttpListen.Port)
	server, err := calendar.NewServer(listenAddr, zap.InfoLevel, mux)
	if err != nil {
		log.Fatalf("Creating server error, %s", err)
	}
	server.Start()
}
