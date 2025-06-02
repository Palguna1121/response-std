package helper

import (
	"fmt"
	"response-std/dto"
)

type ResponseWithData struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ResponseWithoutData struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func Response(params dto.ResponseParams) any {
	var response any
	var status string

	if params.StatusCode >= 200 && params.StatusCode <= 299 {
		status = "success"
	} else {
		status = "error"
	}

	if params.Data != nil {
		response = &ResponseWithData{
			Code:    params.StatusCode,
			Status:  status,
			Message: params.Message,
			Data:    params.Data,
		}
	} else {
		fmt.Printf("Error: %s\n", params.Message)
		response = &ResponseWithoutData{
			Code:    params.StatusCode,
			Status:  status,
			Message: params.Message,
		}
	}

	return response
}

func Success(message string, data any) any {
	return Response(dto.ResponseParams{
		StatusCode: 200,
		Message:    message,
		Data:       data,
	})
}

func Created(message string, data any) any {
	return Response(dto.ResponseParams{
		StatusCode: 201,
		Message:    message,
		Data:       data,
	})
}

func Accepted(message string, data any) any {
	return Response(dto.ResponseParams{
		StatusCode: 202,
		Message:    message,
		Data:       data,
	})
}

func NoContent(message string) any {
	return Response(dto.ResponseParams{
		StatusCode: 204,
		Message:    message,
	})
}

// Error returns a generic error response with the specified status code and message.
// It is used for cases where the error does not fit into a specific category.

func Error(code int, message string) any {
	return Response(dto.ResponseParams{
		StatusCode: code,
		Message:    message,
	})
}

func NotFound(message string) any {
	return Response(dto.ResponseParams{
		StatusCode: 404,
		Message:    message,
	})
}

func BadRequest(message string) any {
	return Response(dto.ResponseParams{
		StatusCode: 400,
		Message:    message,
	})
}

func Unauthorized(message string) any {
	return Response(dto.ResponseParams{
		StatusCode: 401,
		Message:    message,
	})
}

func Forbidden(message string) any {
	return Response(dto.ResponseParams{
		StatusCode: 403,
		Message:    message,
	})
}

func UnprocessableEntity(message string) any {
	return Response(dto.ResponseParams{
		StatusCode: 422,
		Message:    message,
	})
}

func Conflict(message string) any {
	return Response(dto.ResponseParams{
		StatusCode: 409,
		Message:    message,
	})
}

func Gone(message string) any {
	return Response(dto.ResponseParams{
		StatusCode: 410,
		Message:    message,
	})
}

func PreconditionFailed(message string) any {
	return Response(dto.ResponseParams{
		StatusCode: 412,
		Message:    message,
	})
}
func RequestTimeout(message string) any {
	return Response(dto.ResponseParams{
		StatusCode: 408,
		Message:    message,
	})
}

func TooManyRequests(message string) any {
	return Response(dto.ResponseParams{
		StatusCode: 429,
		Message:    message,
	})
}

func InternalServerError(message string) any {
	return Response(dto.ResponseParams{
		StatusCode: 500,
		Message:    message,
	})
}
