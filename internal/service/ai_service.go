package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"smart-home-energy/internal/helper"
	"smart-home-energy/internal/model"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
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

func (s *AIService) ChatWithAI(query, token string) (model.ChatResponse, error) {
	// TODO: answer here
	inputs := model.ChatRequest{
		Model: "mistralai/Mistral-Nemo-Instruct-2407",
		Messages: []model.Message{
			{
				Role:    "user",
				Content: query,
			},
		},
		Temperatur: 0.5,
		MaxToken:   2048,
		TopP:       0.7,
		Stream:     false,
	}

	jsonData, err := json.Marshal(inputs)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return model.ChatResponse{}, err
	}

	request, err := http.NewRequest("POST", HUGGING_URL+"mistralai/Mistral-Nemo-Instruct-2407/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return model.ChatResponse{}, err
	}

	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", "application/json")

	response, err := s.Client.Do(request)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return model.ChatResponse{}, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return model.ChatResponse{}, fmt.Errorf("errors: %s", http.StatusText(response.StatusCode))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error read body:", err)
		return model.ChatResponse{}, err
	}

	var chatResponse model.ChatResponse

	err = json.Unmarshal(body, &chatResponse)
	if err != nil {
		fmt.Println("Error unmarshalling:", err)
		return model.ChatResponse{}, err
	}

	return chatResponse, nil
}

func (s *AIService) GenerateAudioFromElevenLabs(text string) ([]byte, error) {
	ELEVENLAB_TOKEN, err := helper.GetToken("ELEVENLAB_TOKEN")
	if err != nil {
		fmt.Printf("error : %s", err.Error())
		return nil, err
	}

	inputs := model.TextToSpeechPayload{
		Text:    text,
		ModelId: "eleven_multilingual_v2",
	}

	jsonData, err := json.Marshal(inputs)
	if err != nil {
		fmt.Printf("error encoding JSON: %s", err.Error())
		return nil, err
	}

	request, err := http.NewRequest("POST", "https://api.elevenlabs.io/v1/text-to-speech/21m00Tcm4TlvDq8ikWAM", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("error creating request: %s", err.Error())
		return nil, err
	}

	request.Header.Set("xi-api-key", ELEVENLAB_TOKEN)
	request.Header.Set("Content-Type", "application/json")

	response, err := s.Client.Do(request)
	if err != nil {
		fmt.Printf("error creating request: %s", err.Error())
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("errors: %s", http.StatusText(response.StatusCode))
	}

	audioData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return audioData, nil
}

func (s *AIService) UploadAudioToCloudinary(filePath string) (string, error) {
	cld, err := helper.InitCloudinary()
	if err != nil {
		fmt.Printf("Error call InitCloudinary function: %v", err)
	}

	var ctx = context.Background()

	// Unggah file ke Cloudinary
	uploadResult, err := cld.Upload.Upload(ctx, filePath, uploader.UploadParams{
		ResourceType: "video",
		Folder:       "audio_files",
		PublicID:     "my_audio",
	})

	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}
