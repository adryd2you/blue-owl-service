package models

import "time"

type API struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Endpoint    string `json:"endpoint"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Header      string `json:"header"`
	APIType     string `json:"api_type"`
}

type TestRequest struct {
	ProjectID  int `json:"project_id"`
	ServiceID  int `json:"service_id"`
	EndpointID int `json:"endpoint_id"`
	ExampleID  int `json:"example_id"`
}

type TestResult struct {
	ID                 int       `json:"id"`                   // SERIAL PRIMARY KEY
	ExampleID          int       `json:"api_example_id"`       // Foreign Key (api_example_id) references api_examples(id)
	TestExecutedAt     time.Time `json:"test_executed_at"`     // Timestamp for when the test was executed
	IsSuccess          bool      `json:"is_success"`           // Indicates whether the test was successful
	ResponseStatusCode int       `json:"response_status_code"` // HTTP Status Code of the response
	ResponseBody       string    `json:"response_body"`        // The body of the response
	ResponseHeader     string    `json:"response_header"`      // The headers of the response
	TestMessage        string    `json:"test_message"`         // Text field for storing the result, e.g., "Success" or "Failed"
}
