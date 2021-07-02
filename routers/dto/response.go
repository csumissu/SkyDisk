package dto

import (
	"net/http"
)

type compatibleError struct {
	err     error
	Message string `json:"message"`
}

type Response struct {
	HttpStatus int               `json:"-"`
	Data       interface{}       `json:"data,omitempty"`
	Errors     []compatibleError `json:"errors,omitempty"`
}

func Success(data interface{}) Response {
	return Response{
		HttpStatus: http.StatusOK,
		Data:       data,
	}
}

func Failure(httpStatus int, message string) Response {
	return FailureWithCause(httpStatus, message, nil)
}

func FailureWithCause(httpStatus int, message string, cause error) Response {
	return Response{
		HttpStatus: httpStatus,
		Errors: []compatibleError{{
			Message: message,
			err:     cause,
		}},
	}
}

func ErrorResponse(err error) Response {
	return Response{
		Errors: []compatibleError{{
			Message: err.Error(),
		}},
	}
}

func EmptyResponse() Response {
	return Response{}
}
