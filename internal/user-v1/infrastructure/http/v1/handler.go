package http

import (
	"net/http"

	user "github.com/ZyoGo/default-ddd-http/internal/user-v1/core"
	"github.com/ZyoGo/default-ddd-http/internal/user-v1/infrastructure/http/v1/request"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	userSvc user.Service
}

func New(userSvc user.Service) *Handler {
	return &Handler{userSvc}
}

func (h *Handler) CreateUser(c *gin.Context) {
	reqBody := new(request.CreateUser)

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	dto := CreateUserDTO(reqBody)
	result, err := h.userSvc.CreateUser(c.Request.Context(), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "SUCCESS",
		"payload": result,
	})
}
