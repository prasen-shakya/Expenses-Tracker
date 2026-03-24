package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prasen-shakya/todo/internal/app"
)

func main() {
	handler, err := app.Handler()
	if err != nil {
		log.Fatal(err)
	}

	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = "3001"
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", serverPort),
		Handler: handler,
	}

	log.Printf("Server running on 0.0.0.0:%s", serverPort)
	log.Fatal(server.ListenAndServe())
}
