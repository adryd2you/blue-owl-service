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

// DeleteAPIEndpoint deletes an API endpoint by ID
func DeleteAPIEndpoint(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	deleted, err := repositories.DeleteAPIEndpoint(id, transport.Db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error deleting API endpoint: %v", err)})
		return
	}

	if !deleted {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("API endpoint with ID %d not found", id)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("API endpoint with ID %d deleted successfully", id)})
}

func UpdateAPIEndpoint(c *gin.Context) {
	id := c.Param("id")

	// Parse ID to int
	intID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var body struct {
		Name         *string `json:"name"`
		Path         *string `json:"path"`
		Method       *string `json:"method"`
		Description  *string `json:"description"`
		APIType      *string `json:"api_type"`
		APIServiceID *int    `json:"api_service_id"`
	}

	// Bind JSON input to the struct
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid input: %v", err)})
		return
	}

	// Validate fields if provided
	if body.Method != nil && !isValidHTTPMethod(*body.Method) {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid HTTP method: %s", *body.Method)})
		return
	}
	if body.APIType != nil && !isValidAPICategory(*body.APIType) {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid API type: %s", *body.APIType)})
		return
	}

	// Map body to models.APIEndpoint with dereferencing
	updatedData := models.APIEndpoint{
		Name:        dereferenceString(body.Name),
		Path:        dereferenceString(body.Path),
		Method:      dereferenceString(body.Method),
		Description: dereferenceString(body.Description),
		APIType:     dereferenceString(body.APIType),
		ServiceID:   dereferenceInt(body.APIServiceID),
	}

	// Update the API endpoint in the database
	updated, err := repositories.UpdateAPIEndpoint(intID, updatedData, transport.Db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error updating API endpoint: %v", err)})
		return
	}

	if !updated {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("API endpoint with ID %s not found", id)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("API endpoint with ID %s updated successfully", id)})
}

// Helper functions to handle nil pointers
func dereferenceString(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
}

func dereferenceInt(ptr *int) int {
	if ptr != nil {
		return *ptr
	}
	return 0
}

func isValidAPICategory(apiType string) bool {
	validCategories := map[string]bool{
		"Internal": true, "External": true, "ThirdParty": true,
	}
	return validCategories[apiType]
}

// Helper functions to validate enums
func isValidHTTPMethod(method string) bool {
	validMethods := map[string]bool{
		"GET": true, "POST": true, "PUT": true, "DELETE": true, "PATCH": true,
	}
	return validMethods[method]
}
