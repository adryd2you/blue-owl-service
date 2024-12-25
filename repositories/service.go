package repositories

import (
	"blue-owl-service/models"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func GetListServices(db *sqlx.DB, projectID *int) ([]models.APIService, error) {
	// Initialize the base query
	query := `SELECT id, api_project_id, name, endpoint, description, status 
	          FROM api_services`

	// Prepare for filtering by projectID (if provided)
	var args []interface{}
	if projectID != nil {
		query += " WHERE api_project_id = $1"
		args = append(args, *projectID)
	}

	// Execute the query
	var services []models.APIService
	err := db.Select(&services, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve API services: %w", err)
	}

	return services, nil
}

// GetProjectByID fetches a single project by ID from the database
func GetServiceByID(db *sqlx.DB, serviceID int) (*models.APIService, error) {
	// Define the query to retrieve the service by its ID
	query := `SELECT id, api_project_id, name, endpoint, description, status 
	          FROM api_services WHERE id = $1`

	// Define a variable to hold the result
	var service models.APIService

	// Execute the query
	err := db.Get(&service, query, serviceID)
	if err != nil {
		// If no rows are found or another error occurs, return it
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("service with ID %d not found", serviceID)
		}
		return nil, fmt.Errorf("failed to retrieve service by ID: %w", err)
	}

	return &service, nil
}
