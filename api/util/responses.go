package util

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var Responses = responses{}

type responses struct{}

// Return a generic success response.
func (r responses) Success(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`{"status": "success"}`)); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

// Return a success response with some data.
// The response data must be JSON serializable.
func (r responses) SuccessWithData(w http.ResponseWriter, data interface{}) {
	// Encode data to JSON
	encoded, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to encode data to JSON: %v", err)
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(fmt.Sprintf(`{"status": "success", "data": %s}`, string(encoded)))); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

// Return an error response with a reason
func (r responses) Error(w http.ResponseWriter, status int, reason string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write([]byte(fmt.Sprintf(`{"status": "error", "reason": "%s"}`, reason))); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
