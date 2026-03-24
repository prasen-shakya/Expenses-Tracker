package app

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/prasen-shakya/todo/internal/auth"
	"github.com/prasen-shakya/todo/internal/controllers"
	"github.com/prasen-shakya/todo/internal/db"
	"github.com/prasen-shakya/todo/internal/expenses"
	"github.com/prasen-shakya/todo/internal/routes"
	"github.com/prasen-shakya/todo/internal/users"
)

var (
	handlerOnce sync.Once
	handler     http.Handler
	handlerErr  error
)

func Handler() (http.Handler, error) {
	handlerOnce.Do(func() {
		_ = godotenv.Load()

		dbConn, err := db.OpenPostgres()
		if err != nil {
			handlerErr = fmt.Errorf("open postgres: %w", err)
			return
		}

		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			_ = dbConn.Close()
			handlerErr = fmt.Errorf("JWT_SECRET not set")
			return
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

		handler = mux
	})

	return handler, handlerErr
}
