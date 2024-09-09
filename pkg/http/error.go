package http

import (
	"errors"
	"net/http"

	"github.com/ZyoGo/default-ddd-http/pkg/derrors"
	stdLog "github.com/ZyoGo/default-ddd-http/pkg/logger"
)

func RenderErrResp(err error) (resp ErrorResponse) {
	resp = ErrorResponse{
		Code:     http.StatusInternalServerError,
		Status:   "SERVER_ERROR",
		Message:  "internal server error",
		Internal: err,
	}

	var ierr *derrors.Error
	if !errors.As(err, &ierr) {
		stdLog.Get().Debug().Str("internal server error", err.Error())
		return resp
	} else {
		stdLog.Get().Debug().Str("internal server error", err.Error())
		cases := map[derrors.ErrorCode]ErrorResponse{
			derrors.ErrorCodeNotFound:          {Status: "BAD_REQUEST", Code: http.StatusBadRequest, Message: "Data Not Found", Internal: ierr},
			derrors.ErrorCodeCustomBadRequest:  {Status: "BAD_REQUEST", Code: http.StatusBadRequest, Message: ierr.Error(), Internal: ierr},
			derrors.ErrorCodeAlreadyRegistered: {Status: "ALREADY_EXISTS", Code: http.StatusConflict, Message: ierr.Error(), Internal: ierr},
		}

		action, exists := cases[ierr.Code()]
		if exists {
			return action
		}

		return resp
	}
}
