package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"response-std/config"
	"response-std/libs/external/requests"
	"response-std/libs/responses"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

type APIClient struct {
	client *resty.Client
	logger *Logger
	config *config.Config
}

type RequestBuilder struct {
	client      *APIClient
	url         string
	method      string
	headers     map[string]string
	queryParams map[string]string
	body        interface{}
	timeout     *time.Duration
	requestID   string
}

func NewAPIClient(cfg *config.Config, logger *Logger) *APIClient {
	client := resty.New()
	client.SetTimeout(cfg.RequestTimeout)
	client.SetRetryCount(cfg.MaxRetries)
	client.SetRetryWaitTime(cfg.RetryDelay)

	// Set default headers
	client.SetHeaders(map[string]string{
		"Accept": "application/json",
	})

	return &APIClient{
		client: client,
		logger: logger,
		config: cfg,
	}
}

// HTTP Method builders
func (ac *APIClient) Get(url string) *RequestBuilder {
	return &RequestBuilder{
		client:      ac,
		url:         url,
		method:      "GET",
		headers:     make(map[string]string),
		queryParams: make(map[string]string),
		requestID:   uuid.New().String(),
	}
}

func (ac *APIClient) Post(url string) *RequestBuilder {
	return &RequestBuilder{
		client:      ac,
		url:         url,
		method:      "POST",
		headers:     make(map[string]string),
		queryParams: make(map[string]string),
		requestID:   uuid.New().String(),
	}
}

func (ac *APIClient) Put(url string) *RequestBuilder {
	return &RequestBuilder{
		client:      ac,
		url:         url,
		method:      "PUT",
		headers:     make(map[string]string),
		queryParams: make(map[string]string),
		requestID:   uuid.New().String(),
	}
}

func (ac *APIClient) Delete(url string) *RequestBuilder {
	return &RequestBuilder{
		client:      ac,
		url:         url,
		method:      "DELETE",
		headers:     make(map[string]string),
		queryParams: make(map[string]string),
		requestID:   uuid.New().String(),
	}
}

func (ac *APIClient) Patch(url string) *RequestBuilder {
	return &RequestBuilder{
		client:      ac,
		url:         url,
		method:      "PATCH",
		headers:     make(map[string]string),
		queryParams: make(map[string]string),
		requestID:   uuid.New().String(),
	}
}

func (ac *APIClient) Head(url string) *RequestBuilder {
	return &RequestBuilder{
		client:      ac,
		url:         url,
		method:      "HEAD",
		headers:     make(map[string]string),
		queryParams: make(map[string]string),
		requestID:   uuid.New().String(),
	}
}

func (ac *APIClient) Options(url string) *RequestBuilder {
	return &RequestBuilder{
		client:      ac,
		url:         url,
		method:      "OPTIONS",
		headers:     make(map[string]string),
		queryParams: make(map[string]string),
		requestID:   uuid.New().String(),
	}
}

// RequestBuilder methods for chaining
func (rb *RequestBuilder) WithHeader(key, value string) *RequestBuilder {
	rb.headers[key] = value
	return rb
}

func (rb *RequestBuilder) WithHeaders(headers map[string]string) *RequestBuilder {
	for k, v := range headers {
		rb.headers[k] = v
	}
	return rb
}

func (rb *RequestBuilder) WithBearerToken(token string) *RequestBuilder {
	rb.headers["Authorization"] = fmt.Sprintf("Bearer %s", token)
	return rb
}

func (rb *RequestBuilder) WithAPIKey(key string) *RequestBuilder {
	rb.headers["X-API-Key"] = key
	return rb
}

func (rb *RequestBuilder) WithBasicAuth(username, password string) *RequestBuilder {
	rb.headers["Authorization"] = fmt.Sprintf("Basic %s", basicAuth(username, password))
	return rb
}

func (rb *RequestBuilder) WithContentType(contentType string) *RequestBuilder {
	rb.headers["Content-Type"] = contentType
	return rb
}

func (rb *RequestBuilder) WithJSONContentType() *RequestBuilder {
	rb.headers["Content-Type"] = "application/json"
	return rb
}

func (rb *RequestBuilder) WithXMLContentType() *RequestBuilder {
	rb.headers["Content-Type"] = "application/xml"
	return rb
}

func (rb *RequestBuilder) WithFormContentType() *RequestBuilder {
	rb.headers["Content-Type"] = "application/x-www-form-urlencoded"
	return rb
}

func (rb *RequestBuilder) WithBody(body interface{}) *RequestBuilder {
	rb.body = body
	return rb
}

func (rb *RequestBuilder) WithJSONBody(data interface{}) *RequestBuilder {
	rb.body = data
	rb.headers["Content-Type"] = "application/json"
	return rb
}

func (rb *RequestBuilder) WithQuery(key, value string) *RequestBuilder {
	rb.queryParams[key] = value
	return rb
}

func (rb *RequestBuilder) WithQueryParams(params map[string]string) *RequestBuilder {
	for k, v := range params {
		rb.queryParams[k] = v
	}
	return rb
}

func (rb *RequestBuilder) WithTimeout(timeout time.Duration) *RequestBuilder {
	rb.timeout = &timeout
	return rb
}

func (rb *RequestBuilder) WithTimeoutSeconds(seconds int) *RequestBuilder {
	timeout := time.Duration(seconds) * time.Second
	rb.timeout = &timeout
	return rb
}

func (rb *RequestBuilder) WithUserAgent(userAgent string) *RequestBuilder {
	rb.headers["User-Agent"] = userAgent
	return rb
}

func (rb *RequestBuilder) WithReferer(referer string) *RequestBuilder {
	rb.headers["Referer"] = referer
	return rb
}

func (rb *RequestBuilder) WithCustomHeader(key, value string) *RequestBuilder {
	rb.headers[key] = value
	return rb
}

// Execute the request
func (rb *RequestBuilder) Execute() *responses.APIResponse {
	startTime := time.Now()

	response := &responses.APIResponse{
		RequestID: rb.requestID,
		Timestamp: startTime,
	}

	// Validate URL
	if rb.url == "" {
		return &responses.APIResponse{
			Success:   false,
			Error:     "URL is required",
			Message:   "Invalid request",
			Timestamp: time.Now(),
		}
	}

	// Clone client for this request to avoid race conditions
	reqClient := rb.client.client.R()

	// Set custom timeout if provided
	if rb.timeout != nil {
		timeoutCtx, cancel := context.WithTimeout(reqClient.Context(), *rb.timeout)
		defer cancel()
		reqClient.SetContext(timeoutCtx)
	}

	// Set headers
	if len(rb.headers) > 0 {
		reqClient.SetHeaders(rb.headers)
	}

	// Set query parameters
	if len(rb.queryParams) > 0 {
		reqClient.SetQueryParams(rb.queryParams)
	}

	// Set body
	if rb.body != nil {
		reqClient.SetBody(rb.body)
	}

	// Execute request
	var resp *resty.Response
	var err error

	switch rb.method {
	case "GET":
		resp, err = reqClient.Get(rb.url)
	case "POST":
		resp, err = reqClient.Post(rb.url)
	case "PUT":
		resp, err = reqClient.Put(rb.url)
	case "DELETE":
		resp, err = reqClient.Delete(rb.url)
	case "PATCH":
		resp, err = reqClient.Patch(rb.url)
	case "HEAD":
		resp, err = reqClient.Head(rb.url)
	case "OPTIONS":
		resp, err = reqClient.Options(rb.url)
	default:
		err = fmt.Errorf("unsupported HTTP method: %s", rb.method)
	}

	duration := time.Since(startTime)
	response.Duration = duration.String()

	// Handle response
	if err != nil {
		response.Success = false
		response.Error = err.Error()
		response.Message = "Request failed"
	} else {
		response.StatusCode = resp.StatusCode()

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
		response.Headers = make(map[string]string)
		for k, v := range resp.Header() {
			if len(v) > 0 {
				response.Headers[k] = v[0]
			}
		}

		if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
			response.Success = true
			response.Message = "Request successful"
		} else {
			response.Success = false
			response.Message = fmt.Sprintf("Request failed with status code: %d", resp.StatusCode())
		}
	}

	// Log the request if logging is enabled
	if rb.client.config.EnableLogging {
		rb.logRequest(startTime, duration, response)
	}

	return response
}

// Execute async
func (rb *RequestBuilder) ExecuteAsync() <-chan *responses.APIResponse {
	resultChan := make(chan *responses.APIResponse, 1)

	go func() {
		defer close(resultChan)
		resultChan <- rb.Execute()
	}()

	return resultChan
}

// Batch execution
type BatchRequest struct {
	requests []*RequestBuilder
	parallel bool
}

func (ac *APIClient) Batch() *BatchRequest {
	return &BatchRequest{
		requests: make([]*RequestBuilder, 0),
		parallel: false,
	}
}

func (br *BatchRequest) Add(rb *RequestBuilder) *BatchRequest {
	br.requests = append(br.requests, rb)
	return br
}

func (br *BatchRequest) Parallel() *BatchRequest {
	br.parallel = true
	return br
}

func (br *BatchRequest) Execute() *responses.BatchAPIResponse {
	startTime := time.Now()

	batchResponse := &responses.BatchAPIResponse{
		Total:     len(br.requests),
		Timestamp: startTime,
		Results:   make([]responses.APIResponse, len(br.requests)),
	}

	if br.parallel {
		// Execute requests in parallel
		var wg sync.WaitGroup
		var mu sync.Mutex

		for i, req := range br.requests {
			wg.Add(1)
			go func(index int, rb *RequestBuilder) {
				defer wg.Done()

				result := rb.Execute()

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
		for i, req := range br.requests {
			result := req.Execute()
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

// Helper functions
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64Encode([]byte(auth))
}

func base64Encode(data []byte) string {
	const base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

	buf := make([]byte, (len(data)+2)/3*4)

	for i := 0; i < len(data); i += 3 {
		n := uint32(data[i]) << 16
		if i+1 < len(data) {
			n |= uint32(data[i+1]) << 8
		}
		if i+2 < len(data) {
			n |= uint32(data[i+2])
		}

		j := i / 3 * 4
		buf[j] = base64Table[n>>18&63]
		buf[j+1] = base64Table[n>>12&63]
		buf[j+2] = base64Table[n>>6&63]
		buf[j+3] = base64Table[n&63]
	}

	// Add padding
	switch len(data) % 3 {
	case 1:
		buf[len(buf)-2] = '='
		buf[len(buf)-1] = '='
	case 2:
		buf[len(buf)-1] = '='
	}

	return string(buf)
}

func (rb *RequestBuilder) logRequest(startTime time.Time, duration time.Duration, response *responses.APIResponse) {
	requestLog := &requests.RequestLog{
		ID:          rb.requestID,
		Method:      rb.method,
		URL:         rb.url,
		Headers:     rb.headers,
		Body:        rb.body,
		QueryParams: rb.queryParams,
		Timestamp:   startTime,
		Duration:    duration,
		StatusCode:  response.StatusCode,
		Response:    response.Data,
		Error:       response.Error,
	}

	rb.client.logger.LogRequest(requestLog)
}
