package errors

import (
	"encoding/json"
	"net/http"
)

// Описание структуры ответа при ошибке
type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func ErrorHandler(w http.ResponseWriter, message string, status int, err error) {
	errorMessage := "Error!"
	if err != nil {
		errorMessage = err.Error()
	}
	b, _ := json.Marshal(&ErrorResponse{
		Message: message,
		Error:   errorMessage,
	})
	w.Write([]byte(b))
	w.WriteHeader(status)
}
