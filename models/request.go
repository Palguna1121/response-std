package models

import (
	"errors"
	"strings"
	"time"
)

type APIRequest struct {
	Method      string            `json:"method" binding:"required"`
	URL         string            `json:"url" binding:"required,url"`
	Headers     map[string]string `json:"headers,omitempty"`
	Body        interface{}       `json:"body,omitempty"`
	QueryParams map[string]string `json:"query_params,omitempty"`
	Timeout     *int              `json:"timeout,omitempty"` // in seconds
	Retries     *int              `json:"retries,omitempty"`
	RetryDelay  *int              `json:"retry_delay,omitempty"` // in seconds
}

type BatchAPIRequest struct {
	Requests []APIRequest `json:"requests" binding:"required,min=1,max=10"`
	Parallel bool         `json:"parallel,omitempty"`
}

type RequestLog struct {
	ID          string            `json:"id"`
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers"`
	Body        interface{}       `json:"body"`
	QueryParams map[string]string `json:"query_params"`
	Timestamp   time.Time         `json:"timestamp"`
	Duration    time.Duration     `json:"duration"`
	StatusCode  int               `json:"status_code"`
	Response    interface{}       `json:"response"`
	Error       string            `json:"error,omitempty"`
	Retries     int               `json:"retries"`
}

// Validation
func (r *APIRequest) Validate() error {
	r.Method = strings.ToUpper(strings.TrimSpace(r.Method))

	validMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	isValid := false
	for _, method := range validMethods {
		if r.Method == method {
			isValid = true
			break
		}
	}

	if !isValid {
		return errors.New("invalid HTTP method")
	}

	if r.Timeout != nil && *r.Timeout > 300 { // max 5 minutes
		return errors.New("timeout cannot exceed 300 seconds")
	}

	if r.Retries != nil && *r.Retries > 5 { // max 5 retries
		return errors.New("retries cannot exceed 5")
	}

	return nil
}
