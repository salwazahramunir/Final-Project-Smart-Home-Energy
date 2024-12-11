package service

import (
	"net/http"
	"smart-home-energy/internal/model"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AIService struct {
	Client HTTPClient
}

func (s *AIService) AnalyzeData(table map[string][]string, query, token string) (string, error) {
	return "", nil
}

func (s *AIService) ChatWithAI(context, query, token string) (model.ChatResponse, error) {
	// TODO: answer here
	return model.ChatResponse{}, nil
}
