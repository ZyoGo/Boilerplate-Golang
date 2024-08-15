package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/ZyoGo/default-ddd-http/internal/user"
	"github.com/ZyoGo/default-ddd-http/internal/user/handler/http/request"
	helperLib "github.com/ZyoGo/default-ddd-http/pkg/http"
)

type Handler struct {
	userSvc user.Service
}

func New(userSvc user.Service) *Handler {
	return &Handler{userSvc}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3000*time.Millisecond)
	defer cancel()

	reqBody := new(request.CreateUser)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		helperLib.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := json.Unmarshal(body, reqBody); err != nil {
		helperLib.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.userSvc.CreateUser(ctx, user.User{
		Email:    reqBody.Email,
		Password: reqBody.Password,
	})
	if err != nil {
		helperLib.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	helperLib.WriteResponse(w, http.StatusCreated, "Success", result)
}
