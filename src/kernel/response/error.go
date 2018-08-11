package response

import (
	"net/http"
)

type StatusError struct {
	Message string `json:"_message"`
	Code int `json:"_code"`
	Meta map[string]interface{} `json:"_meta"`
}

func NotFound(n string) *StatusError {
	return &StatusError{
		Message: "Resource not found",
		Code:    http.StatusNotFound,
		Meta: map[string]interface{}{
			"_notice": n,
		},
	}
}

func Unauthorized(n string) *StatusError {
	return &StatusError{
		Message: "Unauthorized",
		Code:    401,
		Meta: map[string]interface{}{
			"_notice": n,
		},
	}
}


func TooManyRequests(n string) *StatusError {
	return &StatusError{
		Message: "Too Many Requests",
		Code:    429,
		Meta: map[string]interface{}{
			"_notice": n,
		},
	}
}

func Unprocessable(n string) *StatusError {
	return &StatusError{
		Message: "Entity unprocessable",
		Code:    422,
		Meta: map[string]interface{}{
			"_notice": n,
		},
	}
}

func InternalError(n string) *StatusError {
	return &StatusError{
		Message: "Internal Server Error",
		Code:    500,
		Meta: map[string]interface{}{
			"_notice": n,
		},
	}
}


func BadRequest(n string) *StatusError {
	return &StatusError{
		Message: "Bad request",
		Code:    400,
		Meta: map[string]interface{}{
			"_notice": n,
		},
	}
}

func Forbidden(n string) *StatusError {
	return &StatusError{
		Message: "Forbidden",
		Code:    403,
		Meta: map[string]interface{}{
			"_notice": n,
		},
	}
}
