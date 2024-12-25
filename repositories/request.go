package repositories

import (
	"blue-owl-service/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func InsertAPIRequests(requests []models.RequestContract, db *sqlx.DB) ([]int64, error) {
	// Ensure that there are requests to process
	if len(requests) == 0 {
		return nil, fmt.Errorf("no requests to insert")
	}

	// Check if the api_endpoint_id of the first request exists in the table
	firstRequest := requests[0]
	var existingID int64
	queryCheck := `
		SELECT id FROM api_requests
		WHERE api_endpoint_id = $1
		LIMIT 1
	`
	err := db.Get(&existingID, queryCheck, firstRequest.EndpointID)
	if err == nil {
		// If an error is nil, it means the api_endpoint_id already exists, so we don't insert
		return nil, fmt.Errorf("api_endpoint_id %d already exists in the table", firstRequest.EndpointID)
	}

	// Prepare the insert query for all requests
	queryInsert := `
		INSERT INTO api_requests (
			api_endpoint_id,
			parameter_name,
			data_type,
			description,
			is_required,
			default_value
		)
		VALUES (
			$1, $2, $3, $4, $5, $6
		) RETURNING id
	`

	// Prepare a slice to store the IDs of inserted requests
	var insertedIDs []int64

	// Loop through all requests and insert them
	for _, request := range requests {
		var newID int64
		err = db.Get(&newID, queryInsert, request.EndpointID, request.ParameterName, request.DataType, request.Description, request.IsRequired, request.DefaultValue)
		if err != nil {
			return nil, fmt.Errorf("failed to insert API request for parameter '%s': %w", request.ParameterName, err)
		}

		// Append the new ID to the result slice
		insertedIDs = append(insertedIDs, newID)
	}

	// Return the IDs of the newly inserted records
	return insertedIDs, nil
}

func GetAPIRequests(db *sqlx.DB, apiEndpointID int) ([]models.RequestContract, error) {
	// Query to fetch all API requests associated with the provided api_endpoint_id
	query := `
		SELECT
			id,
			api_endpoint_id,
			parameter_name,
			data_type,
			description,
			is_required,
			default_value
		FROM api_requests
		WHERE api_endpoint_id = $1
	`

	// Slice to hold the result
	var requests []models.RequestContract

	// Execute the query
	err := db.Select(&requests, query, apiEndpointID)
	if err != nil {
		// Return an error if the query execution fails
		return nil, fmt.Errorf("failed to fetch API requests for api_endpoint_id %d: %w", apiEndpointID, err)
	}

	// Return the result
	return requests, nil
}
