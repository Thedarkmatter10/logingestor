package routes

import (
	"logingestor/controllers"

	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine {
	route := gin.Default()
	route.POST("/ingest", controllers.LogIngestHandler)
	route.GET("/search", controllers.SearchLogsHandler)
	return route
}
