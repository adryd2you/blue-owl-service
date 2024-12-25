package route

import (
	"blue-owl-service/services"

	"github.com/gin-gonic/gin"
)

func InitRoute() {
	router := gin.Default()
	router.GET("/projects", services.GetProjects)
	router.GET("/projects/:id", services.GetProjectDetail)

	router.GET("/services", services.GetServices)
	router.GET("/services/:id", services.GetServiceDetail)

	router.GET("/api-endpoints", services.GetAPIEndpoints)
	router.POST("/api-endpoints", services.CreateAPIEndpoint)
	router.DELETE("/api-endpoints/:id", services.DeleteAPIEndpoint)
	router.GET("/api-endpoints/:id", services.GetAPIEndpointDetail)
	router.PUT("/api-endpoints/:id", services.UpdateAPIEndpoint)

	router.POST("/test", services.RunTests)
	router.POST("/project-tests", services.RunProjectTests)
	router.POST("/service-tests", services.RunServiceTests)
	router.POST("/endpoint-tests", services.RunEndpointTests)
	router.POST("/specific-test", services.RunSpecificTest)
	router.POST("/ai", services.HitAI)
	router.Run("localhost:8080")
}
