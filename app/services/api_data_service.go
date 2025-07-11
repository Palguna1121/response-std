package services

import (
	"response-std/config"
	"response-std/libs/external/services"
	"sync"
)

type APIDataService struct {
	client *services.APIClient
}

var (
	log  *services.Logger
	once sync.Once
)

func InitLogger() {
	once.Do(func() {
		log = services.NewLogger(config.ENV.LogLevel, config.ENV.Environment)
	})
}

func NewAPIDataService(cfg *config.Config) *APIDataService {
	// Pastikan logger sudah diinisialisasi
	InitLogger()

	return &APIDataService{
		client: services.NewAPIClient(cfg, log),
	}
}

func (s *APIDataService) GetUserData() interface{} {
	res := s.client.Get("https://dummyjson.com/users").Execute()
	return res.Data
}
