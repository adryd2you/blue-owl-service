package repositories

import (
	"blue-owl-service/models"
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// GetAPIEndpoints retrieves all API endpoints
func GetAPIEndpoints(db *sqlx.DB, serviceID *int) ([]models.APIEndpoint, error) {
	// Initialize the base query with concatenation of s.endpoint and e.path
	query := `
        SELECT 
            e.id, 
            e.api_service_id, 
            e.method, 
            e.name, 
            e.path, 
            e.description, 
            e.api_type,
            p.id AS project_id,
            s.endpoint,
            s.endpoint || e.path AS uri  -- Concatenate service endpoint and path
        FROM 
            api_endpoints e
        JOIN 
            api_services s ON e.api_service_id = s.id
        JOIN 
            api_projects p ON s.api_project_id = p.id
    `

	// Prepare for filtering by serviceID (if provided)
	var args []interface{}
	if serviceID != nil {
		query += " WHERE e.api_service_id = $1"
		args = append(args, *serviceID)
	}

	// Execute the query
	var endpoints []models.APIEndpoint
	err := db.Select(&endpoints, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve API endpoints: %w", err)
	}

	// Return the list of endpoints
	return endpoints, nil
}

// InsertAPIEndpoint creates a new API endpoint and returns its ID
func InsertAPIEndpoint(endpoint models.CreateAPIEndpoint, db *sqlx.DB) (int64, error) {
	// Define the SQL query for inserting a new API endpoint
	query := `
		INSERT INTO api_endpoints (
			api_service_id, 
			api_type, 
			name, 
			path, 
			method, 
			description
		) 
		VALUES (
			$1, 
			$2, 
			$3, 
			$4, 
			$5, 
			$6
		) RETURNING id
	`

	// Execute the query, passing the endpoint values as parameters
	var id int64
	err := db.Get(&id, query, endpoint.ServiceID, endpoint.APIType, endpoint.Name, endpoint.Path, endpoint.Method, endpoint.Description)
	if err != nil {
		return 0, fmt.Errorf("failed to insert API endpoint: %w", err)
	}

	// Return the ID of the newly inserted record
	return id, nil
}

// CheckAPIEndpointExist checks if an API endpoint with the same name, path, and method exists
func CheckAPIEndpointExist(name, path, method string, db *sqlx.DB) bool {
	var existingID int64
	query := `SELECT id FROM api_endpoints WHERE name = $1 AND path = $2 AND method = $3`
	err := db.QueryRow(query, name, path, method).Scan(&existingID)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("CheckAPIEndpointExist: %v", err)
		}
		return false
	}
	return true
}

// DeleteAPIEndpoint deletes an API endpoint by ID
func DeleteAPIEndpoint(id int, db *sqlx.DB) (bool, error) {
	query := `DELETE FROM api_endpoints WHERE id = $1`
	result, err := db.Exec(query, id)
	if err != nil {
		return false, fmt.Errorf("DeleteAPIEndpoint: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("DeleteAPIEndpoint: %v", err)
	}

	return rowsAffected > 0, nil
}

// GetAPIEndpointByID retrieves a single API endpoint by its ID
func GetAPIEndpointByID(id int, db *sqlx.DB) (models.APIEndpoint, error) {
	// Define the query to retrieve the API endpoint by ID, including the associated project and service information
	query := `
		SELECT 
			e.id, 
			e.api_service_id, 
			e.method, 
			e.name, 
			e.path, 
			e.description, 
			e.api_type,
			p.id AS project_id,
			s.endpoint,
			s.endpoint || e.path AS uri  -- Concatenate service endpoint and path directly in the query
		FROM 
			api_endpoints e
		JOIN 
			api_services s ON e.api_service_id = s.id
		JOIN 
			api_projects p ON s.api_project_id = p.id
		WHERE 
			e.id = $1
	`

	// Define a variable to hold the result
	var endpoint models.APIEndpoint

	// Execute the query and retrieve the result
	err := db.Get(&endpoint, query, id)
	if err != nil {
		// If no rows are found or another error occurs, return it
		if err == sql.ErrNoRows {
			return models.APIEndpoint{}, fmt.Errorf("API endpoint with ID %d not found", id)
		}
		return models.APIEndpoint{}, fmt.Errorf("failed to retrieve API endpoint: %w", err)
	}

	// Return the APIEndpoint
	return endpoint, nil
}

// UpdateAPIEndpoint updates the details of an API endpoint by ID
func UpdateAPIEndpoint(id int, updatedData models.APIEndpoint, db *sqlx.DB) (bool, error) {
	query := `UPDATE api_endpoints SET 
				api_service_id = COALESCE(NULLIF($1, 0), api_service_id),
				name = COALESCE(NULLIF($2, ''), name),
				path = COALESCE(NULLIF($3, ''), path),
				method = COALESCE(NULLIF($4, ''), method),
				description = COALESCE(NULLIF($5, ''), description),
				api_type = COALESCE(NULLIF($6, ''), api_type)
			  WHERE id = $7`

	result, err := db.Exec(query, updatedData.ServiceID, updatedData.Name, updatedData.Path, updatedData.Method, updatedData.Description, updatedData.APIType, id)
	if err != nil {
		return false, fmt.Errorf("UpdateAPIEndpoint: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("UpdateAPIEndpoint: %v", err)
	}

	return rowsAffected > 0, nil
}

// APIServiceExists checks if an API service exists
func APIServiceExists(serviceID int, db *sqlx.DB) bool {
	var id int
	query := `SELECT id FROM api_services WHERE id = $1`
	err := db.QueryRow(query, serviceID).Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("APIServiceExists: %v", err)
		}
		return false
	}
	return true
}
