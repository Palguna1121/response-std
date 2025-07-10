package handlers

import (
	"net/http"
	"time"

	"response-std/app/pkg/response"
	"response-std/libs/external/requests"
	"response-std/libs/external/services"
	"response-std/libs/responses"

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

	res := responses.HealthResponse{
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
	var apiReq requests.APIRequest

	if apiReq.URL == "" {
		response.BadRequest(c, "URL is required", nil)
		return
	}

	if err := c.ShouldBindJSON(&apiReq); err != nil {
		response.UnprocessableEntity(c, "Data tidak valid", err)
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
			response.BadRequest(c, res.Message, nil)
		case http.StatusUnauthorized:
			response.Unauthorized(c, res.Message, nil)
		case http.StatusForbidden:
			response.Forbidden(c, res.Message, nil)
		case http.StatusNotFound:
			response.NotFound(c, res.Message, nil)
		case http.StatusUnprocessableEntity:
			response.UnprocessableEntity(c, res.Message, nil)
		case http.StatusConflict:
			response.Conflict(c, res.Message, nil)
		case http.StatusInternalServerError:
			response.InternalServerError(c, res.Message, nil)
		case http.StatusServiceUnavailable:
			response.ServiceUnavailable(c, res.Message, nil)
		case http.StatusTooManyRequests:
			response.TooManyRequests(c, res.Message, nil)
		case http.StatusGone:
			response.Gone(c, res.Message, nil)
		case http.StatusPreconditionFailed:
			response.PreconditionFailed(c, res.Message, nil)
		case http.StatusRequestTimeout:
			response.RequestTimeout(c, res.Message, nil)
		default:
			response.InternalServerError(c, "Unexpected error", nil)
		}
	}
}

// template for get request
func (h *APIHandler) GetRequest(c *gin.Context) {
	var apiReq requests.APIRequest

	if err := c.ShouldBindQuery(&apiReq); err != nil {
		response.UnprocessableEntity(c, "Data tidak valid", err)
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
			response.BadRequest(c, res.Message, nil)
		case http.StatusUnauthorized:
			response.Unauthorized(c, res.Message, nil)
		default:
			response.InternalServerError(c, "Unexpected error", nil)
		}
	}
}
