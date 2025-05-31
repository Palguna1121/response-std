package errorhandler

import (
	"response-std/dto"
	"response-std/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) (int, any) {
	var statusCode int

	switch err.(type) {
	case *NotFoundError:
		statusCode = http.StatusNotFound
	case *BadRequestError:
		statusCode = http.StatusBadRequest
	case *UnauthorizedError:
		statusCode = http.StatusUnauthorized
	case *ForbiddenError:
		statusCode = http.StatusForbidden
	case *UnprocessableEntityError:
		statusCode = http.StatusUnprocessableEntity
	case *ConflictError:
		statusCode = http.StatusConflict
	case *GoneError:
		statusCode = http.StatusGone
	case *PreconditionFailedError:
		statusCode = http.StatusPreconditionFailed
	case *RequestTimeoutError:
		statusCode = http.StatusRequestTimeout
	case *TooManyRequestsError:
		statusCode = http.StatusTooManyRequests
	default:
		statusCode = http.StatusInternalServerError
	}
	response := helper.Response(dto.ResponseParams{
		StatusCode: statusCode,
		Message:    err.Error(),
	})

	return statusCode, response
}
