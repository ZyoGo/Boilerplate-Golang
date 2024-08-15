package http

import (
	"encoding/json"
	"net/http"
)

type DefaultResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func WriteResponse(w http.ResponseWriter, code int, message string, body interface{}) {
	resBody, _ := json.Marshal(DefaultResponse{
		Code:    code,
		Message: message,
		Data:    body,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resBody)
}

func WriteErrorResponse(w http.ResponseWriter, code int, message string) {
	resBody, _ := json.Marshal(DefaultResponse{
		Code:    code,
		Message: message,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resBody)
}
