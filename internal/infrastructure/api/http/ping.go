package infrastructure

import (
	"fmt"
	"io"
	"net/http"
)

func Hello(writer http.ResponseWriter, _ *http.Request) {
	_, err := io.WriteString(writer, fmt.Sprintf("Hello world!"))
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func Main(writer http.ResponseWriter, _ *http.Request) {
	_, err := io.WriteString(writer, fmt.Sprintf("Welcome to the main page!"))
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}
