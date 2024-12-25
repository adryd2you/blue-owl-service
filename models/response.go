package models

type ResponseContract struct {
	ID            int    `json:"id" db:"id"`                                 // Primary key
	EndpointID    int    `json:"api_endpoint_id" db:"api_endpoint_id"`       // Foreign key to api_endpoints table
	ParameterName string `json:"parameter_name" db:"parameter_name"`         // The name of the response parameter
	DataType      string `json:"data_type" db:"data_type"`                   // The data type of the response parameter
	Description   string `json:"description,omitempty" db:"description"`     // Optional description of the response parameter
	IsRequired    bool   `json:"is_required" db:"is_required"`               // Flag to indicate if the response parameter is required
	DefaultValue  string `json:"default_value,omitempty" db:"default_value"` // Default value for the response parameter (optional)
}
