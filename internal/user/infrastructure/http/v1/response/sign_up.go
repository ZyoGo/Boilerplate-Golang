package response

import (
	"net/http"

	common "github.com/ZyoGo/default-ddd-http/pkg/http"
)

type SignUp struct {
	ID string `json:"id"`
}

func SignUpResp(id string) common.DefaultResponse {
	payload := SignUp{
		ID: id,
	}

	return common.DefaultResponse{
		Code:    http.StatusCreated,
		Message: "CREATED",
		Payload: payload,
	}
}
