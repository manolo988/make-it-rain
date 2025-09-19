package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/manuel/make-it-rain/db"
	"github.com/manuel/make-it-rain/services"
	"github.com/rs/zerolog/log"
)

var userMessageService *services.UserMessageService

func init() {
	dbService := db.NewDBService()
	userMessageService = services.NewUserMessageService(dbService)
}

func GetUserMessages(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	messages, err := userMessageService.GetUserMessages(c.Request.Context(), userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user messages")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user messages"})
		return
	}
	c.JSON(http.StatusOK, messages)
}
