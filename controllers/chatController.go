package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/manuel/make-it-rain/db"
	"github.com/manuel/make-it-rain/models"
	"github.com/manuel/make-it-rain/services"
	"github.com/rs/zerolog/log"
)

var chatService *services.ChatService

func init() {
	dbService := db.NewDBService()
	chatService = services.NewChatService(dbService)
}

func CreateChat(c *gin.Context) {
	var req models.CreateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chat, err := chatService.CreateChat(c.Request.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create chat")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
		return
	}

	c.JSON(http.StatusCreated, chat)
}

func SendMessage(c *gin.Context) {
	chatID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	var req models.CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ChatID = chatID
	message, err := chatService.SendMessage(c.Request.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send message")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}
	c.JSON(http.StatusCreated, message)

}

func GetChat(c *gin.Context) {
	chatID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	chat, err := chatService.GetChat(c.Request.Context(), chatID)
	if err != nil {
		if errors.Is(err, db.ErrChatNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
			return
		}
		log.Error().Err(err).Int64("chat_id", chatID).Msg("Failed to get chat")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat"})
		return
	}

	c.JSON(http.StatusOK, chat)
}
