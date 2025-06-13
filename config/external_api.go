package config

import (
	"os"
)

var ExternalAPIBaseURL string

var ExternalAPIEndpoint string

func loadExternalAPIBaseURL() {
	ExternalAPIBaseURL = os.Getenv("EXTERNAL_API_BASE_URL")
	if ExternalAPIBaseURL == "" {
		ExternalAPIBaseURL = "http://localhost:8080"
	}
}

func loadExternalAPIEndpoint() {
	ExternalAPIEndpoint = os.Getenv("EXTERNAL_API_ENDPOINT")
	if ExternalAPIEndpoint == "" {
		ExternalAPIEndpoint = "/api/v1/"
	}
}
