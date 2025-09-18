package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/manuel/make-it-rain/controllers"
	"github.com/manuel/make-it-rain/middleware"
)

func SetupRoutes(r *gin.Engine) {
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	r.GET("/health", HealthCheck)
	r.GET("/ready", ReadinessCheck)

	api := r.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.POST("", controllers.CreateUser)
			users.GET("/:id", controllers.GetUser)
			users.GET("", controllers.GetUsers)
			users.PUT("/:id", controllers.UpdateUser)
			users.DELETE("/:id", controllers.DeleteUser)
		}
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "endpoint not found"})
	})
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}

func ReadinessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}