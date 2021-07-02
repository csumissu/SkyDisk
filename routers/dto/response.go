package dto

import (
	"net/http"
)

type ResponseResult struct {
	HttpStatus int         `json:"-"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Error      string      `json:"error,omitempty"`
}

func Success(data interface{}) ResponseResult {
	return ResponseResult{
		HttpStatus: http.StatusOK,
		Message:    "this operation is successful",
		Data:       data,
	}
}

func Failure(httpStatus int, message string) ResponseResult {
	return FailureWithError(httpStatus, message, nil)
}

func FailureWithError(httpStatus int, message string, err error) ResponseResult {
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}

	return ResponseResult{
		HttpStatus: httpStatus,
		Message:    message,
		Error:      errMsg,
	}
}

func ErrorResponse(err error) ResponseResult {
	return ResponseResult{
		Message: err.Error(),
	}
}
