package http

type DefaultResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

type ErrorResponse struct {
	Code     int    `json:"code"`
	Status   string `json:"status"`
	Message  string `json:"message"`
	Internal error  `json:"-"`
}
