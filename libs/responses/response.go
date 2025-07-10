package responses

import (
	"time"
)

type APIResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data,omitempty"`
	Message    string      `json:"message,omitempty"`
	Error      string      `json:"error,omitempty"`
	StatusCode int         `json:"status_code"`
	RequestID  string      `json:"request_id,omitempty"`
	Timestamp  time.Time   `json:"timestamp"`
	Duration   string      `json:"duration,omitempty"`
}

type BatchAPIResponse struct {
	Success   bool          `json:"success"`
	Results   []APIResponse `json:"results"`
	Total     int           `json:"total"`
	Succeeded int           `json:"succeeded"`
	Failed    int           `json:"failed"`
	Timestamp time.Time     `json:"timestamp"`
	Duration  string        `json:"duration"`
}

type HealthResponse struct {
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
	Version     string    `json:"version"`
	Environment string    `json:"environment"`
	Uptime      string    `json:"uptime"`
}

type ErrorResponse struct {
	Success   bool      `json:"success"`
	Error     string    `json:"error"`
	Message   string    `json:"message"`
	Code      string    `json:"code,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}
