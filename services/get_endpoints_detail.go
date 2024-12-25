package services

import (
	"blue-owl-service/models"
	"blue-owl-service/repositories"
	"blue-owl-service/transport"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAPIEndpointDetail retrieves a single API endpoint by ID
func GetAPIEndpointDetail(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing ID parameter"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	apiEndpoint, err := repositories.GetAPIEndpointByID(id, transport.Db)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("API endpoint with ID %d not found", id)})
		return
	}

	apiRequestsContract, err := repositories.GetAPIRequests(transport.Db, apiEndpoint.ID)
	if err != nil {
		apiRequestsContract = []models.RequestContract{}
	}

	apiResponseContract, err := repositories.GetAPIResponses(transport.Db, apiEndpoint.ID)
	if err != nil {
		apiResponseContract = []models.ResponseContract{}
	}

	apiExamples, err := repositories.GetAPIExampleByEndpointID(transport.Db, apiEndpoint.ID)
	if err != nil {
		apiExamples = []models.APIExampleUnParsed{}
	}

	response := gin.H{
		"api_endpoint":       apiEndpoint,
		"request_contracts":  apiRequestsContract,
		"response_contracts": apiResponseContract,
		"examples":           apiExamples,
	}

	c.JSON(http.StatusOK, response)
}
