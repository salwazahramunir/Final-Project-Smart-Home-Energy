package handler

import (
	"net/http"
	"smart-home-energy/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// Initialize the services
var fileService = &service.FileService{}
var aiService = &service.AIService{Client: &http.Client{}}
var store = sessions.NewCookieStore([]byte("my-key"))

func getSession(r *http.Request) *sessions.Session {
	session, _ := store.Get(r, "chat-session")
	return session
}

func UploadFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"result": "haii"})
	}
}

func ChatAI() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"result": "haii"})
	}
}
