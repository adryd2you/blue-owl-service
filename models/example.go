package models

import "time"

// APIProject represents the api_projects table

type APIExample struct {
	ID                 int                    `db:"id"`
	ApiEndpointID      int                    `db:"api_endpoint_id"`
	Endpoint           string                 `db:"endpoint"`
	Path               string                 `db:"path"`
	APIMethod          string                 `db:"method"`
	RequestHeader      string                 `db:"request_header"`
	RequestBody        map[string]interface{} `db:"request_body"`
	RequestParameter   string                 `db:"request_parameter"`
	ResponseStatusCode int                    `db:"response_status_code"`
	ResponseBody       map[string]interface{} `db:"response_body"`
	ResponseHeader     string                 `db:"response_header"`
	TestState          int                    `db:"test_state"`
	CreatedAt          time.Time              `db:"created_at"`
	Name               string                 `db:"name"`
}

type APIExampleUnParsed struct {
	ID                 int       `db:"id"`
	ApiEndpointID      int       `db:"api_endpoint_id"`
	Endpoint           string    `db:"endpoint"`
	Path               string    `db:"path"`
	APIMethod          string    `db:"method"`
	RequestHeader      string    `db:"request_header"`
	RequestBody        string    `db:"request_body"`
	RequestParameter   string    `db:"request_parameter"`
	ResponseStatusCode int       `db:"response_status_code"`
	ResponseBody       string    `db:"response_body"`
	ResponseHeader     string    `db:"response_header"`
	TestState          int       `db:"test_state"`
	CreatedAt          time.Time `db:"created_at"`
	Name               string    `db:"name"`
}
