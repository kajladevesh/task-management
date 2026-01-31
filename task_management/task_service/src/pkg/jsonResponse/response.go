package pkg

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
}

func JSON(w http.ResponseWriter, statusCode int, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Status:  http.StatusText(statusCode),
		Data:    data,
		Message: message,
	}

	_ = json.NewEncoder(w).Encode(response)
}

func Success(w http.ResponseWriter, data interface{}, message string) {
	JSON(w, http.StatusOK, data, message)
}

func Created(w http.ResponseWriter, data interface{}, message string) {
	JSON(w, http.StatusCreated, data, message)
}

func Error(w http.ResponseWriter, statusCode int, message string) {
	JSON(w, statusCode, nil, message)
}
