package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/manuel/make-it-rain/db"
	"github.com/manuel/make-it-rain/services"
	"github.com/rs/zerolog/log"
)

var jobScheduler *services.JobScheduler

func init() {
	dbService := db.NewDBService()
	jobScheduler = services.NewJobScheduler(dbService)
}

func SubmitJob(c *gin.Context) {
	var req db.SubmitJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job, err := jobScheduler.SubmitJob(c.Request.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to submit job")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit job"})
		return
	}
	c.JSON(http.StatusCreated, job)
}

func GetJobs(c *gin.Context) {
	// Query params
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")
	status := c.DefaultQuery("status", "submitted")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	jobs, err := jobScheduler.GetJobs(
		c.Request.Context(),
		page,
		pageSize,
		sortBy,
		sortOrder,
		status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Jobs not found"})
			return
		}
		log.Error().Err(err).Msg("Failed to get jobs")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get jobs"})
		return
	}

	c.JSON(http.StatusOK, jobs)
}

func GetJobByID(c *gin.Context) {
	jobID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	result, err := jobScheduler.GetJob(c.Request.Context(), jobID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}
		log.Error().Err(err).Int64("job_id", jobID).Msg("Failed to get job")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get job"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func SearchJob(c *gin.Context) {
	// Query params, searched by metadata JSONB
	var metadata map[string]interface{}
	if err := c.ShouldBindJSON(&metadata); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := jobScheduler.SearchJob(c.Request.Context(), metadata)
	if err != nil {
		log.Error().Err(err).Msg("Failed to search job")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search job"})
		return
	}
	c.JSON(http.StatusOK, result)
}

// // TODO
// func DeleteJob(c *gin.Context) {

// 	// c.JSON(http.StatusOK, gin.H{"message": "Job deleted successfully"})
// }
