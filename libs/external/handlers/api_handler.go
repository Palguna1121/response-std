package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"response-std/app/pkg/response"
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

// Example usage in your service/controller
func (h *APIHandler) ExampleUsage(c *gin.Context) {
	// Example 1: Simple GET request
	resp1 := h.apiClient.Get("https://dummyjson.com/users").Execute()

	// Return one of the responses as example
	if resp1.Success {
		response.Success(c, "Request successful", resp1)
	} else {
		response.InternalServerError(c, "Request failed", errors.New(resp1.Error))
	}
}

// Generic proxy endpoint
func (h *APIHandler) ProxyRequest(c *gin.Context) {
	// Get request details from body
	var proxyReq struct {
		Method      string            `json:"method" binding:"required"`
		URL         string            `json:"url" binding:"required"`
		Headers     map[string]string `json:"headers,omitempty"`
		QueryParams map[string]string `json:"query_params,omitempty"`
		Body        interface{}       `json:"body,omitempty"`
		Timeout     int               `json:"timeout,omitempty"`
	}

	if err := c.ShouldBindJSON(&proxyReq); err != nil {
		response.UnprocessableEntity(c, "Invalid request data", err)
		return
	}

	// Build request based on method
	var builder *services.RequestBuilder

	switch proxyReq.Method {
	case "GET":
		builder = h.apiClient.Get(proxyReq.URL)
	case "POST":
		builder = h.apiClient.Post(proxyReq.URL)
	case "PUT":
		builder = h.apiClient.Put(proxyReq.URL)
	case "DELETE":
		builder = h.apiClient.Delete(proxyReq.URL)
	case "PATCH":
		builder = h.apiClient.Patch(proxyReq.URL)
	case "HEAD":
		builder = h.apiClient.Head(proxyReq.URL)
	case "OPTIONS":
		builder = h.apiClient.Options(proxyReq.URL)
	default:
		response.BadRequest(c, "Unsupported HTTP method", nil)
		return
	}

	// Add headers if provided
	if proxyReq.Headers != nil {
		builder = builder.WithHeaders(proxyReq.Headers)
	}

	// Add query parameters if provided
	if proxyReq.QueryParams != nil {
		builder = builder.WithQueryParams(proxyReq.QueryParams)
	}

	// Add body if provided
	if proxyReq.Body != nil {
		builder = builder.WithBody(proxyReq.Body)
	}

	// Set timeout if provided
	if proxyReq.Timeout > 0 {
		builder = builder.WithTimeoutSeconds(proxyReq.Timeout)
	}

	// Execute request
	res := builder.Execute()

	// Return appropriate response
	if res.Success {
		response.Success(c, "Request executed successfully", res)
	} else {
		// Map status codes to appropriate responses
		switch res.StatusCode {
		case http.StatusBadRequest:
			response.BadRequest(c, res.Message, errors.New(res.Error))
		case http.StatusUnauthorized:
			response.Unauthorized(c, res.Message, errors.New(res.Error))
		case http.StatusForbidden:
			response.Forbidden(c, res.Message, errors.New(res.Error))
		case http.StatusNotFound:
			response.NotFound(c, res.Message, errors.New(res.Error))
		case http.StatusUnprocessableEntity:
			response.UnprocessableEntity(c, res.Message, errors.New(res.Error))
		case http.StatusTooManyRequests:
			response.TooManyRequests(c, res.Message, errors.New(res.Error))
		default:
			response.InternalServerError(c, "Request failed", errors.New(res.Error))
		}
	}
}

// Service layer examples
type UserService struct {
	apiClient *services.APIClient
}

func NewUserService(apiClient *services.APIClient) *UserService {
	return &UserService{
		apiClient: apiClient,
	}
}

func (s *UserService) GetUser(userID string, token string) (*responses.APIResponse, error) {
	resp := s.apiClient.Get("https://api.example.com/users/" + userID).
		WithBearerToken(token).
		WithJSONContentType().
		Execute()

	if !resp.Success {
		return nil, fmt.Errorf("failed to get user: %s", resp.Error)
	}

	return resp, nil
}

func (s *UserService) CreateUser(userData map[string]interface{}, apiKey string) (*responses.APIResponse, error) {
	resp := s.apiClient.Post("https://api.example.com/users").
		WithAPIKey(apiKey).
		WithJSONBody(userData).
		WithTimeoutSeconds(30).
		Execute()

	if !resp.Success {
		return nil, fmt.Errorf("failed to create user: %s", resp.Error)
	}

	return resp, nil
}

func (s *UserService) UpdateUser(userID string, userData map[string]interface{}, token string) (*responses.APIResponse, error) {
	resp := s.apiClient.Put("https://api.example.com/users/" + userID).
		WithBearerToken(token).
		WithJSONBody(userData).
		Execute()

	if !resp.Success {
		return nil, fmt.Errorf("failed to update user: %s", resp.Error)
	}

	return resp, nil
}

func (s *UserService) DeleteUser(userID string, token string) (*responses.APIResponse, error) {
	resp := s.apiClient.Delete("https://api.example.com/users/" + userID).
		WithBearerToken(token).
		Execute()

	if !resp.Success {
		return nil, fmt.Errorf("failed to delete user: %s", resp.Error)
	}

	return resp, nil
}

func (s *UserService) GetMultipleUsers(userIDs []string, token string) (*responses.BatchAPIResponse, error) {
	batch := s.apiClient.Batch()

	for _, userID := range userIDs {
		batch.Add(s.apiClient.Get("https://api.example.com/users/" + userID).
			WithBearerToken(token))
	}

	resp := batch.Parallel().Execute()

	if !resp.Success {
		return nil, fmt.Errorf("batch request failed: %d out of %d requests failed", resp.Failed, resp.Total)
	}

	return resp, nil
}

// Example third-party API integrations
type IntegrationService struct {
	apiClient *services.APIClient
}

func NewIntegrationService(apiClient *services.APIClient) *IntegrationService {
	return &IntegrationService{
		apiClient: apiClient,
	}
}

// SendSlackMessage sends a message to Slack
func (s *IntegrationService) SendSlackMessage(webhookURL string, message string) error {
	payload := map[string]interface{}{
		"text": message,
	}

	resp := s.apiClient.Post(webhookURL).
		WithJSONBody(payload).
		WithTimeoutSeconds(10).
		Execute()

	if !resp.Success {
		return fmt.Errorf("failed to send Slack message: %s", resp.Error)
	}

	return nil
}

// SendEmail sends an email via SendGrid
func (s *IntegrationService) SendEmail(apiKey string, to string, subject string, content string) error {
	payload := map[string]interface{}{
		"personalizations": []map[string]interface{}{
			{
				"to": []map[string]string{
					{"email": to},
				},
			},
		},
		"from": map[string]string{
			"email": "noreply@example.com",
		},
		"subject": subject,
		"content": []map[string]string{
			{
				"type":  "text/html",
				"value": content,
			},
		},
	}

	resp := s.apiClient.Post("https://api.sendgrid.com/v3/mail/send").
		WithBearerToken(apiKey).
		WithJSONBody(payload).
		Execute()

	if !resp.Success {
		return fmt.Errorf("failed to send email: %s", resp.Error)
	}

	return nil
}

// GetWeatherData gets weather data from OpenWeatherMap
func (s *IntegrationService) GetWeatherData(apiKey string, city string) (*responses.APIResponse, error) {
	resp := s.apiClient.Get("https://api.openweathermap.org/data/2.5/weather").
		WithQuery("q", city).
		WithQuery("appid", apiKey).
		WithQuery("units", "metric").
		Execute()

	if !resp.Success {
		return nil, fmt.Errorf("failed to get weather data: %s", resp.Error)
	}

	return resp, nil
}
