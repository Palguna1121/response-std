package errorhandler

import (
	"net/http"

	"github.com/Palguna1121/response-std/helper"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) {
	switch err.(type) {
	case *NotFoundError:
		helper.NotFound(c, err.Error())
	case *BadRequestError:
		helper.BadRequest(c, err.Error())
	case *UnauthorizedError:
		helper.Unauthorized(c, err.Error())
	case *ForbiddenError:
		helper.Forbidden(c, err.Error())
	case *UnprocessableEntityError:
		helper.UnprocessableEntity(c, err.Error())
	case *ConflictError:
		helper.Conflict(c, err.Error())
	case *GoneError:
		helper.Error(c, http.StatusGone, err.Error())
	case *PreconditionFailedError:
		helper.Error(c, http.StatusPreconditionFailed, err.Error())
	case *RequestTimeoutError:
		helper.Error(c, http.StatusRequestTimeout, err.Error())
	case *TooManyRequestsError:
		helper.Error(c, http.StatusTooManyRequests, err.Error())
	default:
		helper.InternalServerError(c, err.Error())
	}
}

// func HandleInternalError(c *gin.Context, err error, errMsg string) bool {
// 	if err != nil {
// 		helper.InternalServerError(c, errMsg)
// 		return true
// 	}
// 	return false
// }

// func HandleBadRequest(c *gin.Context, err error, errMsg string) bool {
// 	if err != nil {
// 		helper.BadRequest(c, errMsg)
// 		return true
// 	}
// 	return false
// }

// func HandleNotFound(c *gin.Context, err error, errMsg string) bool {
// 	if err != nil {
// 		helper.NotFound(c, errMsg)
// 		return true
// 	}
// 	return false
// }

// func HandleUnauthorized(c *gin.Context, err error, errMsg string) bool {
// 	if err != nil {
// 		helper.Unauthorized(c, errMsg)
// 		return true
// 	}
// 	return false
// }

// func HandleForbidden(c *gin.Context, err error, errMsg string) bool {
// 	if err != nil {
// 		helper.Forbidden(c, errMsg)
// 		return true
// 	}
// 	return false
// }

// func HandleUnprocessableEntity(c *gin.Context, err error, errMsg string) bool {
// 	if err != nil {
// 		helper.UnprocessableEntity(c, errMsg)
// 		return true
// 	}
// 	return false
// }

// func HandleConflict(c *gin.Context, err error, errMsg string) bool {
// 	if err != nil {
// 		helper.Conflict(c, errMsg)
// 		return true
// 	}
// 	return false
// }

// func HandleGone(c *gin.Context, err error, errMsg string) bool {
// 	if err != nil {
// 		helper.Error(c, http.StatusGone, errMsg)
// 		return true
// 	}
// 	return false
// }

// func HandlePreconditionFailed(c *gin.Context, err error, errMsg string) bool {
// 	if err != nil {
// 		helper.Error(c, http.StatusPreconditionFailed, errMsg)
// 		return true
// 	}
// 	return false
// }

// func HandleRequestTimeout(c *gin.Context, err error, errMsg string) bool {
// 	if err != nil {
// 		helper.Error(c, http.StatusRequestTimeout, errMsg)
// 		return true
// 	}
// 	return false
// }

// func HandleTooManyRequests(c *gin.Context, err error, errMsg string) bool {
// 	if err != nil {
// 		helper.Error(c, http.StatusTooManyRequests, errMsg)
// 		return true
// 	}
// 	return false
// }
