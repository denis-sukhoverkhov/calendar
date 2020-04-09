package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {

	http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println()
		io.WriteString(writer, fmt.Sprintf("Hello world!"))
	})

	http.ListenAndServe(":8080", nil)
}
