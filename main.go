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
	"github.com/prasen-shakya/todo/internal/expenses"
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

	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		log.Fatal("JWT_SECRET not set")
	}

	usersRepo := users.NewRepository(dbConn)
	expensesRepo := expenses.NewRepository(dbConn)

	authService := auth.NewService(usersRepo, []byte(jwtSecret))
	expensesService := expenses.NewService(expensesRepo)

	authController := controllers.NewAuthController(authService)
	expensesController := controllers.NewExpenseController(expensesService)

	mux := http.NewServeMux()
	routes.RegisterAuthRoutes(mux, authController)
	routes.RegisterExpenseRoutes(mux, expensesController, authService, usersRepo)

	serverPort := os.Getenv("PORT")

	if serverPort == "" {
		log.Fatal("Server port not set")
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", serverPort),
		Handler: mux,
	}

	// Run server in a goroutine
	go func() {
		log.Printf("Server running on 0.0.0.0:%s", serverPort)
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
