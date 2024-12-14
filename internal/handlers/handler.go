package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"smart-home-energy/internal/model"
	"smart-home-energy/internal/service"

	"github.com/gin-gonic/gin"
)

// Initialize the services
var fileService = &service.FileService{}
var aiService = &service.AIService{Client: &http.Client{}}

func UploadFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request model.UploadFileRequest

		// bind the form data
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// open the uploaded file
		file, err := request.File.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open file"})
			return
		}
		defer file.Close()

		// read file content into a string
		fileContent, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read file content"})
			return
		}

		// change file content into table data
		dataTable, err := fileService.ProcessFile(string(fileContent))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// analysis of uploaded data
		responseAI, err := aiService.AnalyzeData(dataTable, request.Query)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// format output
		text := fmt.Sprintf("From the provided data, %s: %s", request.Query, responseAI)

		response := model.ResponseSuccess{
			Status: "Success",
			Data:   text,
		}

		c.JSON(http.StatusOK, response)
	}
}

func ChatAI() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request model.ChatAIRequest

		// bind the form data
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// call service to chat with AI
		chatResponse, err := aiService.ChatWithAI(request.Query)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// format output
		response := model.ResponseSuccess{
			Status: "Success",
			Data:   chatResponse.Choices[0].Message.Content,
		}

		c.JSON(http.StatusOK, response)
	}
}

func TextToSpeech() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request model.TextToSpeechRequest

		// bind the form data
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// convert text to audio using elevenlabs
		audioData, err := aiService.GenerateAudioFromElevenLabs(request.Text)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// save temporary audio in local files
		tempFile := "output_audio.mp3"
		err = os.WriteFile(tempFile, audioData, 0644)
		if err != nil {
			fmt.Printf("Error saving audio file: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		defer os.Remove(tempFile) // delete temporary audio in local files

		// upload audio to Cloudinary
		cloudinaryURL, err := aiService.UploadAudioToCloudinary(tempFile)
		if err != nil {
			log.Fatalf("Error uploading to Cloudinary: %v", err)
		}

		// format output
		response := model.ResponseSuccess{
			Status: "Success",
			Data:   cloudinaryURL,
		}

		c.JSON(http.StatusOK, response)
	}
}
