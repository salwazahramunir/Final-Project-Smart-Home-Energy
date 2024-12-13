package model

import "mime/multipart"

// Request From Body

// TAPAS
type UploadFileRequest struct {
	File  *multipart.FileHeader `form:"file" binding:"required"`
	Query string                `form:"question" binding:"required"`
}

// CHAT AI
type ChatAIRequest struct {
	Query string `json:"query"`
}

// TEXT TO SPEECH
type TextToSpeechRequest struct {
	Text string `json:"text"`
}
