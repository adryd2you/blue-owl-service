package services

import (
	"blue-owl-service/models"
	"blue-owl-service/repositories"
	"blue-owl-service/transport"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetProjects retrieves the list of projects
func GetProjects(c *gin.Context) {
	projects, err := repositories.GetListProjects(transport.Db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
		return
	}

	// Respond with the list of projects
	c.JSON(http.StatusOK, gin.H{"data": projects})
}

// GetProjectDetail retrieves the details of a specific project by ID
func GetProjectDetail(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Call the database function to fetch the project
	project, err := repositories.GetProjectByID(transport.Db, projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Respond with the project details
	c.JSON(http.StatusOK, gin.H{"data": project})
}

// CreateProject handles creating a new project
func CreateProject(c *gin.Context) {
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Check if a project with the same name already exists
	exists, err := repositories.CheckProjectExistsByName(transport.Db, project.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate project name"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Project with the same name already exists"})
		return
	}

	// Call the database function to create the project
	projectID, err := repositories.CreateProject(transport.Db, &project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	// Respond with the created project ID
	c.JSON(http.StatusCreated, gin.H{"data": gin.H{"id": projectID}})
}

// UpdateProject handles updating an existing project
func UpdateProject(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	project.ID = projectID // Ensure the project ID is set

	// Call the database function to update the project
	err = repositories.UpdateProject(transport.Db, &project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "Project updated successfully"})
}

// DeleteProject handles deleting a project by ID
func DeleteProject(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Call the database function to delete the project
	err = repositories.DeleteProject(transport.Db, projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}
