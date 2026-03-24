package app

import (
	"fmt"
	"io/fs"
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
	webassets "github.com/prasen-shakya/todo/web"
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
		registerFrontendRoutes(mux)
		routes.RegisterAuthRoutes(mux, authController)
		routes.RegisterExpenseRoutes(mux, expensesController, authService, usersRepo)

		handler = mux
	})

	return handler, handlerErr
}

func registerFrontendRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", serveEmbeddedFile("index.html", "text/html; charset=utf-8"))
	mux.HandleFunc("GET /styles.css", serveEmbeddedFile("styles.css", "text/css; charset=utf-8"))

	jsFS, err := fs.Sub(webassets.Assets, "js")
	if err != nil {
		panic(err)
	}

	mux.Handle("GET /js/", http.StripPrefix("/js/", http.FileServer(http.FS(jsFS))))
}

func serveEmbeddedFile(path, contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)

		fileContents, err := webassets.Assets.ReadFile(path)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		_, _ = w.Write(fileContents)
	}
}
