package util

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

var Responses = responses{}

type responses struct{}

var responseLogger = logrus.WithField("app", "responses")

// Return a generic success response.
func (r responses) Success(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	responseLogger.Trace("Set headers and status code")
	if _, err := w.Write([]byte(`{"status": "success"}`)); err != nil {
		responseLogger.WithError(err).Error("Failed to write generic success response")
	}
	responseLogger.WithFields(logrus.Fields{"status": http.StatusOK, "headers": w.Header()}).Trace("Sent generic success response")
}

// Return a success response with some data.
// The response data must be JSON serializable.
func (r responses) SuccessWithData(w http.ResponseWriter, data interface{}) {
	// Encode data to JSON
	encoded, err := json.Marshal(data)
	if err != nil {
		responseLogger.WithError(err).Error("Failed to encode json data")
		return
	}
	responseLogger.Trace("Encoded JSON response data")

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	responseLogger.Trace("Set headers and status code")
	if _, err := w.Write([]byte(fmt.Sprintf(`{"status": "success", "data": %s}`, string(encoded)))); err != nil {
		responseLogger.WithError(err).Error("Failed to write success with data response")
	}
	responseLogger.WithFields(logrus.Fields{"status": http.StatusOK, "headers": w.Header()}).Trace("Sent success with data response")
}

// Return an error response with a reason
func (r responses) Error(w http.ResponseWriter, status int, reason string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	responseLogger.Trace("Set headers and status code")
	if _, err := w.Write([]byte(fmt.Sprintf(`{"status": "error", "reason": "%s"}`, reason))); err != nil {
		responseLogger.WithError(err).Error("Failed to write error response")
	}
	responseLogger.WithFields(logrus.Fields{"status": status, "headers": w.Header()}).Trace("Sent error response")
}
