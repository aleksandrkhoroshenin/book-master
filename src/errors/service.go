package errors

import (
	"encoding/json"
	"net/http"
)

// Описание структуры ответа при ошибке
type ErrorResponse struct {
	Message string `json:"message"`
	Error   error  `json:"error"`
}

func ErrorHandler(w http.ResponseWriter, message string, status int, err error) {
	b, _ := json.Marshal(&ErrorResponse{
		Message: message,
		Error:   err,
	})
	w.Write([]byte(b))
	w.WriteHeader(status)
}
