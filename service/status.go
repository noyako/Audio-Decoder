package service

import (
	"encoding/json"
	"net/http"
)

type ServiceError struct {
	message string
}

func ProcessForbidden(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(ServiceError{msg})
}

func ProcessServerError(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(ServiceError{msg})
}

func ProcessOk(w http.ResponseWriter, result interface{}) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(result)
}

func ProcessBadFormat(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(ServiceError{msg})
}
