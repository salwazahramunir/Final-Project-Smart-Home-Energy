package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"smart-home-energy/internal/model"
	"smart-home-energy/internal/service"

	"github.com/gin-gonic/gin"
)

// Initialize the services
var fileService = &service.FileService{}
var aiService = &service.AIService{Client: &http.Client{}}

// Retrieve the Hugging Face token from the environment variables
func getToken() (string, error) {
	token := os.Getenv("HUGGINGFACE_TOKEN")

	if token == "" {
		return "", errors.New("HUGGINGFACE_TOKEN is not set in the .env file")
	}

	return token, nil
}

func UploadFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request model.UploadFileRequest

		// Bind the form data
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		file, err := request.File.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open file"})
			return
		}
		defer file.Close()

		// Read file content into a string
		fileContent, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read file content"})
			return
		}

		dataTable, err := fileService.ProcessFile(string(fileContent))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token, err := getToken()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		responseAI, err := aiService.AnalyzeData(dataTable, request.Query, token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		text := fmt.Sprintf("From the provided data, %s: %s", request.Query, responseAI)

		response := model.ResponseSuccess{
			Status: "Success",
			Answer: text,
		}

		c.JSON(http.StatusOK, response)
	}
}

func ChatAI() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request model.ChatAIRequest

		// Bind the form data
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token, err := getToken()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		chatResponse, err := aiService.ChatWithAI(request.Query, token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response := model.ResponseSuccess{
			Status: "Success",
			Answer: chatResponse.Choices[0].Message.Content,
		}

		c.JSON(http.StatusOK, response)
	}
}
