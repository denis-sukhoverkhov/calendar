package api_http

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func Hello(writer http.ResponseWriter, request *http.Request) {
	_, err := io.WriteString(writer, fmt.Sprintf("Hello world!"))
	if err != nil {
		log.Fatalf("/hello, %s", err)
	}
}

func HandleNotFound(writer http.ResponseWriter, request *http.Request) {
	_, err := io.WriteString(writer, fmt.Sprintf("Welcome to the main page!"))
	if err != nil {
		log.Fatalf("/, %s", err)
	}
}