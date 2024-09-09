package http

import (
	"net/http"

	user "github.com/ZyoGo/default-ddd-http/internal/user/core"
	"github.com/ZyoGo/default-ddd-http/internal/user/infrastructure/http/v1/request"
	"github.com/ZyoGo/default-ddd-http/internal/user/infrastructure/http/v1/response"
	commonHTTP "github.com/ZyoGo/default-ddd-http/pkg/http"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	userSvc user.Service
}

func New(userSvc user.Service) *Handler {
	return &Handler{userSvc}
}

func (h *Handler) SignUp(c *gin.Context) {
	reqBody := new(request.SignUp)

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	dto := SignUpDTO(reqBody)
	result, err := h.userSvc.SignUp(c.Request.Context(), dto)
	if err != nil {
		errResp := commonHTTP.RenderErrResp(err)
		c.JSON(errResp.Code, errResp)
		return
	}

	resp := response.SignUpResp(result.ID)
	c.JSON(http.StatusCreated, resp)
}
