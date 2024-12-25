package services

import (
	"blue-owl-service/repositories"
	"blue-owl-service/transport"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAPIEndpoints retrieves all API endpoints
func GetAPIEndpoints(c *gin.Context) {
	// Get the serviceID from query parameters (if available)
	serviceIDStr := c.DefaultQuery("service_id", "") // Default to empty string if not provided

	var serviceID *int
	if serviceIDStr != "" {
		// Convert the service_id query parameter to an integer
		id, err := strconv.Atoi(serviceIDStr)
		if err != nil {
			// If the service_id is not a valid integer, return a 400 Bad Request
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service_id"})
			return
		}
		serviceID = &id
	}

	// Call GetAPIEndpoints with the serviceID (which can be nil if not provided)
	endpoints, err := repositories.GetAPIEndpoints(transport.Db, serviceID) // db is your *sqlx.DB instance
	if err != nil {
		// Handle errors during database retrieval
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the list of API endpoints in JSON format
	c.JSON(http.StatusOK, endpoints)
}
