package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/prasen-shakya/todo/internal/respond"
)

func GetRequestParams[T any](w http.ResponseWriter, r *http.Request) (T, error) {
	var params T

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&params); err != nil {
		respond.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
		return params, err
	}

	return params, nil
}
