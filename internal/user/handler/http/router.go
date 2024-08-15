package http

import "net/http"

func RegisterPath(router *http.ServeMux, h *Handler) {
	if h == nil {
		panic("item controller cannot be nil")
	}

	router.HandleFunc("POST /v1/users", h.CreateUser)
}
