package handler

import (
	"net/http"

	"github.com/prasen-shakya/todo/internal/app"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	handler, err := app.Handler()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handler.ServeHTTP(w, r)
}
