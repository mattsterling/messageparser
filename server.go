package main

import (
	"os"
	"fmt"
	"net/http"
	"github.com/messageparser/http"
)

func main() {

	// TODO: Set up make file
	port := os.Getenv("SERVER_PORT")

	if "" == port {
		fmt.Println("Port not provided will default to 8080.")
		port = "8080"
	}

	fmt.Println("Symbols:", byte('@'), byte('('),  byte(')'), byte(' '), byte('h'))

	// TODO: Load logging dependency
	// TODO: Parse env variables

	// TODO: Create parser.
	// TODO: Create HTTP client to get page title of any links found in the text.
	// TODO: Create tests.
	// TODO: Update Readme.

	http.HandleFunc("/v1/message", message_parser.ParseMessageHandler)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
