package routes

import (
	"net/http"

	"github.com/prasen-shakya/todo/internal/controllers"
)

func RegisterAuthRoutes(mux *http.ServeMux, authController *controllers.AuthController) {
	mux.HandleFunc("POST /register", authController.Register)
	mux.HandleFunc("POST /login", authController.Login)
}
