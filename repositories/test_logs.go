package repositories

import (
	"blue-owl-service/models"
	"log"

	"github.com/jmoiron/sqlx"
)

func InsertTestResult(db *sqlx.DB, testResult models.TestResult) (int, error) {
	query := `
		INSERT INTO test_logs (
			api_example_id, 
			test_executed_at, 
			is_success, 
			response_status_code,
			response_body,
			response_header,
			test_message
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	// Execute the query and get the inserted ID
	var id int
	err := db.QueryRow(query,
		testResult.ExampleID,
		testResult.TestExecutedAt,
		testResult.IsSuccess,
		testResult.ResponseStatusCode,
		testResult.ResponseBody,
		testResult.ResponseHeader,
		testResult.TestMessage,
	).Scan(&id)

	if err != nil {
		log.Printf("Error inserting test result: %v", err)
		return 0, err
	}

	// Return the inserted ID and nil error
	return id, nil
}
