package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/manuel/make-it-rain/db"
	"github.com/manuel/make-it-rain/models"
	"github.com/manuel/make-it-rain/services"
	"github.com/rs/zerolog/log"
)

var auctionService *services.AuctionService

func init() {
	dbService := db.NewDBService()
	auctionService = services.NewAuctionService(dbService)
}

func CreateAuction(c *gin.Context) {
	var req models.CreateAuctionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	auction, err := auctionService.CreateAuction(c.Request.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create auction")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create auction"})
		return
	}

	c.JSON(http.StatusCreated, auction)
}

func CreateAuctionItem(c *gin.Context) {
	var req models.CreateAuctionItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	auctionItem, err := auctionService.CreateAuctionItem(c.Request.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create auction item")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create auction item"})
		return
	}

	c.JSON(http.StatusCreated, auctionItem)
}
func GetAuctions(c *gin.Context) {
	// Params with validation
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	// Validate page and pageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Additional validation for sort parameters (defense in depth)
	// The DB layer also validates these, but we validate here too
	allowedSortColumns := map[string]bool{
		"id":          true,
		"created_at":  true,
		"updated_at":  true,
		"start_date":  true,
		"end_date":    true,
		"start_price": true,
		"title":       true,
		"status":      true,
	}

	if !allowedSortColumns[sortBy] {
		sortBy = "created_at"
	}

	// Normalize sort order to lowercase for consistency
	sortOrder = strings.ToLower(sortOrder)
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	auctions, err := auctionService.GetAuctions(
		c.Request.Context(),
		page,
		pageSize,
		sortBy,
		sortOrder,
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get auctions")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get auctions"})
		return
	}
	c.JSON(http.StatusOK, auctions)
}

// GetAuction by id
func GetAuction(c *gin.Context) {
	auctionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}
	auction, err := auctionService.GetAuction(c.Request.Context(), auctionID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get auction")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get auction"})
		return
	}
	c.JSON(http.StatusOK, auction)
}
