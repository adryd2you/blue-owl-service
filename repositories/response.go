package repositories

import (
	"blue-owl-service/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func InsertAPIResponse(responses []models.ResponseContract, db *sqlx.DB) ([]int64, error) {
	// Prepare a slice to store the IDs of inserted responses
	var insertedIDs []int64

	// Ensure that there are responses to process
	if len(responses) == 0 {
		return nil, fmt.Errorf("no responses to insert")
	}

	// Check if the api_endpoint_id of the first response exists in the table
	firstResponse := responses[0]
	var existingID int64
	queryCheck := `
		SELECT id FROM api_responses
		WHERE api_endpoint_id = $1
		LIMIT 1
	`
	err := db.Get(&existingID, queryCheck, firstResponse.EndpointID)
	if err == nil {
		// If an error is nil, it means the api_endpoint_id already exists, so we don't insert
		return insertedIDs, fmt.Errorf("api_endpoint_id %d already exists in the table", firstResponse.EndpointID)
	}

	// Prepare the insert query for all responses
	queryInsert := `
		INSERT INTO api_responses (
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

	// Loop through all responses and insert them
	for _, response := range responses {
		var newID int64
		err = db.Get(&newID, queryInsert, response.EndpointID, response.ParameterName, response.DataType, response.Description, response.IsRequired, response.DefaultValue)
		if err != nil {
			return nil, fmt.Errorf("failed to insert API response for parameter '%s': %w", response.ParameterName, err)
		}

		// Append the new ID to the result slice
		insertedIDs = append(insertedIDs, newID)
	}

	// Return the IDs of the newly inserted records
	return insertedIDs, nil
}

func GetAPIResponses(db *sqlx.DB, apiEndpointID int) ([]models.ResponseContract, error) {
	// Query to fetch all API responses associated with the provided api_endpoint_id
	query := `
		SELECT
			id,
			api_endpoint_id,
			parameter_name,
			data_type,
			description,
			is_required,
			default_value
		FROM api_responses
		WHERE api_endpoint_id = $1
	`

	// Slice to hold the result
	var responses []models.ResponseContract

	// Execute the query
	err := db.Select(&responses, query, apiEndpointID)
	if err != nil {
		// Return an error if the query execution fails
		return nil, fmt.Errorf("failed to fetch API responses for api_endpoint_id %d: %w", apiEndpointID, err)
	}

	// Return the result
	return responses, nil
}
