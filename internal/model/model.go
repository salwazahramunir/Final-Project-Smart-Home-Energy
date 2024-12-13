package model

// TAPAS STRUCT
type Inputs struct {
	Table map[string][]string `json:"table"`
	Query string              `json:"query"`
}

type AIRequest struct {
	Inputs Inputs `json:"inputs"`
}

type TapasResponse struct {
	Answer      string   `json:"answer"`
	Coordinates [][]int  `json:"coordinates"`
	Cells       []string `json:"cells"`
	Aggregator  string   `json:"aggregator"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AI CHAT STRUCT
type ChatRequest struct {
	Model      string    `json:"model"`
	Messages   []Message `json:"messages"`
	Temperatur float32   `json:"temperature"`
	MaxToken   int       `json:"max_tokens"`
	TopP       float32   `json:"top_p"`
	Stream     bool      `json:"stream"`
}

type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}

type ChatResponse struct {
	Choices []Choice `json:"choices"`
}

type TextToSpeechPayload struct {
	Text    string `json:"text"`
	ModelId string `json:"model_id"`
}
