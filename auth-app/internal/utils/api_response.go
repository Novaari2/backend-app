package utils

type response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ApiResponse(status int, message string, data interface{}) response {
	return response{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
