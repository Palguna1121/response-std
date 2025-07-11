package responses

import (
	"time"
)

// APIResponse represents a single API response
type APIResponse struct {
	RequestID  string            `json:"request_id"`
	Success    bool              `json:"success"`
	StatusCode int               `json:"status_code"`
	Message    string            `json:"message"`
	Data       interface{}       `json:"data,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Error      string            `json:"error,omitempty"`
	Duration   string            `json:"duration"`
	Timestamp  time.Time         `json:"timestamp"`
}

// BatchAPIResponse represents a batch of API responses
type BatchAPIResponse struct {
	Total     int           `json:"total"`
	Succeeded int           `json:"succeeded"`
	Failed    int           `json:"failed"`
	Success   bool          `json:"success"`
	Results   []APIResponse `json:"results"`
	Duration  string        `json:"duration"`
	Timestamp time.Time     `json:"timestamp"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
	Version     string    `json:"version"`
	Environment string    `json:"environment"`
	Uptime      string    `json:"uptime"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	Error     string    `json:"error,omitempty"`
	Code      string    `json:"code,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// ValidationErrorResponse represents a validation error response
type ValidationErrorResponse struct {
	Success   bool              `json:"success"`
	Message   string            `json:"message"`
	Errors    []ValidationError `json:"errors"`
	Timestamp time.Time         `json:"timestamp"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
	Timestamp  time.Time   `json:"timestamp"`
}

// Pagination represents pagination information
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// MetricsResponse represents API metrics
type MetricsResponse struct {
	TotalRequests   int64         `json:"total_requests"`
	SuccessfulReqs  int64         `json:"successful_requests"`
	FailedReqs      int64         `json:"failed_requests"`
	AvgResponseTime time.Duration `json:"avg_response_time"`
	Uptime          string        `json:"uptime"`
	LastRequestTime time.Time     `json:"last_request_time"`
	RequestsPerMin  float64       `json:"requests_per_minute"`
	ErrorRate       float64       `json:"error_rate"`
	Timestamp       time.Time     `json:"timestamp"`
}
