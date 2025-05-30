package utils

import (
	"encoding/json"
	"net/http"
)

// Helper function to standardize JSON responses
func RespondWithJSON(w http.ResponseWriter, statusCode int, response any) {
	if statusCode < 100 || statusCode > 599 {
		// Fallback to 500 if an invalid status code is passed
		statusCode = http.StatusInternalServerError
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)           // Call once
	json.NewEncoder(w).Encode(response) // Do not call WriteHeader again
}
