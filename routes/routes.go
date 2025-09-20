package routes

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/manuel/make-it-rain/controllers"
	"github.com/manuel/make-it-rain/middleware"
)

func SetupRoutes(r *gin.Engine) {
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	// Serve static files in production
	if _, err := os.Stat("frontend/dist"); err == nil {
		r.Static("/assets", "./frontend/dist/assets")
		r.StaticFile("/", "./frontend/dist/index.html")
	}

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
		// Auction endpoints
		auctions := api.Group("/auctions")
		{
			auctions.POST("", controllers.CreateAuction)
			auctions.POST("/items", controllers.CreateAuctionItem)
			auctions.GET("/:id", controllers.GetAuction)
			auctions.GET("", controllers.GetAuctions)
			// auctions.PUT("/:id", controllers.UpdateAuction)
			// auctions.DELETE("/:id", controllers.DeleteAuction)
		}
	}

	// NoRoute handler - serve SPA for client routes or 404 for API routes
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		// If it's an API route, return 404 JSON
		if len(path) >= 4 && path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"error": "endpoint not found"})
			return
		}
		// For non-API routes, serve the SPA if it exists
		if _, err := os.Stat("frontend/dist/index.html"); err == nil {
			c.File("./frontend/dist/index.html")
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "endpoint not found"})
		}
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
