package models

import "encoding/json"

// APIProject represents the api_projects table
type APIProject struct {
	ID          int          `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string       `json:"name" gorm:"not null"`
	Description string       `json:"description,omitempty"`
	Status      string       `json:"status" gorm:"default:'active'"` // Default to 'active'
	APIServices []APIService `json:"api_services,omitempty" gorm:"foreignKey:APIProjectID;constraint:OnDelete:CASCADE"`
}

// APIEndpoint represents the api_endpoints table
type APIEndpoint struct {
	ID                int    `json:"id" db:"id"`
	ProjectID         int    `json:"project_id" db:"project_id"`     // Assuming 'api_project_id' is the column in the table
	ServiceID         int    `json:"service_id" db:"api_service_id"` // Assuming 'api_service_id' is the column in the table
	Method            string `json:"method" db:"method"`             // Enum: 'GET', 'POST', 'PUT', 'DELETE', 'PATCH'
	Name              string `json:"name" db:"name"`
	Endpoint          string `json:"endpoint" db:"endpoint"` // Temporarily hold the endpoint from the query
	Path              string `json:"path" db:"path"`
	URI               string `json:"uri"` // This field will store the combined URI
	Description       string `json:"description,omitempty" db:"description"`
	APIType           string `json:"api_type" db:"api_type"` // Enum: 'Internal', 'External', 'ThirdParty'
	RequestContracts  []RequestContract
	ResponseContracts []ResponseContract
	Examples          []APIExample
}

// Define a KeyValuePair struct for handling key-value pairs in headers and parameters
type KeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CreateAPIEndpoint struct {
	ServiceID          int             `json:"service_id" db:"api_service_id"`                     // The service ID that this API endpoint is part of
	Method             string          `json:"method" db:"method"`                                 // HTTP method (GET, POST, PUT, DELETE, PATCH)
	Name               string          `json:"name" db:"name"`                                     // The name of the API endpoint
	Path               string          `json:"path" db:"path"`                                     // The path of the API endpoint (e.g., /users/{id})
	Description        string          `json:"description,omitempty" db:"description"`             // A brief description of the endpoint (optional)
	APIType            string          `json:"api_type" db:"api_type"`                             // Type of API (Internal, External, or ThirdParty)
	RequestHeader      []KeyValuePair  `json:"request_header,omitempty" db:"request_header"`       // List of key-value pairs for request headers
	RequestBody        json.RawMessage `json:"request_body,omitempty" db:"request_body"`           // The body of the request
	RequestParameter   []KeyValuePair  `json:"request_parameter,omitempty" db:"request_parameter"` // List of key-value pairs for request parameters
	ResponseStatusCode int             `json:"response_status_code" db:"response_status_code"`     // HTTP status code (e.g., 200, 404)
	ResponseBody       json.RawMessage `json:"response_body,omitempty" db:"response_body"`         // The body of the response
	ResponseHeader     []KeyValuePair  `json:"response_header,omitempty" db:"response_header"`     // List of key-value pairs for response headers
}
