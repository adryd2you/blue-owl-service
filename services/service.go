package services

import (
	"blue-owl-service/repositories"
	"blue-owl-service/transport"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetServices(c *gin.Context) {
	// Get the projectID from the query parameters
	projectIDStr := c.DefaultQuery("project_id", "") // Default to empty if not provided

	var projectID *int
	if projectIDStr != "" {
		id, err := strconv.Atoi(projectIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project_id"})
			return
		}
		projectID = &id
	}

	// Call GetListServices with the projectID (which can be nil if not provided)
	services, err := repositories.GetListServices(transport.Db, projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a response structure
	response := struct {
		Data []struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Status      string `json:"status"`
		} `json:"data"`
	}{
		Data: make([]struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Status      string `json:"status"`
		}, len(services)),
	}

	// Populate the response structure with the required fields
	for i, service := range services {
		response.Data[i] = struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Status      string `json:"status"`
		}{
			ID:          service.ID,
			Name:        service.Name,
			Description: service.Description,
			Status:      service.Status,
		}
	}

	// Return the modified response in JSON format
	c.JSON(http.StatusOK, response)
}

func GetServiceDetail(c *gin.Context) {
	// Extract the service ID from the path parameters (assuming the URL is like /services/:id)
	serviceIDParam := c.Param("id")

	// Convert the service ID to an integer
	serviceID, err := strconv.Atoi(serviceIDParam)
	if err != nil {
		// Return an error if the service ID is invalid
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	// Call GetServiceByID to retrieve the service details
	service, err := repositories.GetServiceByID(transport.Db, serviceID)
	if err != nil {
		// Return a 404 error if the service is not found, or a 500 error for other errors
		if err.Error() == fmt.Sprintf("service with ID %d not found", serviceID) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve service: %v", err)})
		}
		return
	}

	// Return the service details in the response
	c.JSON(http.StatusOK, service)
}
