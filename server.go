package main

import (
	"fmt"
	"net/http"
	"os"

	ph "github.com/messageparser/http"
)

func main() {

	// TODO: Set up make file
	port := os.Getenv("SERVER_PORT")

	if "" == port {
		fmt.Println("Port not provided will default to 8080.")
		port = "8080"
	}

	// TODO: Load logging dependency
	// TODO: Parse env variables
	// TODO: Create tests.
	// TODO: Update Readme.

	http.HandleFunc("/v1/message", ph.ParseMessageHandler)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
