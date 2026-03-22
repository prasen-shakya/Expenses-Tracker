package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/prasen-shakya/todo/internal/auth"
	"github.com/prasen-shakya/todo/internal/controllers"
	"github.com/prasen-shakya/todo/internal/db"
	"github.com/prasen-shakya/todo/internal/routes"
	"github.com/prasen-shakya/todo/internal/users"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbConn, err := db.OpenPostgres()
	if err != nil {
		log.Fatalf("Could not connect to postgres: %v", err)
	}
	defer dbConn.Close()

	usersRepo := users.NewRepository(dbConn)
	authService := auth.NewService(usersRepo)
	authController := controllers.NewAuthController(authService)

	mux := http.NewServeMux()
	routes.RegisterAuthRoutes(mux, authController)

	serverPort := os.Getenv("PORT")

	if serverPort == "" {
		log.Fatal("Server port not set")
	}

	server := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", serverPort),
		Handler: mux,
	}

	// Run server in a goroutine
	go func() {
		log.Printf("Server running on http://localhost:%s", serverPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Listen for shutdown signal (Ctrl+C)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop // block until signal received
	log.Println("Shutting down server...")

	// Give server time to finish requests
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Forced shutdown: %v", err)
	}

	log.Println("Server exited cleanly")
}
