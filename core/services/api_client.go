package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"response-std/config"
	"response-std/core/models"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

type APIClient struct {
	client *resty.Client
	logger *Logger
	config *config.Config
}

func NewAPIClient(cfg *config.Config, logger *Logger) *APIClient {
	client := resty.New()
	client.SetTimeout(cfg.RequestTimeout)
	client.SetRetryCount(cfg.MaxRetries)
	client.SetRetryWaitTime(cfg.RetryDelay)

	// Set default headers
	client.SetHeaders(map[string]string{
		"Accept":     "application/json",
	})

	return &APIClient{
		client: client,
		logger: logger,
		config: cfg,
	}
}

func (ac *APIClient) ExecuteRequest(apiReq *models.APIRequest) *models.APIResponse {
	// Validate request
	if err := apiReq.Validate(); err != nil {
		return &models.APIResponse{
			Success:   false,
			Error:     err.Error(),
			Message:   "Invalid request",
			Timestamp: time.Now(),
		}
	}

	requestID := uuid.New().String()
	startTime := time.Now()

	response := &models.APIResponse{
		RequestID: requestID,
		Timestamp: startTime,
	}

	// Create request log
	requestLog := &models.RequestLog{
		ID:          requestID,
		Method:      apiReq.Method,
		URL:         apiReq.URL,
		Headers:     apiReq.Headers,
		Body:        apiReq.Body,
		QueryParams: apiReq.QueryParams,
		Timestamp:   startTime,
	}

	// Clone client for this request to avoid race conditions
	reqClient := ac.client.R()

	// Set custom timeout if provided
	if apiReq.Timeout != nil {
		timeoutCtx, cancel := context.WithTimeout(reqClient.Context(), time.Duration(*apiReq.Timeout)*time.Second)
		defer cancel()
		reqClient.SetContext(timeoutCtx)
	}

	// Set headers
	if apiReq.Headers != nil {
		reqClient.SetHeaders(apiReq.Headers)
	}

	// Set query parameters
	if apiReq.QueryParams != nil {
		reqClient.SetQueryParams(apiReq.QueryParams)
	}

	// Set body
	if apiReq.Body != nil {
		reqClient.SetBody(apiReq.Body)
	}

	// Execute request
	var resp *resty.Response
	var err error

	switch apiReq.Method {
	case "GET":
		resp, err = reqClient.Get(apiReq.URL)
	case "POST":
		resp, err = reqClient.Post(apiReq.URL)
	case "PUT":
		resp, err = reqClient.Put(apiReq.URL)
	case "DELETE":
		resp, err = reqClient.Delete(apiReq.URL)
	case "PATCH":
		resp, err = reqClient.Patch(apiReq.URL)
	case "HEAD":
		resp, err = reqClient.Head(apiReq.URL)
	case "OPTIONS":
		resp, err = reqClient.Options(apiReq.URL)
	default:
		err = fmt.Errorf("unsupported HTTP method: %s", apiReq.Method)
	}

	duration := time.Since(startTime)
	response.Duration = duration.String()

	// Handle response
	if err != nil {
		response.Success = false
		response.Error = err.Error()
		response.Message = "Request failed"
		requestLog.Error = err.Error()
	} else {
		response.StatusCode = resp.StatusCode()
		requestLog.StatusCode = resp.StatusCode()

		// Parse response
		var responseData interface{}
		if len(resp.Body()) > 0 {
			contentType := resp.Header().Get("Content-Type")
			if strings.Contains(contentType, "application/json") {
				if json.Unmarshal(resp.Body(), &responseData) != nil {
					responseData = string(resp.Body())
				}
			} else {
				responseData = string(resp.Body())
			}
		}

		response.Data = responseData
		requestLog.Response = responseData

		if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
			response.Success = true
			response.Message = "Request successful"
		} else {
			response.Success = false
			response.Message = fmt.Sprintf("Request failed with status code: %d", resp.StatusCode())
		}
	}

	// Update request log
	requestLog.Duration = duration

	// Log the request if logging is enabled
	if ac.config.EnableLogging {
		ac.logger.LogRequest(requestLog)
	}

	return response
}

func (ac *APIClient) ExecuteBatchRequest(batchReq *models.BatchAPIRequest) *models.BatchAPIResponse {
	startTime := time.Now()

	batchResponse := &models.BatchAPIResponse{
		Total:     len(batchReq.Requests),
		Timestamp: startTime,
		Results:   make([]models.APIResponse, len(batchReq.Requests)),
	}

	if batchReq.Parallel {
		// Execute requests in parallel
		var wg sync.WaitGroup
		var mu sync.Mutex

		for i, req := range batchReq.Requests {
			wg.Add(1)
			go func(index int, apiReq models.APIRequest) {
				defer wg.Done()

				result := ac.ExecuteRequest(&apiReq)

				mu.Lock()
				batchResponse.Results[index] = *result
				if result.Success {
					batchResponse.Succeeded++
				} else {
					batchResponse.Failed++
				}
				mu.Unlock()
			}(i, req)
		}

		wg.Wait()
	} else {
		// Execute requests sequentially
		for i, req := range batchReq.Requests {
			result := ac.ExecuteRequest(&req)
			batchResponse.Results[i] = *result

			if result.Success {
				batchResponse.Succeeded++
			} else {
				batchResponse.Failed++
			}
		}
	}

	batchResponse.Duration = time.Since(startTime).String()
	batchResponse.Success = batchResponse.Failed == 0

	return batchResponse
}
