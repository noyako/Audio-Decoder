package service

import (
	"encoding/json"
	"net/http"
)

// Error represents general server error
type Error struct {
	Message string `json:"error"`
}

// ProcessForbidden returns forbidden response
func ProcessForbidden(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(Error{msg})
}

// ProcessServerError returns service error response
func ProcessServerError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(Error{msg})
}

// ProcessOk returns service ok 200 response with json body
func ProcessOk(w http.ResponseWriter, result interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// ProcessOkString returns service ok 200 response with string body
func ProcessOkString(w http.ResponseWriter, result []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

// ProcessOkFile returns service ok 200 response
func ProcessOkFile(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

// ProcessBadFormat returns bad request error
func ProcessBadFormat(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(Error{msg})
}
