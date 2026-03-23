package controllers

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func GetRequestParams[T any](w http.ResponseWriter, r *http.Request) (T, error) {
	var params T

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&params); err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
		return params, err
	}

	return params, nil
}
