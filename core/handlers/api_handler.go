package handlers

import (
	"net/http"
	"time"

	"response-std/core/models"
	"response-std/core/response"
	"response-std/core/services"

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

	res := models.HealthResponse{
		Status:      "healthy",
		Timestamp:   time.Now(),
		Version:     "1.0.0",
		Environment: "starter",
		Uptime:      uptime.String(),
	}
	response.Success(c, "API is healthy", res)
}

// Execute single API request
func (h *APIHandler) ExecuteRequest(c *gin.Context) {
	var apiReq models.APIRequest

	if apiReq.URL == "" {
		response.BadRequest(c, "URL is required")
		return
	}

	if err := c.ShouldBindJSON(&apiReq); err != nil {
		response.UnprocessableEntity(c, err.Error())
		return
	}

	// Execute the request
	res := h.apiClient.ExecuteRequest(&apiReq)

	// Return appropriate status code based on the response
	if res.Success {
		switch res.StatusCode {
		case http.StatusCreated:
			response.Created(c, "Request executed successfully", res)
		case http.StatusAccepted:
			response.Accepted(c, "Request executed successfully", res)
		case http.StatusNoContent:
			response.NoContent(c)
		default:
			response.Success(c, "Request executed successfully", res)
		}
	} else {
		switch res.StatusCode {
		case http.StatusBadRequest:
			response.BadRequest(c, res.Message)
		case http.StatusUnauthorized:
			response.Unauthorized(c, res.Message)
		case http.StatusForbidden:
			response.Forbidden(c, res.Message)
		case http.StatusNotFound:
			response.NotFound(c, res.Message)
		case http.StatusUnprocessableEntity:
			response.UnprocessableEntity(c, res.Message)
		case http.StatusConflict:
			response.Conflict(c, res.Message)
		case http.StatusInternalServerError:
			response.InternalServerError(c, res.Message)
		case http.StatusServiceUnavailable:
			response.ServiceUnavailable(c, res.Message)
		case http.StatusTooManyRequests:
			response.TooManyRequests(c, res.Message)
		case http.StatusGone:
			response.Gone(c, res.Message)
		case http.StatusPreconditionFailed:
			response.PreconditionFailed(c, res.Message)
		case http.StatusRequestTimeout:
			response.RequestTimeout(c, res.Message)
		default:
			response.InternalServerError(c, "Unexpected error")
		}
	}
}

// template for get request
func (h *APIHandler) GetRequest(c *gin.Context) {
	var apiReq models.APIRequest

	if err := c.ShouldBindQuery(&apiReq); err != nil {
		response.UnprocessableEntity(c, err.Error())
		return
	}

	// Execute the request
	res := h.apiClient.ExecuteRequest(&apiReq)

	// Return appropriate status code based on the response
	if res.Success {
		response.Success(c, "Request executed successfully", res)
	} else {
		switch res.StatusCode {
		case http.StatusBadRequest:
			response.BadRequest(c, res.Message)
		case http.StatusUnauthorized:
			response.Unauthorized(c, res.Message)
		default:
			response.InternalServerError(c, "Unexpected error")
		}
	}
}
