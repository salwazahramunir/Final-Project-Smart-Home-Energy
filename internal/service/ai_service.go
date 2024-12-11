package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"smart-home-energy/internal/model"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AIService struct {
	Client HTTPClient
}

const HUGGING_URL = "https://api-inference.huggingface.co/models/"

func (s *AIService) AnalyzeData(table map[string][]string, query, token string) (string, error) {
	if len(table) == 0 {
		return "", errors.New("table data is empty")
	}

	inputs := model.AIRequest{
		Inputs: model.Inputs{
			Table: table,
			Query: query,
		},
	}

	jsonData, err := json.Marshal(inputs)

	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return "", err
	}

	request, err := http.NewRequest("POST", HUGGING_URL+"google/tapas-base-finetuned-wtq", bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", "application/json")

	response, err := s.Client.Do(request)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error: %s", http.StatusText(response.StatusCode))
	}

	var tapasResponse model.TapasResponse

	err = json.NewDecoder(response.Body).Decode(&tapasResponse)
	if err != nil {
		fmt.Println("Error decode JSON:", err)
		return "", err
	}

	return tapasResponse.Cells[0], nil
}

func (s *AIService) ChatWithAI(context, query, token string) (model.ChatResponse, error) {
	// TODO: answer here
	return model.ChatResponse{}, nil
}
