package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/prasen-shakya/todo/internal/auth"
	"github.com/prasen-shakya/todo/internal/users"
)

type AuthParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthController struct {
	authService *auth.Service
}

func NewAuthController(authService *auth.Service) *AuthController {
	return &AuthController{authService: authService}
}

func getAuthParams(w http.ResponseWriter, r *http.Request) (AuthParams, error) {
	var params AuthParams

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&params); err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
		return AuthParams{}, err
	}

	return params, nil
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	params, err := getAuthParams(w, r)
	if err != nil {
		return
	}

	user, err := c.authService.Login(r.Context(), params.Username, params.Password)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Could not log in"

		if errors.Is(err, auth.ErrInvalidCredentials) {
			status = http.StatusUnauthorized
			message = "Invalid credentials"
		}

		WriteJSON(w, status, map[string]string{"error": message})
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{
		"message":  "Login successful",
		"username": user.Username,
	})
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	params, err := getAuthParams(w, r)
	if err != nil {
		return
	}

	user, err := c.authService.Register(r.Context(), params.Username, params.Password)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Could not register user"

		switch {
		case errors.Is(err, auth.ErrInvalidUsername), errors.Is(err, auth.ErrInvalidPassword):
			status = http.StatusBadRequest
			message = err.Error()
		case errors.Is(err, users.ErrUsernameTaken):
			status = http.StatusConflict
			message = "Username already exists"
		}

		WriteJSON(w, status, map[string]string{"error": message})
		return
	}

	WriteJSON(w, http.StatusCreated, map[string]any{
		"message": "User registered",
		"user": map[string]any{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}
