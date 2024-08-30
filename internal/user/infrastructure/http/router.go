package http

import (
	"log"

	v1 "github.com/ZyoGo/default-ddd-http/internal/user/infrastructure/http/v1"
	"github.com/gin-gonic/gin"
)

func RegisterPath(router *gin.Engine, hV1 *v1.Handler) {
	if hV1 == nil {
		log.Fatal("handler v1 cannot be nil")
	}

	router.POST("v1/sign-up", hV1.CreateUser)
}
