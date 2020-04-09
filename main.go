package main

import (
	"calendar/calendar"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	config := calendar.InitConfig()

	http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println()
		_, err := io.WriteString(writer, fmt.Sprintf("Hello world!"))
		if err != nil {
			log.Fatalf("/hello, %s", err)
		}
	})

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.HttpListen.Ip, config.HttpListen.Port), nil)
	if err != nil {
		log.Fatalf("Run server error, %s", err)
	}
}
