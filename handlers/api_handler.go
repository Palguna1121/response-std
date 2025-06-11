package handlers

import (
	"net/http"
	"time"

	"github.com/Palguna1121/response-std/helper"
	"github.com/Palguna1121/response-std/models"
	"github.com/Palguna1121/response-std/services"

	"github.com/gin-gonic/gin"
)

type APIHandler struct {
	apiClient *services.APIClient
	logger    *services.Logger
	startTime time.Time
}

func NewAPIHandler(apiClient *services.APIClient, logger *services.Logger) *APIHandler {
	return &APIHandler{
		apiClient: apiClient,
		logger:    logger,
		startTime: time.Now(),
	}
}

// Health check endpoint
func (h *APIHandler) HealthCheck(c *gin.Context) {
	uptime := time.Since(h.startTime)

	response := models.HealthResponse{
		Status:      "healthy",
		Timestamp:   time.Now(),
		Version:     "1.0.0",
		Environment: "starter",
		Uptime:      uptime.String(),
	}
	helper.Success(c, "API is healthy", response)
}

// Execute single API request
func (h *APIHandler) ExecuteRequest(c *gin.Context) {
	var apiReq models.APIRequest

	if err := c.ShouldBindJSON(&apiReq); err != nil {
		helper.UnprocessableEntity(c, err.Error())
		return
	}

	// Execute the request
	response := h.apiClient.ExecuteRequest(&apiReq)

	// Return appropriate status code based on the response
	if response.Success {
		switch response.StatusCode {
		case http.StatusCreated:
			helper.Created(c, "Request executed successfully", response)
		case http.StatusAccepted:
			helper.Accepted(c, "Request executed successfully", response)
		case http.StatusNoContent:
			helper.NoContent(c)
		default:
			helper.Success(c, "Request executed successfully", response)
		}
	} else {
		switch response.StatusCode {
		case http.StatusBadRequest:
			helper.BadRequest(c, response.Message)
		case http.StatusUnauthorized:
			helper.Unauthorized(c, response.Message)
		case http.StatusForbidden:
			helper.Forbidden(c, response.Message)
		case http.StatusNotFound:
			helper.NotFound(c, response.Message)
		case http.StatusUnprocessableEntity:
			helper.UnprocessableEntity(c, response.Message)
		case http.StatusConflict:
			helper.Conflict(c, response.Message)
		case http.StatusInternalServerError:
			helper.InternalServerError(c, response.Message)
		case http.StatusServiceUnavailable:
			helper.ServiceUnavailable(c, response.Message)
		case http.StatusTooManyRequests:
			helper.TooManyRequests(c, response.Message)
		case http.StatusGone:
			helper.Gone(c, response.Message)
		case http.StatusPreconditionFailed:
			helper.PreconditionFailed(c, response.Message)
		case http.StatusRequestTimeout:
			helper.RequestTimeout(c, response.Message)
		default:
			helper.InternalServerError(c, "Unexpected error")
		}
	}
}

// template for get request
func (h *APIHandler) GetRequest(c *gin.Context) {
	var apiReq models.APIRequest

	if err := c.ShouldBindQuery(&apiReq); err != nil {
		helper.UnprocessableEntity(c, err.Error())
		return
	}

	// Execute the request
	response := h.apiClient.ExecuteRequest(&apiReq)

	// Return appropriate status code based on the response
	if response.Success {
		helper.Success(c, "Request executed successfully", response)
	} else {
		switch response.StatusCode {
		case http.StatusBadRequest:
			helper.BadRequest(c, response.Message)
		case http.StatusUnauthorized:
			helper.Unauthorized(c, response.Message)
		default:
			helper.InternalServerError(c, "Unexpected error")
		}
	}
}

// template for post request
func (h *APIHandler) PostRequest(c *gin.Context) {
	var apiReq models.APIRequest

	if err := c.ShouldBindJSON(&apiReq); err != nil {
		helper.UnprocessableEntity(c, err.Error())
		return
	}

	// Execute the request
	response := h.apiClient.ExecuteRequest(&apiReq)

	// Return appropriate status code based on the response
	if response.Success {
		helper.Success(c, "Request executed successfully", response)
	} else {
		switch response.StatusCode {
		case http.StatusBadRequest:
			helper.BadRequest(c, response.Message)
		case http.StatusUnauthorized:
			helper.Unauthorized(c, response.Message)
		default:
			helper.InternalServerError(c, "Unexpected error")
		}
	}
}

// template for put request
func (h *APIHandler) PutRequest(c *gin.Context) {
	var apiReq models.APIRequest

	if err := c.ShouldBindJSON(&apiReq); err != nil {
		helper.UnprocessableEntity(c, err.Error())
		return
	}
	// Execute the request
	response := h.apiClient.ExecuteRequest(&apiReq)

	// Return appropriate status code based on the response
	if response.Success {
		helper.Success(c, "Request executed successfully", response)
	} else {
		switch response.StatusCode {
		case http.StatusBadRequest:
			helper.BadRequest(c, response.Message)
		case http.StatusUnauthorized:
			helper.Unauthorized(c, response.Message)
		default:
			helper.InternalServerError(c, "Unexpected error")
		}
	}
}

// template for delete request
func (h *APIHandler) DeleteRequest(c *gin.Context) {
	var apiReq models.APIRequest

	if err := c.ShouldBindJSON(&apiReq); err != nil {
		helper.UnprocessableEntity(c, err.Error())
		return
	}
	// Execute the request
	response := h.apiClient.ExecuteRequest(&apiReq)
	// Return appropriate status code based on the response
	if response.Success {
		helper.NoContent(c)
	} else {
		switch response.StatusCode {
		case http.StatusBadRequest:
			helper.BadRequest(c, response.Message)
		case http.StatusUnauthorized:
			helper.Unauthorized(c, response.Message)
		default:
			helper.InternalServerError(c, "Unexpected error")
		}
	}
}

// template for patch request
func (h *APIHandler) PatchRequest(c *gin.Context) {
	var apiReq models.APIRequest

	if err := c.ShouldBindJSON(&apiReq); err != nil {
		helper.UnprocessableEntity(c, err.Error())
		return
	}
	// Execute the request
	response := h.apiClient.ExecuteRequest(&apiReq)
	// Return appropriate status code based on the response
	if response.Success {
		helper.Success(c, "Request executed successfully", response)
	} else {
		switch response.StatusCode {
		case http.StatusBadRequest:
			helper.BadRequest(c, response.Message)
		case http.StatusUnauthorized:
			helper.Unauthorized(c, response.Message)
		default:
			helper.InternalServerError(c, "Unexpected error")
		}
	}
}

// template for head request
func (h *APIHandler) HeadRequest(c *gin.Context) {
	var apiReq models.APIRequest

	if err := c.ShouldBindQuery(&apiReq); err != nil {
		helper.UnprocessableEntity(c, err.Error())
		return
	}

	// Execute the request
	response := h.apiClient.ExecuteRequest(&apiReq)
	// Return appropriate status code based on the response
	if response.Success {
		c.Header("Content-Type", "application/json")
		c.Status(response.StatusCode)
	} else {
		switch response.StatusCode {
		case http.StatusBadRequest:
			helper.BadRequest(c, response.Message)
		case http.StatusUnauthorized:
			helper.Unauthorized(c, response.Message)
		default:
			helper.InternalServerError(c, "Unexpected error")
		}
	}
}

// template for options request
func (h *APIHandler) OptionsRequest(c *gin.Context) {
	var apiReq models.APIRequest

	if err := c.ShouldBindQuery(&apiReq); err != nil {
		helper.UnprocessableEntity(c, err.Error())
		return
	}
	// Execute the request
	response := h.apiClient.ExecuteRequest(&apiReq)
	// Return appropriate status code based on the response
	if response.Success {
		c.Header("Content-Type", "application/json")
		c.Status(response.StatusCode)
	} else {
		switch response.StatusCode {
		case http.StatusBadRequest:
			helper.BadRequest(c, response.Message)
		case http.StatusUnauthorized:
			helper.Unauthorized(c, response.Message)
		default:
			helper.InternalServerError(c, "Unexpected error")
		}
	}
}
