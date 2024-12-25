package repositories

import (
	"blue-owl-service/models"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

func GetAPIExampleByID(db *sqlx.DB, exampleID int) (*models.APIExampleUnParsed, error) {
	// Define the SQL query
	query := `
		SELECT 
			ae.id,
			ae.api_endpoint_id,
			ae.request_header,
			ae.request_body,
			ae.request_parameter,
			ae.response_status_code,
			ae.response_body,
			ae.response_header,
			ae.test_state,
			asv.endpoint,
			aept.path,
			aept.method,
			ae.name
		FROM 
			api_examples ae
		JOIN 
			api_endpoints aept ON ae.api_endpoint_id = aept.id
		JOIN 
			api_services asv ON aept.api_service_id = asv.id
		WHERE 
			ae.id = $1;
	`

	// Prepare a variable to hold the result
	var apiExample models.APIExampleUnParsed

	// Execute the query and bind the result to the apiExample struct
	err := db.Get(&apiExample, query, exampleID)
	if err != nil {
		log.Printf("Error querying API example by ID %d: %v", exampleID, err)
		return nil, err
	}

	// Return the result
	return &apiExample, nil
}

func UpdateTestState(db *sqlx.DB, apiExampleID int, testState int) error {
	// Prepare the SQL update query
	query := `
		UPDATE api_examples
		SET test_state = $1
		WHERE id = $2
		RETURNING id, test_state
	`

	// Execute the query and retrieve the updated id and test_state
	var id int
	var updatedTestState int
	err := db.QueryRow(query, testState, apiExampleID).Scan(&id, &updatedTestState)
	if err != nil {
		log.Printf("Error updating test_state: %v", err)
		return err
	}

	// Log the update (optional)
	fmt.Printf("Updated test_state of api_example with ID: %d to %d\n", id, updatedTestState)

	// Return nil error indicating success
	return nil
}

func InsertAPIExample(endpoint models.CreateAPIEndpoint, db *sqlx.DB, apiEndpointID int) (int64, error) {
	// Serialize RequestHeader and RequestParameter to JSON strings
	requestHeaderJSON, err := json.Marshal(endpoint.RequestHeader)
	if err != nil {
		return 0, fmt.Errorf("failed to serialize requestHeader: %w", err)
	}

	requestParameterJSON, err := json.Marshal(endpoint.RequestParameter)
	if err != nil {
		return 0, fmt.Errorf("failed to serialize requestParameter: %w", err)
	}

	// Serialize ResponseHeader to JSON string
	responseHeaderJSON, err := json.Marshal(endpoint.ResponseHeader)
	if err != nil {
		return 0, fmt.Errorf("failed to serialize responseHeader: %w", err)
	}

	// Prepare the SQL query for inserting the new API example
	query := `
		INSERT INTO api_examples (
			api_endpoint_id,
			name,
			request_header,
			request_body,
			request_parameter,
			response_status_code,
			response_body,
			response_header
		) 
		VALUES (
			$1, 
			$2, 
			$3, 
			$4, 
			$5, 
			$6, 
			$7, 
			$8
		) RETURNING id
	`

	// Execute the query and insert the API example into the database
	var id int64
	err = db.Get(
		&id,
		query,
		apiEndpointID,
		"Success Example",
		string(requestHeaderJSON),
		endpoint.RequestBody,
		string(requestParameterJSON),
		endpoint.ResponseStatusCode,
		endpoint.ResponseBody,
		string(responseHeaderJSON),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert API example: %w", err)
	}

	// Return the ID of the newly inserted record
	return id, nil
}

func GetAPIExampleByEndpointID(db *sqlx.DB, endpointID int) ([]models.APIExampleUnParsed, error) {
	// Query to fetch all API examples associated with the given api_endpoint_id
	query := `
		SELECT
			id,
			api_endpoint_id,
			name, -- Include the name field
			request_header,
			request_body,
			request_parameter,
			response_status_code,
			response_body,
			response_header,
			test_state
		FROM api_examples
		WHERE api_endpoint_id = $1
	`

	// Slice to hold the result
	var examples []models.APIExampleUnParsed

	// Execute the query
	rows, err := db.Queryx(query, endpointID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch API examples for api_endpoint_id %d: %w", endpointID, err)
	}
	defer rows.Close()

	// Iterate over rows to manually parse response_body
	for rows.Next() {
		var example models.APIExampleUnParsed
		// var responseBodyString string
		// var requestBodyString string

		// Scan all fields, storing response_body as a string
		err := rows.Scan(
			&example.ID,
			&example.ApiEndpointID,
			&example.Name,
			&example.RequestHeader,
			&example.RequestBody,
			&example.RequestParameter,
			&example.ResponseStatusCode,
			&example.ResponseBody, // Fetch response_body as a string
			&example.ResponseHeader,
			&example.TestState,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// // Parse the responseBodyString into map[string]interface{}
		// if responseBodyString != "" {
		// 	// Attempt unmarshalling the response body string into a map
		// 	err = json.Unmarshal([]byte(responseBodyString), &example.ResponseBody)
		// 	if err != nil {
		// 		return nil, fmt.Errorf("failed to unmarshal response_body into map: %w", err)
		// 	}
		// } else {
		// 	// Initialize ResponseBody as an empty map if the response_body is empty
		// 	example.ResponseBody = make(map[string]interface{})
		// }

		// if requestBodyString != "" {
		// 	err = json.Unmarshal([]byte(requestBodyString), &example.RequestBody)
		// 	if err != nil {
		// 		return nil, fmt.Errorf("failed to unmarshal response_body into map: %w", err)
		// 	}
		// } else {
		// 	// Initialize ResponseBody as an empty map if the response_body is empty
		// 	example.RequestBody = make(map[string]interface{})
		// }
		// Append the populated example to the results slice
		examples = append(examples, example)
	}

	// Return the list of API examples
	return examples, nil
}
